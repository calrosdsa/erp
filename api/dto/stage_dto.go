package dto

import "erp/gen/db/model"

type (
	StagesRequest struct {
		DefaultListParams
		EntityID string `query:"entity_id" required:"true"`
		Name     string `query:"name"`
	}

	StageTransitionRequest struct {
		Body StageTransitionData
	}

	EntityTransitionRequest struct {
		Body EntityTransitionData
	}

	EntityTransitionData struct {
		SourceName         string `json:"source_name"`
		SourceIndex        int32  `json:"source_index"`
		SourceStageID      int64  `json:"source_stage_id"`
		DestionationName   string `json:"destination_name"`
		DestinationIndex   int32  `json:"destination_index"`
		DestinationStageID int64  `json:"destination_stage_id"`
		ID                 int64  `json:"id"`
	}

	StageTransitionData struct {
		SourceID         int32 `json:"source_id" required:"true"`
		SourceIndex      int32 `json:"source_index" required:"true"`
		DestinationID    int32 `json:"destination_id" required:"true"`
		DestinationIndex int32 `json:"destionation_index" required:"true"`
	}

	StageDataRequest struct {
		Body StageData
	}

	StageData struct {
		ID     int32       `json:"id"`
		Fields StageFields `json:"fields"`
	}

	StageFields struct {
		Name     string `json:"name" required:"true"`
		EntityID int32  `json:"entity_id" required:"true"`
		Color    string `json:"color" required:"true"`
		Index    int32  `json:"index" required:"true"`
	}

	StageDto struct {
		ID       int32  `json:"id"`
		Name     string `json:"name"`
		Color    string `json:"color"`
		EntityID int64  `json:"entity_id"`
		Index    int32  `json:"index"`
	}
)

func StageFromModel(d model.Stage) StageDto {
	return StageDto{
		ID:   d.ID,
		Name: d.Name,
	}
}
