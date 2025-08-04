package connection

import (
	"context"
	"erp/gen/db/query"
	"erp/internal/app/config"
	"erp/internal/app/service/helpers"
	_logger "erp/pkg/logger"
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

type Connection struct {
	Db *gorm.DB
	// Adapter   *Adapter
	convertor helpers.ConvertorHelper
	generator helpers.Generator
	logger    _logger.Logger
	Q         *query.Query
}

func NewDbConnection(
	config *config.ConfigModule,
	helpers *helpers.Helpers,
	logger *_logger.Logger,
) *Connection {
	dbName := viper.GetString("db.name")
	dbPass := viper.GetString("db.pass")
	dbUser := viper.GetString("db.user")
	dbPort := viper.GetInt("db.port")
	dbHost := viper.GetString("db.host")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", dbHost,
		dbUser, dbPass, dbName, dbPort)
	fmt.Println("DB NAME MODULE", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	// adapter, err := NewAdapterByDB(db)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	Q := query.Use(db)
	// Q.Profile.First()
	// migrateEntitiesSchema(db, dbConfig.CustomEntities...)
	//Migrate the schemas
	return &Connection{
		Db: db,
		// Adapter:   adapter,
		convertor: helpers.Convertor,
		generator: helpers.Generator,
		logger:    *logger,
		Q:         Q,
	}
}

// func (c *Connection)

// func (c *Connection) GetActions(db *gorm.DB,entityID int,entities interface{})(any,error){
// 	err := db.Where(&entity.Action{EntityID: uint(entityID)}).Find(&entities).Error
// 	return err
// }

func (c *Connection) Order(params map[string]string) func(db *gorm.DB) *gorm.DB {
	fmt.Println(params)
	return func(db *gorm.DB) *gorm.DB {
		order := params["order"]
		column := params["column"]
		fmt.Println(order, column)
		if order != "" && column != "" {
			return db.Order(fmt.Sprintf("%s %s", column, order))
		}
		return db
	}
}

func (c *Connection) PaginateDao(params map[string]string) func(d gen.Dao) gen.Dao {
	return func(d gen.Dao) gen.Dao {
		page := c.convertor.StrtoInt(params["page"])
		size := c.convertor.StrtoInt(params["size"])
		// order := params["order"]
		// column := params["column"]
		// page = page+1
		// offset := (page - 1) * size
		offset := page * size
		// if order != "" && column != "" {
		// 	d = d.Order(fmt.Sprintf("%s %s", column, order))
		// }
		d = d.Offset(int(offset)).Limit(int(size))
		return d
	}
}

func (c *Connection) Paginate(params map[string]string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := c.convertor.StrtoInt(params["page"])
		size := c.convertor.StrtoInt(params["size"])
		order := params["order"]
		column := params["column"]
		// page = page+1
		// offset := (page - 1) * size
		offset := page * size
		b := db
		if order != "" && column != "" {
			b = db.Order(fmt.Sprintf("%s %s", column, order))
		}
		b = db.Offset(int(offset)).Limit(int(size))
		return b
	}
}

// func (c *Connection) GetEntity(ctx context.Context, db *gorm.DB, entity interface{}, entityId uint, companyId uint) error {
// 	if err := db.WithContext(ctx).First(entity, entityId).Where("company_id = ?", companyId).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

func (c *Connection) GetEntity(ctx context.Context, db *gorm.DB, entity interface{}, args ...interface{}) error {
	var preload func(db *gorm.DB) *gorm.DB
	if len(args) == 1 {
		v, ok := args[0].(func(db *gorm.DB) *gorm.DB)
		if ok {
			preload = v
		}
	}
	builder := db.WithContext(ctx).Where(entity)

	if preload != nil {
		builder.Scopes(preload)
	}

	err := builder.First(entity).Error
	if err != nil {
		return err
	}
	return err
}

// func (c *Connection) GetEntity(ctx context.Context, db *gorm.DB, entity interface{}) error {
// 	if err := db.WithContext(ctx).Where(entity).First(entity).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

func (c *Connection) PreloadColumns(columns []string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(columns)
	}
}

func (s *Connection) GenerateCode(ctx context.Context, db *gorm.DB, entity interface{}, companyId int64) (code string) {
	for {
		code = s.generator.GenerateCode()
		// Check if the code already exists in the database
		var count int64
		err := db.Model(entity).
			Where("company_id = ? AND code = ?", companyId, code).
			Count(&count).Error
		if err != nil {
			s.logger.LogError(
				err,
				_logger.OptionsLog.WithMethod("GenerateCode"),
				_logger.OptionsLog.WithOperation("Connection"),
			)
		}

		// If the code is unique, break the loop
		if count == 0 {
			break
		}
	}

	return
}

// func GetEntity[T any](ctx context.Context, db *gorm.DB, entity *T) error {
//     if err := db.WithContext(ctx).First(entity).Error; err != nil {
//         return err
//     }
//     return nil
// }
