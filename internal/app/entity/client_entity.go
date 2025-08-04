package entity

import (
	"erp/gen/db/model"
	"time"

	"gorm.io/gorm"
)

type Client struct {
	// PartyID        uint           `gorm:"primarykey"`
	// Base
	ID               uint `gorm:"primarykey"`
	// Party            Party `gorm:"references:ID"`
	CreatedAt        time.Time      `gorm:"default:now()"`
	DeletedAt        gorm.DeletedAt `gorm:"index" required:"false"`
	UpdatedAt        time.Time      `gorm:"default:now()"`
	Code             string
	Uuid             string `gorm:"type:uuid;default:uuid_generate_v4()"`
	GivenName        string `gorm:"not null;default:null"`
	FamilyName       string `gorm:"not null;default:null"`
	OrganizationName string
	EmailAddress     string `gorm:"not null;default:null"`
	PhoneNumber      string
	CountryCode      string
	UserID           uint `gorm:"not null;default:null"`
	// User             User `json:"-"`
	CompanyID        uint `gorm:"not null;default:null"`
	Company          model.Company
	Payload          interface{} `gorm:"-" json:"-"`

	ClientKeyValueData []ClientKeyValueData `gorm:"foreignKey:BaseID"`
}

type ClientKeyValueData struct {
	ID     uint   `gorm:"primarykey" required:"false"`
	Key    string `json:"key"`
	Value  string `json:"value"`
	BaseID uint   `json:"baseId" required:"false"`
}
