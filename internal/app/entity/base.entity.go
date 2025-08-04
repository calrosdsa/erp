package entity

import (
	"time"

	"gorm.io/gorm"
)

// `gorm:"unique;not null;type:varchar(100);default:null"`
type Base struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:"default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index" required:"false"`
	UpdatedAt time.Time      `gorm:"default:now()"`
	// Uuid string `gorm:"type:uuid;default:uuid_generate_v4()"`
}

type BaseWithoutUuid struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:"default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index" required:"false"`
	UpdatedAt time.Time      `gorm:"default:now()"`
}

type BaseWithoutID struct {
	CreatedAt time.Time      `gorm:"default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index" required:"false"`
	UpdatedAt time.Time      `gorm:"default:now()"`
}

type TranslationBase struct {
	LanguageCode string `gorm:"not null;default:null"`
	BaseId       uint   `gorm:"not null;default:null"`
}

type EntityPluginBase struct {
	Plugin string `gorm:"primaryKey;"`
	BaseID uint   `gorm:"primaryKey;" required:"true"`
	Data   string `gorm:"-"`
}

type Ordinal int32

const (
	PARENT_ORDINAL     Ordinal = 0
	SUB_PARENT_ORDINAL Ordinal = 1
	DEPTH_TWO Ordinal = 2
)



var (
	COMPANY_ENTITY_ID = 1
	ITEM_ENTITY_ID = 2
    ITEM_PRICE_ENTITY_ID = 3
    ITEM_GROUP_ENTITY_ID = 4
    ITEM_STOCK_ENTITY_ID = 5
    ITEM_ATTRIBUTES_ENTITY_ID = 6
    ITEM_WAREHOUSE_ENTITY_ID = 7
    TAX_ENTITY_ID = 8
    PRICE_LIST_ENTITY_ID = 9
	ROLE_ENTITY_ID = 10	
	USERS_ENTITY_ID = 11
)