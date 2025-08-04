package entitysquare

import (
	"erp/gen/db/model"
	"time"

	"gorm.io/gorm"
)

type SquareObject struct {
	CreatedAt         time.Time      `gorm:"default:now()"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	UpdatedAt         time.Time      `gorm:"default:now()"`
	ItemGroupId       uint           `gorm:"primaryKey"`
	ObjectVariationId string         `gorm:"primaryKey"`
	ObjectId          string
	ItemPriceID       uint
	ItemPrice         model.ItemPrice
}

