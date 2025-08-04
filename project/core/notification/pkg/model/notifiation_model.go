package notification_model

type NotificationPayload struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}