package helpers

import (
	"erp/pkg/exporter"
	"erp/pkg/logger"

	"erp/pkg/config"

	"github.com/asaskevich/EventBus"
)

type Helpers struct {
	Convertor ConvertorHelper
	Session   SessionHelper
	Validator *ValidatorHelper
	Locale    Locale
	Currency  CurrencyHelper
	Tax       *TaxHelper
	Generator Generator
	Event     *EventHelper
	Logger    *LoggerHelper
	Permify   *PermifyHelper
	Cache     *CacheHelper
	Error ErrorHelper
	Jwt JwtHelper
	ExcelExporter exporter.ExcelExporter
	PdfExporter exporter.PdfExporter
	Util Util
	Query QueryHelper
}

func Init(
	logger logger.Logger,
	bus *EventBus.Bus,
	appConfig *config.AppConfig,
) *Helpers {
	loggerHelper := NewLoggerHelper(logger)
	locale := NewLocaleHelper()
	return &Helpers{
		Convertor: NewConvertorHelper(),
		Session:   NewSessionHelper(loggerHelper),
		Validator: NewValidator(),
		Locale:    locale,
		Currency:  NewCurrencyHelper(),
		Jwt: NewJwtHelper(appConfig),
		Tax:       NewTaxHelper(),
		Generator: NewGeneratorHelper(logger),
		Event:     NewEventBus(bus),
		Logger:    loggerHelper,
		Cache:     NewCacheHelper(),
		Error: NewErrorHelper(locale),
	}
}


func InitHelpers(
	logger logger.Logger,
	appConfig *config.AppConfig,
)*Helpers {
	loggerHelper := NewLoggerHelper(logger)
	locale := NewLocaleHelper()
	convertor := NewConvertorHelper()
	return &Helpers{
		Convertor: convertor,
		Session:   NewSessionHelper(loggerHelper),
		Validator: NewValidator(),
		Locale:    locale,
		Currency:  NewCurrencyHelper(),
		Tax:       NewTaxHelper(),
		Generator: NewGeneratorHelper(logger),
		Logger:    loggerHelper,
		Cache:     NewCacheHelper(),
		Error: NewErrorHelper(locale),
		Jwt: NewJwtHelper(appConfig),
		ExcelExporter: exporter.NewExcelExporter(),
		PdfExporter: exporter.NewPdfExporter(),
		Util: NewUtil(),
		Query:NewQueryHelpers(convertor),
	}
}