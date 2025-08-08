package task_ucase

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
	task_repo "erp/project/projet_module/task/repository"
)

type TaskUseCase interface {
	CreateTask(req *common.RequestContext, d dto.TaskData) (res dto.TaskDto, err error)
	EditTask(req *common.RequestContext, d dto.TaskData) (err error)
	GetTask(req *common.RequestContext, d dto.RequestEntity) (res dto.ResultEntity[dto.TaskDetailDto], err error)
	GetTasks(req *common.RequestContext, d dto.TasksRequest) (res dto.ResponseDataList[[]dto.TaskDto], err error)
	TaskTransition(req *common.RequestContext, d dto.EntityTransitionData) (err error)
}

type taskUseCase struct {
	permission repository.PermissionService
	core       repository.CoreService
	repo       task_repo.TaskRepository
	emitLog    logger.EmitLog
	bus        bus.Bus
	c          di.Container
}

func NewTaskUseCase(
	permission repository.PermissionService,
	core repository.CoreService,
	repo task_repo.TaskRepository,
	logger logger.Logger,
	bus bus.Bus,
	c di.Container,
) TaskUseCase {
	return &taskUseCase{
		permission: permission,
		core:       core,
		repo:       repo,
		emitLog:    logger.EmitLog("task-usecase"),
		bus:        bus,
		c:          c,
	}
}

func (u *taskUseCase) TaskTransition(req *common.RequestContext, d dto.EntityTransitionData) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DbKey).(*query.Query).Begin()
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("TaskTransition"))
		}
		err = u.closeTx(tx, err)
	}(tx)

	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.TASK, domain.EDIT); err != nil {
		return err
	}
	err = u.repo.TaskTransition(tx, req, d)
	if err != nil {
		return
	}
	err = u.bus.Emit(req.Ctx, domain.PartyStageChange, event.ChangeStageEventData{
		Tx:              tx,
		ProfileID:       req.Profile.ID,
		StageTransition: d,
	})
	return
}

func (u *taskUseCase) CreateTask(req *common.RequestContext, d dto.TaskData) (res dto.TaskDto, err error) {
	tx := u.c.Get(domain.DbKey).(*query.Query).Begin()
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateTask"))
		}
		err = u.closeTx(tx, err)
	}(tx)
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.TASK, domain.CREATE); err != nil {
		return res, domain.ACTION_NOT_ALLOWED
	}
	task, err := u.repo.CreateTask(tx, req, d)
	if err != nil {
		return
	}
	
	res = dto.TaskFromModel(task)
	err = u.bus.Emit(req.Ctx, domain.TaskCreatedEvent, event.TaskEventData{
		Tx:   tx,
		Data: d,
		Task: task,
		Req:  *req,
	})
	return
}

func (u *taskUseCase) EditTask(req *common.RequestContext, d dto.TaskData) (err error) {
	tx := u.c.Get(domain.DbKey).(*query.Query).Begin()
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditTask"))
		}
		err = u.closeTx(tx, err)
	}(tx)

	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.TASK, domain.EDIT); err != nil {
		return err
	}
	err = u.repo.EditTask(tx, req, d)
	if err != nil {
		return
	}
	err = u.bus.Emit(req.Ctx, domain.TaskEditedEvent, event.TaskEventData{
		Tx:   tx,
		Data: d,
		Req:  *req,
	})
	return
}

func (u *taskUseCase) GetTask(req *common.RequestContext, d dto.RequestEntity) (res dto.ResultEntity[dto.TaskDetailDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetTask"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.TASK, domain.VIEW); err != nil {
		return res, err
	}

	res.Entity, err = u.repo.GetTask(req, d)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req, res.Entity.Task.ID)
	res.Contacts = u.core.GetPartyContacts(req, res.Entity.Task.ID)
	return
}

func (u *taskUseCase) GetTasks(req *common.RequestContext, d dto.TasksRequest) (res dto.ResponseDataList[[]dto.TaskDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetTasks"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.TASK, domain.VIEW); err != nil {
		return res, err
	}

	res.Body.Result, err = u.repo.GetTasks(req, d)
	if err != nil {
		return
	}
	return
}

func (u *taskUseCase) closeTx(tx *query.QueryTx, err error) error {
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
