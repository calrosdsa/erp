package domain

const (
	DEFAULT_LIMIT = 20
	DEFAULT_ACTIVITY_LIMIT = 30
	MAX_LIMIT     = 1000
	DEFAULT_UOM   = 5
	DEFAULT_TZ = "America/La_Paz"
	DEFAULT_BATCH_SIZE = 50
)

// Dependency Injection Keys
const (
	DatabaseTransactionKey = "tx"
	DbKey = "db"
	ItemUseCase = "item_usecase"
	ItemPriceUseCase = "item_price_usecase"
	PriceListUseCase = "price_list_usecase"
	OrderUseCase = "order_usecase"
	QuotationUseCase = "quotation_usecase"
	NotificationUseCase = "notification_usecase"
	ActivityUseCase = "activity_usecase"
	ContactUseCase = "contact_usecase"
	ModuleUseCase = "module_usecase"
)

const (
	DEFAULT_CURRENCY = "BOB"
	DEFAULT_LANGUAGE = "en"
)


type contextKey string

const (
	SESSION_KEY contextKey = "session"
)