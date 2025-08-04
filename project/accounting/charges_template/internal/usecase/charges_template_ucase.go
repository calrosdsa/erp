package charge_template_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/internal/domain/repository"
	"erp/pkg/bus"
	"erp/pkg/di"
	"erp/pkg/logger"
	charges_template_repo "erp/project/accounting/charges_template/internal/repository"
	"fmt"
)

type ChargesTemplateUseCase interface {
	GetChargesTemplate(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.ChargesTemplateDto], err error)
	CreateChargesTemplate(req *common.RequestContext, d *dto.CreateChargesTemplateRequest) (
		res dto.ChargesTemplateDto, err error)
	GetChargesTemplates(req *common.RequestContext, d *dto.RequestPaginationData) (
		res dto.PaginationResult[[]dto.ChargesTemplateDto], err error)
	EditChargesTemplate(req *common.RequestContext, d *dto.EditChargesTemplateRequest) (err error)	
}

type chargesTemplateUcase struct {
	emitLog            logger.EmitLog
	chargeTemplateRepo charges_template_repo.ChargesTemplateRepository
	permission         repository.PermissionService
	core               repository.CoreService
	bus                bus.Bus
	c                  di.Container
}

func NewChargesTemplateUcase(
	logger logger.Logger,
	chargeTemplateRepo charges_template_repo.ChargesTemplateRepository,
	permission repository.PermissionService,
	core repository.CoreService,
	bus bus.Bus,
	c di.Container,
) ChargesTemplateUseCase {
	return &chargesTemplateUcase{
		emitLog:            logger.EmitLog("charges-template-usecase"),
		chargeTemplateRepo: chargeTemplateRepo,
		permission:         permission,
		core:               core,
		bus:                bus,
		c:                  c,
	}
}

func(u *chargesTemplateUcase) EditChargesTemplate(req *common.RequestContext, d *dto.EditChargesTemplateRequest) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditChargesTemplate"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx,req,domain.CHARGES_TEMPLATE,domain.EDIT)
	if err != nil {
		return
	}
	err = u.chargeTemplateRepo.EditChargesTemplate(req,d)
	return 
}

func (u *chargesTemplateUcase) GetChargesTemplate(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.ChargesTemplateDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetChargesTemplate"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.CHARGES_TEMPLATE, domain.VIEW)
	if err != nil {
		return
	}
	res, err = u.chargeTemplateRepo.GetChargesTemplate(req, d)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req,res.Entity.ID)
	
	return
}
func (u *chargesTemplateUcase) CreateChargesTemplate(req *common.RequestContext, d *dto.CreateChargesTemplateRequest) (
	res dto.ChargesTemplateDto, err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateChargesTemplate"))
		}
		err = u.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err)
	}(tx)
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.CHARGES_TEMPLATE, domain.CREATE)
	if err != nil {
		return
	}
	chargeTemplate, err := u.chargeTemplateRepo.CreateChargesTemplate(tx,req, d)
	if err != nil {
		return
	}
	res = dto.ChargesTemplateFromModel(&chargeTemplate)
	err = u.bus.Emit(req.Ctx, domain.ChargesTemplateCreatedEvent, event.ChargesTemplateEventData{
		Tx:                 tx,
		ChargeTemplateData: d.Body,
		ChargesTemplate:    chargeTemplate,
	})

	return
}
func (u *chargesTemplateUcase) GetChargesTemplates(req *common.RequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.ChargesTemplateDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetChargesTemplates"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.CHARGES_TEMPLATE, domain.VIEW)
	if err != nil {
		return res, domain.ACTION_NOT_ALLOWED
	}

	res, err = u.chargeTemplateRepo.GetChargesTemplates(req, d)
	if err != nil {
		return
	}
	return
}
func (s *chargesTemplateUcase) closeTx(tx *query.QueryTx, err error) error {
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
