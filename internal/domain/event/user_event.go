package event

import (
	"erp/api/common"
	"erp/gen/db/model"
	"erp/gen/db/query"
)

type UserCreatedEventData struct {
	LanguageCode common.LanguageCode
	UseRelation model.UserRelation
	Tx *query.QueryTx
}
