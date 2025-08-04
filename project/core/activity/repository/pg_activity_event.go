package activity_repo

import (
	"context"
	"encoding/json"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/domain/event"
	"erp/pkg/db"
)

type ActivityEventRepository interface {
	OnStagePartyChange(ctx context.Context, d event.ChangeStageEventData) error
}

type activityEventRepo struct {
	Q *query.Query
}

func NewActivityEventRepo(
	conn db.Connection,
) ActivityEventRepository {
	return &activityEventRepo{}
}

func (r *activityEventRepo) OnStagePartyChange(ctx context.Context, d event.ChangeStageEventData) (err error) {
	tx := d.Tx
	activity := model.Activity{}
	activity.PartyID = d.StageTransition.ID
	activity.ProfileID = d.ProfileID
	activity.Type = proto.ActivityType_STAGE.String()
	if d.StageTransition.SourceName != d.StageTransition.DestionationName {
		activityStageData := map[string]interface{}{
			"source":      d.StageTransition.SourceName,
			"destination": d.StageTransition.DestionationName,
		}
		jsonBytes, _ := json.Marshal(activityStageData)
		jsonString := string(jsonBytes)
		activity.Data = &jsonString
	}
	err = tx.Activity.Save(&activity)
	return
}
