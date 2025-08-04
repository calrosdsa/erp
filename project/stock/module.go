package stock

import (
	"context"
	"erp/pkg/system"
	"erp/project/stock/batch_bundle"
	"erp/project/stock/item"
	"erp/project/stock/item_inventory"
	"erp/project/stock/itemline"
	"erp/project/stock/itemprice"
	pricelist "erp/project/stock/price_list"
	"erp/project/stock/serial_no"
	"erp/project/stock/stock_entry"
	stocke_ledger "erp/project/stock/stock_ledger"
	"erp/project/stock/warehouse"
)

type Module struct{}

type monolith struct {
	system.Service
	modules []system.Module
}


func (m Module) Startup(ctx context.Context,svc system.Service) error {
	return Root(ctx,svc)
}
func Root(ctx context.Context,svc system.Service) error{
	m := monolith{
		Service: svc,
		modules: []system.Module{
			&serial_no.Module{},
			&batch_bundle.Module{},
			&itemline.Module{},
			&itemprice.Module{},
			&item.Module{},
			&item_inventory.Module{},
			&stocke_ledger.Module{},
			&warehouse.Module{},
			&stock_entry.Module{},
			&pricelist.Module{},
		},
	}

	return m.startupModules()
}


func (m *monolith) startupModules() error {

	for _, module := range m.modules {
		ctx := m.Waiter().Context()
		if err := module.Startup(ctx, m); err != nil {
			return err
		}
	}
	return nil
}
