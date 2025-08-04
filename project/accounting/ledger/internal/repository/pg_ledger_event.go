package ledger_repo

import (
	"context"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/event"
	pkg_ledger "erp/project/accounting/ledger/internal/pkg"
)

type LedgerEventRepo interface {
	CreateChartOfAccountsCompany(ctx context.Context, payload event.CreatedCompanyEventData) (err error)
}

type ledgerEventRepo struct {
	locale         helpers.Locale
	accountBuilder *pkg_ledger.AccountBuilder
}

func NewLedgerEventRepo(
	helpers *helpers.Helpers,
	accountBuilder *pkg_ledger.AccountBuilder,
) LedgerEventRepo {
	return &ledgerEventRepo{
		locale:         helpers.Locale,
		accountBuilder: accountBuilder,
	}
}

func (r *ledgerEventRepo) CreateChartOfAccountsCompany(ctx context.Context,
	payload event.CreatedCompanyEventData) (err error) {
	tx := payload.Tx
	err = r.accountBuilder.CreateChartofAccount(ctx,tx,payload)
	return
}
