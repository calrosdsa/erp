package entity


type Country struct {
	Code string `gorm:"primaryKey"`
	Name string 
}