package ws

import (
	"context"
	"encoding/json"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/repository"
	notification_model "erp/project/core/notification/pkg/model"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/coder/websocket"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)


type WsConn interface {
	PublishToSubscribers(ctx context.Context,ids []int64,msg string,messageType proto.MessageType) 
}

// chatServer enables broadcasting to a set of subscribers.
type chatServer struct {
	// subscriberMessageBuffer controls the max number
	// of messages that can be queued for a subscriber
	// before it is kicked.
	//
	// Defaults to 16.
	subscriberMessageBuffer int

	// publishLimiter controls the rate limit applied to the publish endpoint.
	//
	// Defaults to one publish every 100ms with a burst of 8.
	publishLimiter *rate.Limiter

	// logf controls where logs are sent.
	// Defaults to log.Printf.
	logf func(f string, v ...interface{})

	// serveMux routes the various endpoints to the appropriate handler.

	subscribersMu sync.Mutex
	subscribers   map[*subscriber]struct{}
	jwt helpers.JwtHelper
	session repository.SessionService
}

// newChatServer constructs a chatServer with the defaults.
func NewChatServer(e *echo.Echo,helpers *helpers.Helpers,session repository.SessionService) WsConn{
	cs := &chatServer{
		subscriberMessageBuffer: 16,
		logf:                    log.Printf,
		subscribers:             make(map[*subscriber]struct{}),
		publishLimiter:          rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
		jwt: helpers.Jwt,
		session: session,
	}
	e.POST("/ws/publish", cs.Publish)
	e.GET("/ws/subscribe", func(c echo.Context) error {
		cs.Subscribe(c)
		return nil
	})
	return cs
}

// subscriber represents a subscriber.
// Messages are sent on the msgs channel and if the client
// cannot keep up with the messages, closeSlow is called.
type subscriber struct {
	id int64
	msgs      chan []byte
	closeSlow func()
}


func (cs *chatServer) PublishToSubscribers(ctx context.Context,ids []int64,msg string,messageType proto.MessageType){
	cs.subscribersMu.Lock()
	defer cs.subscribersMu.Unlock()
	cs.publishLimiter.Wait(ctx)
	var subscribers []*subscriber
	for subscriber := range cs.subscribers {
		if subscriber != nil && lo.Contains(ids,subscriber.id){
			subscribers = append(subscribers, subscriber)
		}
	}

	message := notification_model.NotificationPayload{
		Type: messageType.String(),
		Message: msg,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
	  fmt.Println("Error converting struct to JSON:", err)
	}

	fmt.Println("SUNSCRIBERS",len(subscribers))
	for _,s := range subscribers {
		select {
		case s.msgs <- jsonData:
		default:
			go s.closeSlow()
		}
	}
}

func (cs *chatServer) Subscribe(c echo.Context) (error) {
	token := c.QueryParam("token")
	sessionUUID := c.QueryParam("session_uuid")
	fmt.Println("TOKEN",token)
	_,err := cs.jwt.ExtractClaimsAdmin(token)
	if err != nil {
		return nil
	}
	userRelation, err := cs.session.GetUserRelation(context.Background(), sessionUUID)
	if err != nil {
		return nil
	}

	fmt.Println("PROFILE",sessionUUID,userRelation.ProfileID)
	
	err = cs.subscribe(c.Response().Writer, c.Request(),userRelation.ProfileID)
    if errors.Is(err, context.Canceled) || 
        websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
        websocket.CloseStatus(err) == websocket.StatusGoingAway {
        return nil // Ignore normal closure/cancellation
    }
    if err != nil {
        cs.logf("WebSocket error: %v", err)
        return nil // Suppress error to prevent Echo from writing to hijacked connection
    }
    return nil
}
func (cs *chatServer) Publish(c echo.Context) error {
	id := c.Param("id")
	body := http.MaxBytesReader(c.Response().Writer, c.Request().Body, 8192)
	msg, err := io.ReadAll(body)
	fmt.Println("MESSAGE", string(msg))
	if err != nil {
		return c.String(http.StatusRequestEntityTooLarge, err.Error())
	}
	cs.publish(msg,id)
	return c.String(http.StatusOK, "Publish")
}




// subscribe subscribes the given WebSocket to all broadcast messages.
// It creates a subscriber with a buffered msgs chan to give some room to slower
// connections and then registers the subscriber. It then listens for all messages
// and writes them to the WebSocket. If the context is cancelled or
// an error occurs, it returns and deletes the subscription.
//
// It uses CloseRead to keep reading from the connection to process control
// messages and cancel the context if the connection drops.
func (cs *chatServer) subscribe(w http.ResponseWriter, r *http.Request,id int64) error {
	var mu sync.Mutex
	var c *websocket.Conn
	var closed bool
	s := &subscriber{
		msgs: make(chan []byte, cs.subscriberMessageBuffer),
		closeSlow: func() {
			mu.Lock()
			defer mu.Unlock()
			closed = true
			if c != nil {
				c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
			}
		},
		id: id,
	}
	cs.addSubscriber(s)
	defer cs.deleteSubscriber(s)

	c2, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return err
	}
	mu.Lock()
	if closed {
		mu.Unlock()
		return net.ErrClosed
	}
	c = c2
	mu.Unlock()
	defer c.CloseNow()

	ctx := c.CloseRead(context.Background())

	for {
		select {
		case msg := <-s.msgs:
			err := writeTimeout(ctx, time.Second*5, c, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// publish publishes the msg to all subscribers.
// It never blocks and so messages to slow subscribers
// are dropped.
func (cs *chatServer) publish(msg []byte,id string) {
	cs.subscribersMu.Lock()
	defer cs.subscribersMu.Unlock()
	cs.publishLimiter.Wait(context.Background())
	fmt.Println("SUBSCRIBES",cs.subscribers)
	for s := range cs.subscribers {
		if s != nil {
			fmt.Println("SUBSCRIBER ID",s.id)
		}
		select {
		case s.msgs <- msg:
		default:
			go s.closeSlow()
		}
	}
}

// addSubscriber registers a subscriber.
func (cs *chatServer) addSubscriber(s *subscriber) {
	cs.subscribersMu.Lock()
	cs.subscribers[s] = struct{}{}
	cs.subscribersMu.Unlock()
}

// deleteSubscriber deletes the given subscriber.
func (cs *chatServer) deleteSubscriber(s *subscriber) {
	cs.subscribersMu.Lock()
	delete(cs.subscribers, s)
	cs.subscribersMu.Unlock()
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)
}
