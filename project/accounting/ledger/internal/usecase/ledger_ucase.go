package ledger_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/logger"
	ledger_repo "erp/project/accounting/ledger/internal/repository"
)

type LedgerUseCase interface {
	CreateLedger(req *common.RequestContext, i dto.LedgerData) (dto.LedgerDto, error)
	GetLedgersAccounts(req *common.RequestContext, i dto.LedgersRequest) (
		dto.ResponseDataList[[]dto.LedgerDto], error)
	GetLedgerDetail(req *common.RequestContext, i *dto.RequestEntity) (
		dto.ResultEntity[dto.LedgerDetailDto], error)
	GetGeneralLedgerReport(req *common.RequestContext, i *dto.RequestGeneralLedger) (
		[]dto.GeneralLedgerEntryDto, error)
	GetLedgerAccountsTree(req *common.RequestContext) (
		[]dto.TreeEntryDto, error)
	EditLedger(req *common.RequestContext, d dto.LedgerData) (err error)
}

type ledgerUseCase struct {
	ledgerRepo ledger_repo.LedgerRepository
	permission repository.PermissionService
	emitLog    logger.EmitLog
	core repository.CoreService
}

func NewLedgerUseCase(
	ledgerRepo ledger_repo.LedgerRepository,
	logger logger.Logger,
	permission repository.PermissionService,
	core repository.CoreService,
) LedgerUseCase {
	return &ledgerUseCase{
		ledgerRepo: ledgerRepo,
		permission: permission,
		emitLog:    logger.EmitLog("ledger-usecase"),
		core: core,
	}
}
func (u *ledgerUseCase) EditLedger(req *common.RequestContext, d dto.LedgerData) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditLedger"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.LEDGER, domain.EDIT)
	if err != nil {
		return
	}
	err = u.ledgerRepo.EditLedger(req, d)
	return
}

func (u *ledgerUseCase) GetLedgerAccountsTree(req *common.RequestContext) (
	res []dto.TreeEntryDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetLedgerAccountsTree"))
		}
	}()
	res, err = u.ledgerRepo.GetLedgerAccountsTree(req)
	return
}

func (u *ledgerUseCase) GetGeneralLedgerReport(req *common.RequestContext, i *dto.RequestGeneralLedger) (
	res []dto.GeneralLedgerEntryDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetGeneralLedgerReport"))
		}
	}()
	// if allow := u.permission.
	res, err = u.ledgerRepo.GetGeneralLedgerReport(req, i)

	return res, err
}

func (u *ledgerUseCase) GetLedgerDetail(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[dto.LedgerDetailDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetLedgerDetail"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.LEDGER, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = u.ledgerRepo.GetLedgerDetail(req, i)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req,res.Entity.ID)
	return
}

func (u *ledgerUseCase) CreateLedger(req *common.RequestContext, i dto.LedgerData) (
	res dto.LedgerDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateLedger"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.LEDGER, domain.CREATE); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = u.ledgerRepo.CreateLedger(req, i)
	return
}

func (u *ledgerUseCase) GetLedgersAccounts(req *common.RequestContext, d dto.LedgersRequest) (
	res dto.ResponseDataList[[]dto.LedgerDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetLedgersAccounts"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.LEDGER, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res.Body.Result, err = u.ledgerRepo.GetLedgersAccounts(req, d)

	return res, err
}
