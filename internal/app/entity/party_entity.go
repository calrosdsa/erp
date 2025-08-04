package entity

// type PartyType struct {
// 	Code string `gorm:"primaryKey"`
// 	Name string 
// 	Type string 
// }

const (
	PARTY_ADMIN = "admin"
)

// type Party struct {
// 	Base		 
// 	PartyTypeCode string 
// 	PartyType PartyType `gorm:"references:Code"`

// }

type PartyTypeCode string

const (
	PARTY_CLIENT_CODE PartyTypeCode = "CLIENT"
	PARTY_ORG_CODE PartyTypeCode = "ORG"
)