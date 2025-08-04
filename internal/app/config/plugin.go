package config

import (
	clientconfig "erp/internal/app/config/client_config"
	sellingconfig "erp/internal/app/config/selling_config"
	stockconfig "erp/internal/app/config/stock_config"
)

type PluginApp struct {
	Name string
	Strategies 
}

type Strategies struct {
	ItemStrategy stockconfig.ItemStrategy `json:"-"`
	ClientStrategy clientconfig.ClientStrategy `json:"-"`
	SalesOrderStrategy sellingconfig.SalesOrderStrategy `json:"-"`
	ItemPriceStrategy stockconfig.ItemPriceStrategy `json:"-"`
}

const (
	PLUGIN_SQUARE = "square"
	PLUGIN_EMAIL = "email"
)
