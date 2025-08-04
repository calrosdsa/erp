package connection

import (
	"erp/internal/app/entity"
	// "fmt"

	"gorm.io/gorm"
)

func migrateEntitiesSchema(db *gorm.DB, customEntities ...interface{}) {
	// customEntities = append(customEntities, &entity.CasbinRule{})
	// customEntities = append(customEntities, &entity.User{})
	// customEntities = append(customEntities, &entity.Entity{})
	// customEntities = append(customEntities, &entity.Action{})
	// customEntities = append(customEntities, &entity.Role{})
	// customEntities = append(customEntities, &entity.RoleActions{})

	// customEntities = append(customEntities, &entity.UserRelation{})

	// customEntities = append(customEntities, &entity.PartyType{})
	// customEntities = append(customEntities, &entity.Party{})

	// customEntities = append(customEntities, &entity.ClientKeyValueData{})
	// customEntities = append(customEntities, &entity.Client{})
	// customEntities = append(customEntities, &entity.Administrator{})
	// customEntities = append(customEntities, &entity.Company{})

	// customEntities = append(customEntities, &entity.UnitOfMeasure{})
	// customEntities = append(customEntities, &entity.UnitOfMeasureTranslation{})

	// customEntities = append(customEntities, &entity.ItemGroup{})
	// customEntities = append(customEntities, &entity.Item{})
	// customEntities = append(customEntities, &entity.ItemVariant{})
	// customEntities = append(customEntities, &entity.ItemAttribute{})
	// customEntities = append(customEntities, &entity.ItemPriceList{})
	// customEntities = append(customEntities, &entity.ItemAttributeValue{})
	// customEntities = append(customEntities, &entity.ItemPrice{})
	// customEntities = append(customEntities, &entity.ItemPricePlugin{})
	// customEntities = append(customEntities, &entity.WareHouse{})
	// customEntities = append(customEntities, &entity.StockLevel{})
	// customEntities = append(customEntities, &entity.StockMovement{})

	customEntities = append(customEntities, &entity.CompanyPlugins{})
	customEntities = append(customEntities, &entity.Plugin{})

	customEntities = append(customEntities, &entity.SalesOrder{})
	customEntities = append(customEntities, &entity.SalesItemLine{})
	customEntities = append(customEntities, &entity.SalesInvoice{})
	customEntities = append(customEntities, &entity.SalesOrderPlugin{})
	customEntities = append(customEntities, &entity.Address{})

	//BUYING MODULE
	// customEntities = append(customEntities, &entity.Group{})
	// customEntities = append(customEntities, &entity.Supplier{})

	// customEntities = append(customEntities, &entity.Tax{})

	// err := db.AutoMigrate(
	// 	customEntities...,
	// )
	// if err != nil {
	// 	fmt.Println("FAIL TO AUTOMIGRATE User table")
	// }
}
