package main

import (
	// "erp/cmd/gen/internal"
	// "erp/gen/db/model"
	"erp/cmd/gen-models/internal"
	"fmt"
	"log"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func init(){
	viper.SetConfigFile(`../../configs/config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	g := gen.NewGenerator(gen.Config{
		OutPath:       "../../gen/db/query",
		Mode:          gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
		FieldNullable: true,
		// WithUnitTest: true,
	})
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", viper.GetString("db.host"),
	viper.GetString("db.user"), viper.GetString("db.pass"), viper.GetString("db.name"), viper.GetInt("db.port"))
	gormdb, _ := gorm.Open(postgres.Open(dsn))
	g.UseDB(gormdb) // reuse your gorm db

	var models []interface{}

	adminModels := internal.AdminModels(g)
	accountingModels := internal.AccountingModels(g)
	companyModels := internal.CompanyModels(g)
	buyingModels := internal.BuyingModels(g)
	documentModels := internal.DocumentModels(g)
	stockModels := internal.StockModels(g)
	coreModels := internal.CoreModels(g)
	regateModels := internal.RegateModels(g)
	partyModels := internal.PartModels(g)
	projectModels := internal.ProjectModels(g)
	pianoForm := internal.PianoModels(g)
	sellingModels := internal.SellingModels(g)
	invoicingModels:= internal.InvoicingModels(g)
	pricingModels := internal.PricingModels(g)
	crmModels := internal.CrmModels(g)
	chatModels := internal.ChatModels(g)
	authModels := internal.AuthModels(g)
	models = append(models, authModels...)
	models = append(models, invoicingModels...)
	models= append(models, companyModels...)
	models = append(models, accountingModels...)
	models = append(models, buyingModels...)
	models = append(models, stockModels...)
	models = append(models, regateModels...)
	models = append(models, coreModels...)
	models = append(models, documentModels...)
	models = append(models, partyModels...)
	models = append(models, pianoForm...)
	models = append(models, adminModels...)
	models = append(models, projectModels...)
	models = append(models, sellingModels...)
	models = append(models, pricingModels...)
	models = append(models, crmModels...)
	models = append(models, chatModels...)
	all := false
	if all {
		g.ApplyBasic(
			models...,
		)
	}else {
		g.ApplyBasic(projectModels...)
	}
	// g.ApplyInterface(func(internal.Querier) {},model.Group{}, model.Supplier{})

	// g.ApplyBasic(
	// 	// Generate structs from all tables of current database
	// 	g.GenerateAllTable()...,
	// )
	// Generate the code
	g.Execute()
}

func getInt(e string) int {
	val, err := strconv.Atoi(e)
	if err != nil {
		log.Fatal(err)
	}
	return val
}
