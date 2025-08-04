package notification_repo

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/domain"
	"erp/pkg/db"
	"fmt"

	"github.com/samber/lo"
)

type NotificationRepository interface {
	BatchCreateNotifications(ctx context.Context, notifications []*model.Notification) error
	UpdateSendedNotifications() (res []dto.NotificationDto, err error)
	GetNotifications(req *common.RequestContext, d dto.NotificationsRequest) (res []dto.NotificationDto, err error)
	NotifiationCount(req *common.RequestContext)(count int64,err error)
}

type notificationRepo struct {
	Q *query.Query
}

func NewNotificationRepo(
	db db.Connection,
) NotificationRepository {
	return &notificationRepo{
		Q: db.GetQ(),
	}
}
func (r *notificationRepo) NotifiationCount(req *common.RequestContext)(count int64,err error){
	notificationQ := r.Q.Notification
	count,err = notificationQ.WithContext(req.Ctx).Where(
		notificationQ.CompanyID.Eq(req.CompanyDefaults.CompanyID),
		notificationQ.Read.Is(false),
		notificationQ.ProfileID.Eq(req.Profile.ID),
		notificationQ.Sended.Is(true),
	).Count()
	return 
}


func (r *notificationRepo) BatchCreateNotifications(ctx context.Context, notifications []*model.Notification) (err error) {
	err = r.Q.Notification.WithContext(ctx).CreateInBatches(notifications, domain.DEFAULT_BATCH_SIZE)
	return
}

func (r *notificationRepo) UpdateSendedNotifications() (res []dto.NotificationDto, err error) {
	query := `
	WITH selected AS (
    SELECT n.id,n.payload,n.type,n.send_at,n.profile_id
    FROM notifications as n
	WHERE n.send_at < now() and n.sended = false
    ORDER BY send_at desc
    LIMIT 100
    FOR UPDATE 
	)
	UPDATE notifications
	SET sended = true
	FROM selected
	WHERE notifications.id = selected.id
	RETURNING selected.*
	`
	err = r.Q.Notification.UnderlyingDB().Raw(query).Scan(&res).Error
	if err != nil {
		return
	}
	return
}

func (r *notificationRepo) GetNotifications(req *common.RequestContext, d dto.NotificationsRequest) (res []dto.NotificationDto, err error) {
	query := `
	WITH selected AS (
    SELECT n.id,n.payload,n.type,n.send_at,
	 p.given_name as profile_gn, p.family_name as profile_fn,
	 n.profile_id
    FROM notifications as n
	JOIN profiles as p on p.id = n.profile_id
	WHERE n.profile_id = $1 and send_at < now()
    ORDER BY send_at desc
    LIMIT 100
    FOR UPDATE 
	)
	UPDATE notifications
	SET read = true
	FROM selected
	WHERE notifications.id = selected.id
	RETURNING selected.*
	`
	err = r.Q.Notification.WithContext(req.Ctx).UnderlyingDB().Raw(query, req.Profile.ID).Scan(&res).Error
	if err != nil {
		return
	}
	notificationIds := lo.Map(res, func(n dto.NotificationDto, i int) int64 {
		return n.ID
	})
	var notificationMentions []dto.MentionDto
	entityQ := r.Q.Entity
	nmQ := r.Q.Mention
	err = nmQ.WithContext(req.Ctx).Select(
		entityQ.Name.As("entity_name"), entityQ.Href.As("entity_href"),entityQ.HasModal,
		nmQ.EntityID.As("entity_id"), nmQ.PartyID, nmQ.PartyName,
		nmQ.StartIndex, nmQ.EndIndex, nmQ.ReferenceID,
	).Join(
		entityQ, entityQ.ID.EqCol(nmQ.EntityID),
	).Where(
		nmQ.ReferenceID.In(notificationIds...),
	).Scan(&notificationMentions)
	if err != nil {
		return
	}
	fmt.Println("NOTIFICATION MENTIONS", notificationMentions)
	groups := lo.GroupBy(notificationMentions, func(i dto.MentionDto) int64 {
		return i.ReferenceID
	})
	for i, notification := range res {
		if mentions, ok := groups[notification.ID]; ok {
			res[i].Mentions = mentions
		}
	}
	return
}
