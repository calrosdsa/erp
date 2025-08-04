package entity

type CompanyPlugins struct {
	CompanyID int `gorm:"primaryKey;not null;default:null" required:"false"`
	Plugin string `gorm:"primaryKey;not null;default:null"`
	Credentials string `required:"false"`
}
