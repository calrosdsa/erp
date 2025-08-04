package accounting

import (
	"context"
	"erp/pkg/system"
	"erp/project/accounting/banking"
	"erp/project/accounting/cash_outflow"
	"erp/project/accounting/charges_template"
	"erp/project/accounting/cost_center"
	"erp/project/accounting/journal"
	"erp/project/accounting/ledger"
	"erp/project/accounting/payment"
	acct_report "erp/project/accounting/report"
	"erp/project/accounting/tax_and_charges"
	"erp/project/accounting/transaction"
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
			&payment.Module{},
			&cash_outflow.Module{},
			&transaction.Module{},
			&ledger.Module{},
			&acct_report.Module{},
			&journal.Module{},
			&cost_center.Module{},
			&tax_and_charges.Module{},
			&charges_template.Module{},
			//Baking
			&banking.Module{},
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
