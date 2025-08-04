package entitysquare

import "erp/gen/db/model"

type SquareOrder struct {
	SquareOrderId string `gorm:"primaryKey"`
	SalesOrderID  uint   `gorm:"primaryKey"`
	Party         model.Party
	PartyID       uint
}
