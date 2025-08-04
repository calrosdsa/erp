package domain

import (
	"encoding/json"
	"erp/gen/db/query"
)

func CloseTx(tx *query.QueryTx, err error) error {
	if p := recover(); p != nil {
		_ = tx.Rollback()
		panic(p)
	} else if err != nil {
		_ = tx.Rollback()
		return err
	} else {
		return tx.Commit()
	}
}

func JsonStringify(data map[string]interface{}) string {
	jsonBytes, _ := json.Marshal(data)
	return string(jsonBytes)
}
