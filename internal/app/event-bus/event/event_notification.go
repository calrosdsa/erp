package event

import "erp/api/common"

type NotificationData struct {
	NotificationEventType string
	Data NotificationPayload
}

type NotificationPayload struct {
	RequestContext common.RequestContext	
	Payload        interface{}
}

type NotificationType string

const (
	NOTIFY_NEW_CLIENT = "notify_new_client"
)

const (
	EMAIL_TYPE NotificationType = "email"
)

const (
	NOTIFICATION_EVENT = "email_event"
)
