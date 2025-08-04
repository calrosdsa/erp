package tac_usecase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/pkg/logger"
	tac_repo "erp/project/accounting/tax_and_charges/repository"
)

type TacUseCase interface {
	GetTACLines(req *common.RequestContext, d *dto.RequestTaxLines) (
		res []dto.TaxAndChargeLineDto, err error)
	EditTaxAndChargeLine(req *common.RequestContext,
		d *dto.EditTaxLineRequest) (err error)
	CreateTaxAndChargeLine(req *common.RequestContext,
		d *dto.AddTaxLineRequest) (err error)
	DeleteTaxAndChargeLine(req *common.RequestContext,d *dto.DeleteTaxLineRequest)(err error)
}

type tacUseCase struct {
	emitLog logger.EmitLog
	tacRepo tac_repo.TacRepository
}

func NewTacUseCase(
	logger logger.Logger,
	tacRepo tac_repo.TacRepository,
) TacUseCase {
	return &tacUseCase{
		emitLog: logger.EmitLog("taxes-and-charges"),
		tacRepo: tacRepo,
	}
}

func(u *tacUseCase) DeleteTaxAndChargeLine(req *common.RequestContext,d *dto.DeleteTaxLineRequest)(err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("DeleteTaxAndChargeLine"))
		}
	}()
	err = u.tacRepo.DeleteTaxAndChargeLine(req, d)
	return 
}

func (u *tacUseCase) GetTACLines(req *common.RequestContext, d *dto.RequestTaxLines) (
	res []dto.TaxAndChargeLineDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetTACLines"))
		}
	}()
	res, err = u.tacRepo.GetTACLines(req, d)
	return
}

func (u *tacUseCase) EditTaxAndChargeLine(req *common.RequestContext,
	d *dto.EditTaxLineRequest) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditTaxAndChargeLine"))
		}
	}()
	err = u.tacRepo.EditTaxAndChargeLine(req, d)
	return

}
func (u *tacUseCase) CreateTaxAndChargeLine(req *common.RequestContext,
	d *dto.AddTaxLineRequest) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateTaxAndChargeLine"))
		}
	}()
	err = u.tacRepo.CreateTaxAndChargeLine(req, d)
	return

}
