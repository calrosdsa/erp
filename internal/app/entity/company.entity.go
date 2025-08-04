package entity


// type Company struct {
// 	Base
// 	Uuid string `gorm:"type:uuid;default:uuid_generate_v4()"`
// 	Code string 
// 	Name     string `gorm:"not null;default:null"`
// 	IsParent bool
// 	ParentID *int64
// 	Parent   *Company    `gorm:"foreignKey:ParentID;references:ID"`
// 	Users    []User `gorm:"many2many:user_companies;"`
// 	Ordinal Ordinal 
// 	CompanyDepartments []*Company `gorm:"many2many:company_departments"`
// 	// ItemGroups    []ItemGroup `gorm:"many2many:company_item_groups;"`
// 	CompanyPlugins []CompanyPlugins
	
// 	Logo string 
// 	SiteUrl string 
// }


var (
	COMPANY_PARENT_ID = 1
)


