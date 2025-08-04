package notification_ucase

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/proto"
	"erp/pkg/logger"
	"erp/pkg/ws"
	notification_repo "erp/project/core/notification/repository"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type NotificationUcase interface {
	GetNotifications(req *common.RequestContext, d dto.NotificationsRequest) (
		res dto.ResponseDataList[[]dto.NotificationDto], err error)
	NotifiationCount(req *common.RequestContext) (count int64, err error)
}

type notificationUcase struct {
	emitLog   logger.EmitLog
	repo      notification_repo.NotificationRepository
	scheduler gocron.Scheduler
	ws        ws.WsConn
}

func NewNotificationUcase(
	logger logger.Logger,
	repo notification_repo.NotificationRepository,
	scheduler gocron.Scheduler,
	ws ws.WsConn,
) NotificationUcase {
	usecase := &notificationUcase{
		emitLog:   logger.EmitLog("notification-ucase"),
		repo:      repo,
		scheduler: scheduler,
		ws:        ws,
	}
	usecase.initNotificationScheduler()
	return usecase
}

func (s *notificationUcase) NotifiationCount(req *common.RequestContext) (count int64, err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetNotifications"))
		}
	}()
	count,err = s.repo.NotifiationCount(req)
	return
}

func (s *notificationUcase) GetNotifications(req *common.RequestContext, d dto.NotificationsRequest) (
	res dto.ResponseDataList[[]dto.NotificationDto], err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetNotifications"))
		}
	}()
	res.Body.Result, err = s.repo.GetNotifications(req, d)
	return
}

func (s *notificationUcase) initNotificationScheduler() {
	job, _ := s.scheduler.NewJob(
		gocron.DurationJob(
			5*time.Second,
		),
		gocron.NewTask(func() {
			notifications, err := s.repo.UpdateSendedNotifications()
			if err != nil {
				s.emitLog.Err(err, logger.OptionsLog.WithMethod("initNotificationScheduler"))
			}
			for _, notification := range notifications {
				s.ws.PublishToSubscribers(context.Background(), []int64{notification.ProfileID},
					notification.Payload, proto.MessageType_Notification)
			}
		}),
	)
	if err := job.RunNow(); err != nil {
		fmt.Println("FAIL TO RUN JOB", err)
	}

}

func (s *notificationUcase) BatchCreateNotifications(ctx context.Context, notifications []*model.Notification) (err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("BatchCreateNotifications"))
		}
	}()
	err = s.repo.BatchCreateNotifications(ctx, notifications)
	return
}
