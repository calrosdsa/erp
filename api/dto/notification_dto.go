package dto

import "time"

type (
	NotificationsRequest struct {
		DefaultListParams
	}

	NotificationCountDto struct {
		Body struct {
			Count int64 `json:"count"`
		}
	}

	NotificationDto struct {
		ID        int64        `json:"id"`
		Payload   string       `json:"payload"`
		Type      string       `json:"type"`
		SendAt    time.Time    `json:"send_at"`
		ProfileGN string       `json:"profile_gn"`
		ProfileFN string       `json:"profile_fn"`
		ProfileID int64        `json:"profile_id"`
		Mentions  []MentionDto `gorm:"-" json:"mentions"`
	}

	MentionDto struct {
		ReferenceID int64  `json:"reference_id"`
		EntityName  string `json:"entity_name"`
		EntityID    int    `json:"entity_id"`
		EntityHref  string `json:"entity_href"`
		HasModal bool `json:"has_modal"`
		PartyID     string `json:"party_id"`
		PartyName   string `json:"party_name"`
		StartIndex  int    `json:"start_index"`
		EndIndex    int    `json:"end_index"`
	}
)
