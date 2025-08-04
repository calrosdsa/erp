package main

import (
	"erp/cmd/gen-annotations/internal"
	"erp/gen/db/model"
	"fmt"
	"log"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gen"

	// "gorm.io/gen/field"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
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
	})
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", viper.GetString("db.host"),
	viper.GetString("db.user"), viper.GetString("db.pass"), viper.GetString("db.name"), viper.GetInt("db.port"))
	gormdb, _ := gorm.Open(postgres.Open(dsn))
	g.UseDB(gormdb) // reuse your gorm db


	g.ApplyInterface(func(internal.Querier) {},
	//admin
	model.RoleTemplate{},
	//crm
	model.Deal{},
	model.DealParticipant{},
	//core
	model.Stage{},
	model.Activity{},
	model.ActivityDeadline{},
	model.ActivityComment{},
	model.ActivityMention{},
	model.Mention{},
	model.CompanyEntity{},
	model.Entity{},
	model.Action{},
	model.CurrencyExchange{},
	model.Workspace{},
	model.WorkspaceModule{},
	model.Module{},
	model.ModuleSection{},
	model.Notification{},
	model.NotificationMention{},
	model.Connection{},
	//acounting
	model.Ledger{},
	model.LedgerAccount{},
	model.CashOutflow{},
	model.Payment{},
	model.PaymentReference{},
	model.JournalEntry{},
	model.JournalEntryLine{},
	model.TransactionLedger{},
	model.Tax{},
	model.AccountSetting{},
	model.CostCenter{},
	model.ChargesTemplate{},
	model.Bank{},
	model.BankAccount{},
	//Project 
	model.Project{},
	//Invoicing
	model.PurchaseRecord{},
	model.SalesRecord{},
	//Chat
	model.Chat{},
	model.ChatMember{},
	model.ChatMessage{},

	//Pricing
	model.PricingCharge{},
	model.Pricing{},
	model.PricingLineItem{},

	model.UnitOfMeasure{},
	model.UnitOfMeasureTranslation{},
	model.Currency{},
	model.Party{},
	model.PartyReference{},
	model.Address{},
	model.Contact{},
	//Company
	model.Company{},
	model.CompanyDefault{},

	model.KeyValue{},
	model.User{},
	model.Role{},
	model.RoleAction{},
	model.UserRelation{},
	model.Profile{},
	model.Group{},
	model.Supplier{},

	//Documents
	model.ProgressOrder{},
	model.Order{},
	model.PaymentTerm{},
	model.PaymentTermsTemplate{},
	model.PaymentTermsLine{},
	model.TermsAndCondition{},
	model.AddressAndContact{},
	model.DocAccount{},
	model.DocTerm{},

	// model.OrderLine{},
	model.ItemLine{},
	model.TaxAndChargeLine{},
	model.ItemLineReceipt{},
	model.DeliveryLineItem{},
	model.ItemLineStockEntry{},
	model.InvoicedItemLine{},
	model.Invoice{},
	model.ProgressInvoice{},
	model.Quotation{},
	model.Receipt{},

	model.PriceList{},
	model.Customer{},

	model.WareHouse{},
	model.Item{},
	model.ItemInventorySetting{},
	model.ItemAttribute{},
	model.ItemAttributeValue{},
	model.ItemVariant{},
	model.StockLevel{},
	model.StockMovement{},
	model.ItemPrice{},
	model.StockSetting{},
	model.StockDefault{},
	model.StockTransaction{},
	model.StockEntry{},
	model.SerialNoReference{},
	model.SerialNo{},
	model.SerialNoTransaction{},
	model.BatchBundle{},

	//Regate
	model.Court{},
	model.CourtRate{},
	model.Booking{},
	model.BookingSlot{},
	model.BookingPrice{},
	model.BookingEvent{},
	model.EventBooking{},
	//Piano
	model.PianoForm{},
)


	g.Execute()
}


func getInt(e string) int {
	val, err := strconv.Atoi(e)
	if err != nil {
		log.Fatal(err)
	}
	return val
}
