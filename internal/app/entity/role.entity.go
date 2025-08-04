package entity

// import (
// 	"erp/gen/db/model"

// 	"gorm.io/gorm"
// )

// type Role struct {
// 	Base
// 	Code        string `gorm:"unique;not null;default:null"`
// 	Description string
// 	RoleActions []RoleActions
// 	CompanyID   int64
// 	Company     model.Company
// }

// type RoleActions struct {
// 	RoleID    uint `gorm:"primaryKey"`
// 	Role      Role `json:"-"`
// 	Action    Action
// 	ActionID  uint           `gorm:"primaryKey"`
// 	DeletedAt gorm.DeletedAt `gorm:"index" required:"false" json:"-"`
// }

// type Entity struct {
// 	ID   uint `gorm:"primaryKey"`
// 	Name string
// }

// type Action struct {
// 	ID       uint `gorm:"primaryKey"`
// 	Name     string
// 	EntityID uint
// 	Entity   Entity
// }

const (
	ROLE_CLIENT = "client"
	ROLE_ADMIN  = "admin"
)


