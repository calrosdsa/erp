package entitysquare

import (
	"erp/gen/db/model"
)

type SquareCustomer struct {
	CustomerId string         `gorm:"primaryKey;not null;default:null"`
	PartyID    uint           `gorm:"primaryKey;not null;default:null"`
	Party    model.Party `gorm:"foreignKey:PartyID;constraint:OnDelete:CASCADE;"`
}
