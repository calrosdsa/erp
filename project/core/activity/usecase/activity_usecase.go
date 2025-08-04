package activity_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/di"
	"erp/pkg/logger"
	activity_repo "erp/project/core/activity/repository"
)

type ActivityUseCase interface {
	CreateActivity(req *common.RequestContext, i dto.ActivityData) (
		res dto.ActivityDto,err error)
	EditActivity(req *common.RequestContext, i dto.ActivityData) error
	DeleteActivity(req *common.RequestContext, i dto.DeleteRequest) error
	CreateActivityStatus(tx *query.QueryTx,from string,to string,activity *model.Activity)(err error)
	// GerActivitiesByPartyID(req *common.RequestContext, i *dto.RequestEntity) (
	// 	res []dto.ActivityDto, err error)
}

type activityUseCase struct {
	emitLog      logger.EmitLog
	activityRepo activity_repo.ActivityRepository
	bus          bus.Bus
	c            di.Container
}

func NewActivityUseCase(
	logger logger.Logger,
	activityRepo activity_repo.ActivityRepository,
	bus bus.Bus,
	c di.Container,
) ActivityUseCase {
	return &activityUseCase{
		emitLog:      logger.EmitLog("activity-usecase"),
		activityRepo: activityRepo,
		bus:          bus,
		c:            c,
	}
}

func(r *activityUseCase) CreateActivityStatus(tx *query.QueryTx,from string,to string,activity *model.Activity)(err error) {
	defer  func(){
		if err != nil {
			r.emitLog.Err(err,logger.OptionsLog.WithMethod("CreateActivityStatus"))
		}
	}()
	err = r.activityRepo.CreateActivityStatus(tx,from,to,activity)
	return
}

func (r *activityUseCase) EditActivity(req *common.RequestContext, i dto.ActivityData) (err error) {
	ctx := r.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DbKey).(*query.Query).Begin()
	defer func(tx *query.QueryTx) {
		if err != nil {
			r.emitLog.Err(err, logger.OptionsLog.WithMethod("EditActivity"))
		}
		err = r.closeTx(tx, err)
	}(tx)

	err = r.activityRepo.EditActivity(tx, req, i)
	return
}
func (r *activityUseCase) DeleteActivity(req *common.RequestContext, i dto.DeleteRequest) (err error) {
	ctx := r.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DbKey).(*query.Query).Begin()
	defer func(tx *query.QueryTx) {
		if err != nil {
			r.emitLog.Err(err, logger.OptionsLog.WithMethod("DeleteActivity"))
		}
		err = r.closeTx(tx, err)
	}(tx)
	err = r.activityRepo.DeleteActivity(tx, req, i)
	return
}

func (r *activityUseCase) CreateActivity(req *common.RequestContext, i dto.ActivityData) (
	res dto.ActivityDto,err error) {
	ctx := r.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DbKey).(*query.Query).Begin()
	defer func(tx *query.QueryTx) {
		if err != nil {
			r.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateActivity"))
		}
		err = r.closeTx(tx, err)

	}(tx)

	res,err = r.activityRepo.CreateActivity(tx, req, i)
	if err != nil {
		return
	}
	payload := event.ActivityEventData{
		Data:i,
		Tx: tx,
		ReqCtx: *req,
	}
	r.bus.Emit(req.Ctx,domain.ActivityCreated,payload)
	return
}

func (u *activityUseCase) closeTx(tx *query.QueryTx, err error) error {
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
