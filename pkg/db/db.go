package db

import (
	"context"
	"crypto/rand"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/domain"
	"erp/pkg/logger"
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type connection struct {
	Db      *gorm.DB
	Q       *query.Query
	emitLog logger.EmitLog
	logger  logger.Logger
}

func NewDbConnection(
	logger logger.Logger,
) Connection {
	dbName := viper.GetString("db.name")
	dbPass := viper.GetString("db.pass")
	dbUser := viper.GetString("db.user")
	dbPort := viper.GetInt("db.port")
	dbHost := viper.GetString("db.host")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", dbHost,
		dbUser, dbPass, dbName, dbPort)
	fmt.Println("DB", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	Q := query.Use(db)

	return &connection{
		Db:      db,
		logger:  logger,
		Q:       Q,
		emitLog: logger.EmitLog("db-connection"),
	}
}

func (c *connection) GetDbHelper() DbHelper {
	return NewDbHelper(c.logger)
}

func (c *connection) ValidateCode(ctx context.Context, entity interface{}, companyId int64, code string) (err error) {
	var count int64
	err = c.Db.Model(entity).
		Where("company_id = ? AND code = ?", companyId, code).
		Count(&count).Error
	if err != nil {
		c.emitLog.Err(err, logger.OptionsLog.WithMethod("GenerateCode"))
	}
	if count > 0 {
		return domain.ERROR_ITEM_CODE_TAKEN
	}
	return
}

func (c *connection) GenerateCode(ctx context.Context, entity interface{}, companyId int64) (code string) {
	for {
		code = c.generateCode()
		// Check if the code already exists in the database
		var count int64
		err := c.Db.Model(entity).
			Where("company_id = ? AND code = ?", companyId, code).
			Count(&count).Error
		if err != nil {
			c.emitLog.Err(err, logger.OptionsLog.WithMethod("GenerateCode"))
		}
		// If the code is unique, break the loop
		if count == 0 {
			break
		}
	}

	return
}

func (c *connection) generateCode() string {
	n := 4
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%X", b)
}

func (c *connection) InsertReference(ctx context.Context, tx *query.QueryTx, partyID int64, referenceID int64) error {
	err := tx.PartyReference.WithContext(ctx).Save(&model.PartyReference{
		PartyID:     partyID,
		ReferenceID: referenceID,
	})
	return err
}

func (c *connection) GetPartyByUUID(ctx context.Context, companyID int64, partyType string, uuid string) *int64 {
	switch partyType {
	case proto.PartyType_supplier.String():
		supplier, err := c.Q.Supplier.WithContext(ctx).Select(c.Q.Supplier.ID).Where(
			c.Q.Supplier.UUID.Eq(uuid),
			c.Q.Supplier.CompanyID.Eq(companyID),
		).First()
		if err != nil {
			return nil
		}
		return &supplier.ID

	case proto.PartyType_customer.String():
		customer, err := c.Q.Customer.WithContext(ctx).Select(c.Q.Customer.ID).Where(
			c.Q.Customer.UUID.Eq(uuid),
			c.Q.Customer.CompanyID.Eq(companyID),
		).First()
		if err != nil {
			return nil
		}
		return &customer.ID
	case proto.PartyType_purchaseOrder.String():
		order, err := c.Q.Order.WithContext(ctx).Select(c.Q.Order.ID).Where(
			c.Q.Order.Code.Eq(uuid),
			c.Q.Order.CompanyID.Eq(companyID),
		).First()
		if err != nil {
			return nil
		}
		return &order.ID
	}
	return nil
}

func (c *connection) GetDB() *gorm.DB {
	return c.Db
}

func (c *connection) GetQ() *query.Query {
	Q := query.Use(c.Db)
	return Q
}
