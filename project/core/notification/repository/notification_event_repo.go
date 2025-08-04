package notification_repo

import (
	"context"
	"erp/gen/db/model"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
)

type NotificationEventRepo interface {
	HandlerActivityNotification(ctx context.Context, payload event.ActivityEventData) error
}

type notificationEventRepo struct {
	locale helpers.Locale
}

func NewNotificationEventRepo(
	helpers *helpers.Helpers,
) NotificationEventRepo {
	return &notificationEventRepo{
		locale: helpers.Locale,
	}
}

func (r *notificationEventRepo) HandlerActivityNotification(ctx context.Context, payload event.ActivityEventData) (err error) {
	switch payload.Data.Type {
	case proto.ActivityType_ACTIVITY.String():
		err = r.handleDeadLineNotification(ctx, payload)
	case proto.ActivityType_COMMENT.String():
		err = r.handleCommentNotification(ctx, payload)
	}
	return
}

func (r *notificationEventRepo) handleCommentNotification(ctx context.Context, payload event.ActivityEventData) (err error) {
	tx := payload.Tx
	data := payload.Data
	entity, err := tx.Entity.WithContext(ctx).Where(
		tx.Entity.ID.Eq(int64(data.EntityID)),
	).First()
	if err != nil {
		return
	}
	notificationText := r.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Notification.ActivityComment"),
		helpers.OptionsLocale.WithLang(string(payload.ReqCtx.LanguageCode)),
		helpers.OptionsLocale.WithTemplate(map[string]interface{}{
			"EntityName": entity.Name,
			"Comment":    data.ActivityComment.Fields.Comment,
			"Name":       data.PartyName,
		}),
	)
	mentions := make([]*model.Mention, len(data.ActivityComment.Mentions))
	for i, commentMention := range data.ActivityComment.Mentions {
		notification := &model.Notification{
			Type:      proto.NotificationType_MENTION.String(),
			ProfileID: commentMention.Fields.ProfileID,
			CompanyID: payload.ReqCtx.ActiveCompany.ID,
			SendAt:    time.Now(),
			Payload:   notificationText,
		}
		err = tx.Notification.WithContext(ctx).Save(notification)
		if err != nil {
			return
		}
		startIndex, endIndex := r.findSubstringIndices(notificationText, payload.Data.PartyName)
		notificationMention := &model.Mention{
			EntityID:   int32(entity.ID),
			PartyID:    strconv.Itoa(int(data.PartyID)),
			PartyName:  payload.Data.PartyName,
			StartIndex: int32(startIndex),
			EndIndex:   int32(endIndex),
			ReferenceID: notification.ID,
		}
		mentions[i] = notificationMention
	}
	err = tx.Mention.WithContext(ctx).CreateInBatches(mentions,domain.DEFAULT_BATCH_SIZE)
	
	return
}

func (r *notificationEventRepo) handleDeadLineNotification(ctx context.Context, payload event.ActivityEventData) (err error) {
	if lo.IsEmpty(&payload.Data.ActivityDeadLine.Fields.Deadline) {
		return nil
	}
	tx := payload.Tx
	data := payload.Data
	entity, err := tx.Entity.WithContext(ctx).Where(
		tx.Entity.ID.Eq(int64(data.EntityID)),
	).First()
	if err != nil {
		return
	}
	notificationText := r.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Notification.ActivityReminder"),
		helpers.OptionsLocale.WithLang(string(payload.ReqCtx.LanguageCode)),
		helpers.OptionsLocale.WithTemplate(map[string]interface{}{
			"ActivityName": data.ActivityDeadLine.Fields.Title,
			"EntityName":   entity.Name,
			"Name":         payload.Data.PartyName,
			"Deadline":     data.ActivityDeadLine.Fields.Deadline.Format(time.DateTime),
		}),
	)
	if(data.ActivityDeadLine.Fields.ProfileID == nil) {
		return
	}

	notification := model.Notification{
		Type:      proto.NotificationType_MENTION.String(),
		ProfileID: *data.ActivityDeadLine.Fields.ProfileID,
		CompanyID: payload.ReqCtx.ActiveCompany.ID,
		SendAt:    payload.Data.ActivityDeadLine.Fields.Deadline,
		Payload:   notificationText,
	}
	err = tx.Notification.Save(&notification)
	if err != nil {
		return
	}
	startIndex, endIndex := r.findSubstringIndices(notificationText, payload.Data.PartyName)
	notificationMention := model.Mention{
		EntityID:   int32(entity.ID),
		PartyID:    strconv.Itoa(int(data.PartyID)),
		PartyName:  payload.Data.PartyName,
		StartIndex: int32(startIndex),
		EndIndex:   int32(endIndex),
		ReferenceID: notification.ID,
	}
	err = tx.Mention.Save(&notificationMention)
	return err
}

func (r *notificationEventRepo) findSubstringIndices(s, target string) (int, int) {
	if target == "" {
		return 0, 0
	}
	startIndex := strings.Index(s, target)
	if startIndex == -1 {
		return -1, -1
	}
	return startIndex, startIndex + len(target) - 1
}
