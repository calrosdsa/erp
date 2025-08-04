package event

import (
	"erp/gen/db/model"
	"erp/gen/db/query"
)

type StatusJournalEntryEventData struct {
	JournalEntry model.JournalEntry
	Tx           *query.QueryTx
	Lines      []*model.JournalEntryLine
}
