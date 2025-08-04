package services

import (
	"erp/internal/app/config"
	"erp/internal/app/connection"
	"erp/internal/app/domain/repository"
	"erp/internal/app/plugin"
	"erp/internal/app/service/helpers"
	"erp/internal/app/service/services/account_service"
	accountingservice "erp/internal/app/service/services/accounting_service"
	// addressservice "erp/internal/app/service/services/address_service"
	appservice "erp/internal/app/service/services/app_service"
	"erp/internal/app/service/services/buying"
	_clientservice "erp/internal/app/service/services/client_service"
	"erp/internal/app/service/services/company_service"
	// groupservice "erp/internal/app/service/services/group_service"

	// integrationservice "erp/internal/app/service/services/integration_service"
	"erp/internal/app/service/services/jwt_service"
	partyservice "erp/internal/app/service/services/party_service"
	pluginservice "erp/internal/app/service/services/plugin_service"

	// _sellingservice "erp/internal/app/service/services/selling_service"
	stockservice "erp/internal/app/service/services/stock_service"
	uomservice "erp/internal/app/service/services/uom_service"
	_userservice "erp/internal/app/service/services/user_service"
	"erp/pkg/cache"
	"erp/pkg/db"
	_logger "erp/pkg/logger"
	"erp/pkg/permission"
)

type Services struct {
	DomainService *appservice.DomainService
	// GroupService  *groupservice.GroupService
	OrderService  repository.OrderService

	PermissionService permission.PermissionService
	UserService       *_userservice.UserService
	ProfileService    *_userservice.ProfileService
	SessionService    *account_service.SessionService
	// Partyservice      *partyservice.PartyService
	PartyServices     *repository.PartyServices
	JwtService        *jwt_service.JwtService
	AccountService    *account_service.AccountService
	RoleService       *account_service.RoleService
	CompanyService    *company_service.CompanyService
	ClientService     *_clientservice.ClientService
	PluginService     *pluginservice.PluginService

	// AddressService *addressservice.AddressService

	// ItemGroupService     *stockservice.ItemGroupService
	ItemVariantService   *stockservice.ItemVariantService
	ItemAttributeService *stockservice.ItemAttributeService
	ItemStockService     *stockservice.ItemStockService
	PriceListService     *stockservice.PriceListService
	// WareHouseService     *warehouseservice.WareHouseService

	UOMService *uomservice.UOMService

	InvoiceService repository.InvoiceService

	// SalesOrderService *_sellingservice.SalesOrderService
	// SalesOrderInvoice *_sellingservice.SalesInvoiceService

	TaxService *accountingservice.TaxService

	// TecluMobilityService *integrationservice.TecluMobilityService
	Buying  *repository.BuyingServices
}

func Init(
	connection *connection.Connection,
	db db.Connection,
	configModule *config.ConfigModule,
	helpers *helpers.Helpers,
	logger _logger.Logger,
	plugins *plugin.PluginModule,
	repositories *repository.Repositories,
) *Services {
	cache := cache.NewCache(connection)
	configService := configModule.ConfigService
	timeout := configService.GetApiOptions().TimeoutAPICall

	domainService := appservice.NewDomainService(connection, timeout, helpers)
	permissionService := permission.NewPermissionService(db,logger)

	userService := _userservice.NewUserService(connection, timeout, configService, helpers)
	profileService := _userservice.NewProfileService(connection, timeout, helpers,
		permissionService, userService,logger,
	)
	// groupService := groupservice.NewGroupService(connection, timeout, helpers, permissionService,logger)
	//Company
	companyService := company_service.NewCompanyService(connection, &timeout, helpers, userService, permissionService,logger)
	// roleService := account_service.NewRoleService(connection, timeout, helpers, companyService, permissionService,logger)
	userRelationService := account_service.NewSessionService(connection, timeout, helpers)

	// partyService := partyservice.NewPartyService(connection, timeout, helpers)
	partyServices := partyservice.NewPartyServices(connection, timeout, helpers, repositories, permissionService,logger)

	clientService := _clientservice.NewClientService(connection,
		&timeout, plugins, configService, userService, helpers,
	)

	// addressService := addressservice.NewAddressService(connection, timeout, helpers)

	
	itemAttribyteService := stockservice.NewItemAttributeService(connection, timeout, helpers, permissionService,logger)
	itemVariantService := stockservice.NewItemVaraintService(connection, timeout, helpers, permissionService, cache,logger)
	itemStockService := stockservice.NewItemStockService(connection, timeout, helpers, cache, permissionService,logger)
	// salesOrderService := _sellingservice.NewSalesOrderService(connection, &timeout, helpers, plugins)
	// salesInvoiceService := _sellingservice.NewSalesInvoiceService(connection, &timeout, salesOrderService)

	//Accounting
	taxService := accountingservice.NewTaxService(connection, timeout, helpers, permissionService,logger)

	//INTEGRATIONS
	// teclumobilityService := integrationservice.NewTecluMobilityService(connection, timeout, helpers)

	buying := buying.NewBuyingServices(connection, timeout, helpers, permissionService,repositories,logger)
	return &Services{
		DomainService:     domainService,
		// GroupService:      groupService,
		PermissionService: permissionService,
		UserService:       userService,
		ProfileService:    profileService,
		SessionService:    userRelationService,
		// Partyservice:      partyService,
		PartyServices:     partyServices,
		JwtService:        jwt_service.NewJwtService(configService),
		AccountService:    account_service.NewAccountService(&timeout, connection, configService, helpers),
		// RoleService:       roleService,
		CompanyService:    companyService,

		// AddressService: addressService,
		PluginService: pluginservice.NewPluginService(configService, &timeout, connection),
		// ItemGroupService:     stockservice.NewItemGroupService(connection, &timeout, helpers),
		ItemVariantService:   itemVariantService,
		ItemAttributeService: itemAttribyteService,
		ItemStockService:     itemStockService,
		// WareHouseService:     wareHouseService,
		PriceListService:     stockservice.NewPriceListServer(connection, timeout, helpers, permissionService,logger),
		UOMService:           uomservice.NewUOMService(connection, timeout),
		ClientService:        clientService,
		// SalesOrderService:    salesOrderService,
		// SalesOrderInvoice:    salesInvoiceService,
		// ErrorService: &errorService,
		// TecluMobilityService: teclumobilityService,

		TaxService: taxService,
		Buying:     buying,
	}
}
