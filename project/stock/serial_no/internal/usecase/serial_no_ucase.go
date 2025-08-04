package serial_no_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/logger"
	serial_no_repo "erp/project/stock/serial_no/internal/repository"
)

type SerialNoUseCase interface {
	GetSerialNo(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.SerialNoDto], err error)
	GetSerialNos(req *common.RequestContext, d *dto.RequestSerialNos) (
		res dto.PaginationResult[[]dto.SerialNoDto], err error)
	GetSerialNoTransactions(req *common.RequestContext, d *dto.RequestSerialNoTransactions) (
		res []dto.SerialNoTransactionDto, err error)
}

type serialNoUcase struct {
	emitLog      logger.EmitLog
	serialNoRepo serial_no_repo.SerialNoRepository
	permission   repository.PermissionService
	core         repository.CoreService
}

func NewSerialUcase(
	logger logger.Logger,
	serialNoRepo serial_no_repo.SerialNoRepository,
	permission repository.PermissionService,
	core repository.CoreService,
) SerialNoUseCase {
	return &serialNoUcase{
		emitLog:      logger.EmitLog("serial-no-usecase"),
		serialNoRepo: serialNoRepo,
		permission:   permission,
		core:         core,
	}
}

func (u *serialNoUcase) GetSerialNoTransactions(req *common.RequestContext, d *dto.RequestSerialNoTransactions) (
	res []dto.SerialNoTransactionDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetSerialNoTransactions"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.SERIAL_NO, domain.VIEW)
	if err != nil {
		return
	}
	res, err = u.serialNoRepo.GetSerialNoTransactions(req, d)
	return
}

func (u *serialNoUcase) GetSerialNo(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.SerialNoDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetSerialNo"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.SERIAL_NO, domain.VIEW)
	if err != nil {
		return
	}
	res, err = u.serialNoRepo.GetSerialNo(req, d)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req,res.Entity.ID)
	return
}

func (u *serialNoUcase) GetSerialNos(req *common.RequestContext, d *dto.RequestSerialNos) (
	res dto.PaginationResult[[]dto.SerialNoDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetSerialNos"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.SERIAL_NO, domain.VIEW)
	if err != nil {
		return
	}
	res, err = u.serialNoRepo.GetSerialNos(req, d)
	if err != nil {
		return
	}
	return
}
