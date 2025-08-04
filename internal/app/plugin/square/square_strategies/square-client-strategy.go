package squarestrategies

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"erp/api/common"
// 	clientconfig "erp/internal/app/config/client_config"
// 	"erp/internal/app/connection"
// 	"erp/internal/app/entity"
// 	entitysquare "erp/internal/app/plugin/square/entitiy_square"
// 	squareservice "erp/internal/app/plugin/square/square_service"
// 	squaretypes "erp/internal/app/plugin/square/square_types"
// 	"errors"
// 	"fmt"
// 	"net/http"
// 	"time"

// 	"gorm.io/gorm"
// )

// var ( //interface comformance checks
// 	_ clientconfig.ClientStrategy = (*SquareClientStrategy)(nil)
// )

// type SquareClientStrategy struct {
// 	conn          *connection.Connection
// 	squareService *squareservice.SquareService
// 	timeout       *time.Duration
// }

// func NewSquareClientStrategy(conn *connection.Connection, squareService *squareservice.SquareService, timeout *time.Duration) *SquareClientStrategy {
// 	return &SquareClientStrategy{
// 		conn:          conn,
// 		squareService: squareService,
// 		timeout:       timeout,
// 	}
// }

// func (s *SquareClientStrategy) CreateCustomer(req *common.RequestContext, d *entity.Client, metadata string) error {
// 	var squareMetadata squaretypes.SquareCustomerMetadata
// 	err := json.Unmarshal([]byte(metadata), &squareMetadata)
// 	if err != nil {
// 		return err
// 	}
// 	if squareMetadata.Type == squaretypes.SQUARE_TYPE_SUBSCRIPTION {
// 		credentials, err := s.squareService.GetCredentials(req)
// 		if err != nil {
// 			return err
// 		}
// 		customer, err := s.createOrRetrieveSquareCustomer(credentials, d)
// 		if err != nil {
// 			return err
// 		}
// 		if metadata == "" {
// 			return err
// 		}
// 		err = s.squareService.CreateSubscriptionSquare(&customer, &credentials, &squareMetadata)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (s *SquareClientStrategy) createOrRetrieveSquareCustomer(credentials squaretypes.SquareCredentials, e *entity.Client) (
// 	res squaretypes.SquareCustomerResponse, err error) {
// 	var squareCustomer entitysquare.SquareCustomer
// 	err = s.conn.Db.Where("client_id = ?", e.ID).First(&squareCustomer).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			res, err = s.createSquareCustomer(credentials, e)
// 			return
// 		}
// 	}
// 	res.Customer.ID = squareCustomer.CustomerId
// 	return
// }

// func (s *SquareClientStrategy) createSquareCustomer(credentials squaretypes.SquareCredentials, e *entity.Client) (
// 	res squaretypes.SquareCustomerResponse, err error) {
// 	url := s.squareService.GetBaseUrl() + "customers"
// 	customerRequest := &squaretypes.SquareCustomerRequest{}
// 	customerRequest.FromEntity(e)
// 	body, err := json.Marshal(customerRequest)
// 	if err != nil {
// 		return
// 	}
// 	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
// 	if err != nil {
// 		return
// 	}
// 	r.Header.Add("Content-Type", "application/json")
// 	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", credentials.AccessToken))
// 	r.Header.Add("Square-Version", credentials.ApiVersion)
// 	client := &http.Client{
// 		Timeout: *s.timeout,
// 	}
// 	resp, err := client.Do(r)
// 	if err != nil {
// 		return
// 	}

// 	defer resp.Body.Close()
// 	derr := json.NewDecoder(resp.Body).Decode(&res)
// 	if derr != nil {
// 		return
// 	}
// 	fmt.Println("CREATE SQUARE CLIENT", resp.Status, res)

// 	var squareCustomer entitysquare.SquareCustomer

// 	squareCustomer.CustomerId = res.Customer.ID
// 	squareCustomer.PartyID = e.ID
// 	fmt.Println("SQUARE CUSTOMER", "CLIENT", e.ID, squareCustomer)
// 	err = s.conn.Db.Save(&squareCustomer).Error
// 	fmt.Println("SQUARE CUSTOMER CREATED")
// 	if err != nil {
// 		s.deleteCustomer(&res, credentials)
// 		//ROLLBACK
// 		//DELETE CUSTOMER SQUARE
// 		return
// 	}

// 	return
// }

// func (s SquareClientStrategy) deleteCustomer(customer *squaretypes.SquareCustomerResponse, credentials squaretypes.SquareCredentials) (err error) {
// 	client := &http.Client{
// 		Timeout: *s.timeout,
// 	}
// 	url := s.squareService.GetBaseUrl() + fmt.Sprintf("customers/%s", customer.Customer.ID)
// 	req, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, url, nil)
// 	if err != nil {
// 		fmt.Println("Failed to create request", err)
// 		return
// 	}

// 	// Set the headers
// 	req.Header.Add("Content-Type", "application/json")
// 	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", credentials.AccessToken))
// 	req.Header.Add("Square-Version", credentials.ApiVersion)

// 	// Send the request
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Failed to send request", err)
// 		return
// 	}
// 	defer resp.Body.Close()
// 	// Read the response
// 	// Print the status code and response body
// 	fmt.Printf("Response status: %s\n", resp.Status)
// 	return
// }
