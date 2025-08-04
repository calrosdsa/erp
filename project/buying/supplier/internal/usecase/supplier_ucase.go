package supplier_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/internal/domain/repository"
	"erp/pkg/bus"
	"erp/pkg/di"
	"erp/pkg/fsm"
	"erp/pkg/logger"
	supplier_repo "erp/project/buying/supplier/internal/repository"
)

type SupplierUseCase interface {
	GetSupplier(req *common.RequestContext, i *dto.RequestEntity) (
		dto.ResultEntity[dto.SupplierDto], error)
	GetSuppliers(req *common.RequestContext, i *dto.RequestPaginationData) (
		dto.PaginationResult[[]dto.SupplierDto], error)
	CreateSupplier(req *common.RequestContext, i dto.SupplierData) (dto.SupplierDto, error)
	EditSupplier(req *common.RequestContext, d dto.SupplierData) (err error)
	UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent) error
}

type supplierUseCase struct {
	emitLog      logger.EmitLog
	permission   repository.PermissionService
	supplierRepo supplier_repo.SupplierRepository
	core         repository.CoreService
	bus          bus.Bus
	c            di.Container
	fsm          fsm.FsmState
}

func NewSupplierUseCase(
	permission repository.PermissionService,
	supplierRepo supplier_repo.SupplierRepository,
	logger logger.Logger,
	core repository.CoreService,
	bus bus.Bus,
	c di.Container,
	fsm fsm.FsmState,
) SupplierUseCase {
	return &supplierUseCase{
		emitLog:      logger.EmitLog("supplier-service"),
		permission:   permission,
		supplierRepo: supplierRepo,
		core:         core,
		c:            c,
		bus:          bus,
		fsm:          fsm,
	}
}

func (u *supplierUseCase) EditSupplier(req *common.RequestContext, d dto.SupplierData) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditSupplier"))
		}
		err = u.closeTx(tx, err)
	}(tx)
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.SUPPLIER, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	err = u.supplierRepo.EditSupplier(tx, req, d)
	if err != nil {
		return
	}

	err = u.bus.Emit(ctx, domain.SupplierEditedEvent, event.SupplierEventData{
		Req:  *req,
		Tx:   tx,
		Data: d,
	})
	return err
}
func (u *supplierUseCase) UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdateStatus"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.SUPPLIER, domain.EDIT); err != nil {
		return
	}
	nextState, err := u.fsm.NextState(d.Body.CurrentState, d.Body.Events)
	if err != nil {
		return err
	}
	err = u.supplierRepo.UpdateStatus(req, d, nextState)
	return
}

func (s *supplierUseCase) GetSupplier(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[dto.SupplierDto], err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetSupplier"))
		}
	}()
	if allow := s.permission.CheckPermission(req.Ctx, req, domain.SUPPLIER, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = s.supplierRepo.GetSupplier(req, i)
	if err != nil {
		return
	}
	res.Activities = s.core.GerActivitiesByPartyID(req, res.Entity.ID)
	res.Addresses = s.core.GetPartyAddresses(req, res.Entity.ID)
	res.Contacts = s.core.GetPartyContacts(req, res.Entity.ID)
	return res, err
}

func (s *supplierUseCase) GetSuppliers(req *common.RequestContext, i *dto.RequestPaginationData) (
	dto.PaginationResult[[]dto.SupplierDto], error) {
	var (
		err error
	)
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetSuppliers"))
		}
	}()

	var response dto.PaginationResult[[]dto.SupplierDto]
	if allow := s.permission.CheckPermission(req.Ctx, req, domain.SUPPLIER, domain.VIEW); !allow {
		return response, domain.ACTION_NOT_ALLOWED
	}
	response, err = s.supplierRepo.GetSuppliers(req, i)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (u *supplierUseCase) CreateSupplier(req *common.RequestContext, i dto.SupplierData) (res dto.SupplierDto, err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateSupplier"))
		}
		err = u.closeTx(tx, err)
	}()
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.SUPPLIER, domain.CREATE); err != nil {
		return res, err
	}
	res, err = u.supplierRepo.CreateSupplier(req, tx, i)
	if err != nil {
		return
	}
	err = u.bus.Emit(ctx, domain.SupplierCreatedEvent, event.SupplierEventData{
		Req:      *req,
		Tx:       tx,
		Data:     i,
		Supplier: res,
	})
	return
}

func (s *supplierUseCase) closeTx(tx *query.QueryTx, err error) error {
	if p := recover(); p != nil {
		_ = tx.Rollback()
		panic(p)
	} else if err != nil {
		_ = tx.Rollback()
		return err
	} else {
		return tx.Commit()
	}
}
