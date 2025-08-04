package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	DealsRequest struct {
		DefaultListParams
		CreatedAt string `query:"created_at" required:"false"`
		UpdatedAt string `query:"updated_at" required:"false"`
	}

	DealDataRequest struct {
		Body DealData
	}

	DealData struct {
		ID           int64             `json:"id"`
		Fields       DealFields        `json:"fields"`
		Contacts     []ContactData     `json:"contacts"`
		Participants []ParticipantData `json:"participants"`
	}
	ParticipantData struct {
		ID     int64  `json:"id"`
		Action string `json:"action" required:"false"`
	}

	DealFields struct {
		Name                 string     `json:"name" required:"true"`
		StageID              int64      `json:"stage_id" required:"true"`
		Amount               int64      `json:"amount" required:"true"`
		Currency             string     `json:"currency" required:"true"`
		DealType             *string    `json:"deal_type" required:"false"`
		Source               *string    `json:"source" required:"false"`
		SourceInformation    *string    `json:"source_information" required:"false"`
		StartDate            time.Time  `json:"start_date"`
		EndDate              *time.Time `json:"end_date" required:"false"`
		ResponsibleID        int64      `json:"responsible_id"`
		AvailableForEveryone bool       `json:"available_for_everyone"`
		Index                int32      `json:"index" required:"true"`
		CustomerID           *int64      `json:"customer_id" required:"false"`
	}

	DealDetailDto struct {
		Deal         DealDto      `json:"deal"`
		Participants []ProfileDto `json:"participants"`
	}

	DealDto struct {
		ID        int64     `json:"id"`
		UUID      string    `json:"uuid"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`

		StageID    int32  `json:"stage_id"`
		StageIndex int32  `json:"stage_index"`
		Stage      string `json:"stage"`

		Amount               int64      `json:"amount"`
		Currency             string     `json:"currency"`
		DealType             *string    `json:"deal_type"`
		Source               *string    `json:"source"`
		SourceInformation    *string    `json:"source_information"`
		StartDate            time.Time  `json:"start_date"`
		EndDate              *time.Time `json:"end_date"`
		AvailableForEveryone bool       `json:"available_for_everyone"`

		ResponsibleID         int64  `json:"responsible_id"`
		ResponsibleGivenName  string `json:"responsible_given_name"`
		ResponsibleFamilyName string `json:"responsible_family_name"`
		ResponsibleUUID       string `json:"responsible_uuid"`

		Customer *string `json:"customer"`
		CustomerID *int64 `json:"customer_id"`

		Index int32 `json:"index"`
	}
)

func DealFromModel(d model.Deal) DealDto {
	return DealDto{
		ID:   d.ID,
		UUID: d.UUID,
		Name: d.Name,
	}
}
