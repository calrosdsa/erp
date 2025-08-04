package workspace_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/di"
	"erp/pkg/fsm"
	"erp/pkg/logger"
	module_ucase "erp/project/core/module/usecase"
	workspace_repo "erp/project/core/workspace/repository"
	"strconv"
)

type WorkSpaceUseCase interface {
	EditWorkSpace(req *common.RequestContext, d dto.WorkSpaceData) (err error)
	CreateWorkSpace(req *common.RequestContext, d dto.WorkSpaceData) (res dto.WorkSpaceDto, err error)
	GetWorkSpaces(req *common.RequestContext, d dto.WorkSpaceRequest) (res dto.ResponseDataList[[]dto.WorkSpaceDto], err error)
	GetWorkSpace(req *common.RequestContext, d dto.RequestEntity) (res dto.ResultEntity[dto.WorkSpaceDto], err error)
	UpdateStatus(req *common.RequestContext, d dto.UpdateStatusWithEvent) (err error)
}

type workSpaceUseCase struct {
	permission    repository.PermissionService
	core          repository.CoreService
	repo          workspace_repo.WorkSpaceRepository
	emitLog       logger.EmitLog
	fsm           fsm.FsmState
	c             di.Container
	moduleUseCase module_ucase.ModuleUsecase
}

func NewWorkSpaceUseCase(
	permission repository.PermissionService,
	core repository.CoreService,
	repo workspace_repo.WorkSpaceRepository,
	logger logger.Logger,
	c di.Container,
	fsm fsm.FsmState,
) WorkSpaceUseCase {
	return &workSpaceUseCase{
		permission:    permission,
		core:          core,
		repo:          repo,
		emitLog:       logger.EmitLog("workspace-ucase"),
		c:             c,
		fsm:           fsm,
		moduleUseCase: c.Get(domain.ModuleUseCase).(module_ucase.ModuleUsecase),
	}
}

func (u *workSpaceUseCase) UpdateStatus(req *common.RequestContext, d dto.UpdateStatusWithEvent) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdateStatus"))
		}
		err = domain.CloseTx(tx, err)
	}()
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.WORKSPACE, domain.EDIT); err != nil {
		return domain.ACTION_NOT_ALLOWED
	}
	nextState, err := u.fsm.NextState(d.Body.CurrentState, d.Body.Events)
	if err != nil {
		return err
	}
	_, err = u.repo.UpdateStatus(tx, req, d, nextState)
	if err != nil {
		return
	}
	return
}

func (u *workSpaceUseCase) EditWorkSpace(req *common.RequestContext, d dto.WorkSpaceData) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DbKey).(*query.Query).Begin()

	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditWorkSpace"))
		}
		err = domain.CloseTx(tx, err)
	}(tx)
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.WORKSPACE, domain.CREATE); err != nil {
		return err
	}
	err = u.repo.EditWorkSpace(tx, req, d)
	return
}
func (u *workSpaceUseCase) CreateWorkSpace(req *common.RequestContext, d dto.WorkSpaceData) (res dto.WorkSpaceDto, err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DbKey).(*query.Query).Begin()

	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateWorkSpace"))
		}
		err = domain.CloseTx(tx, err)
	}(tx)
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.WORKSPACE, domain.CREATE); err != nil {
		return res, err
	}
	res, err = u.repo.CreateWorkSpace(tx, req, d)
	if err != nil {
		return
	}

	return
}
func (u *workSpaceUseCase) GetWorkSpaces(req *common.RequestContext, d dto.WorkSpaceRequest) (res dto.ResponseDataList[[]dto.WorkSpaceDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetWorkSpaces"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.WORKSPACE, domain.VIEW); err != nil {
		return res, domain.ACTION_NOT_ALLOWED
	}

	res.Body.Result, err = u.repo.GetWorkSpaces(req, d)
	if err != nil {
		return
	}

	return
}
func (u *workSpaceUseCase) GetWorkSpace(req *common.RequestContext, d dto.RequestEntity) (
	res dto.ResultEntity[dto.WorkSpaceDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetWorkSpace"))
		}
	}()
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.WORKSPACE, domain.VIEW); err != nil {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res.Entity, err = u.repo.GetWorkSpace(req, d)
	if err != nil {
		return
	}
	res.Entity.Modules, err = u.moduleUseCase.GetModules(req, dto.ModulesRequest{
		WorkSpaceID: strconv.Itoa(int(res.Entity.ID)),
		DefaultListParams: dto.DefaultListParams{
			OrderColumn: "wm.created_at",
		},
	})
	

	res.Activities = u.core.GerActivitiesByPartyID(req, res.Entity.ID)
	return
}
