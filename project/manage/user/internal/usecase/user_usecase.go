package user_ucase

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
	user_repo "erp/project/manage/user/internal/repository"
)

type UserUseCase interface {
	CreateUser(req *common.RequestContext, i *dto.CreateUserRequest) error
}

type  userUseCase struct {
	emitLog logger.EmitLog
	permission  repository.PermissionService
	userRepo user_repo.UserRepository
	bus bus.Bus
	c di.Container
}

func NewUserUseCase(
	bus bus.Bus,
	logger logger.Logger,
	permission repository.PermissionService,
	userRepo user_repo.UserRepository,
	c di.Container,
)UserUseCase{
	return &userUseCase{
		emitLog: logger.EmitLog("user-usecase"),
		permission: permission,
		userRepo: userRepo,
		bus: bus,
		c: c,
	}
}
func (u *userUseCase) CreateUser(req *common.RequestContext, i *dto.CreateUserRequest) (err error){
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func(tx *query.QueryTx) {
		if err != nil {

			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateUser"))
		}
		err = domain.CloseTx(tx, err)
	}(tx)
	allow := u.permission.CheckPermission(req.Ctx, req, domain.USER, domain.CREATE)
	if !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	userRelation,err := u.userRepo.CreateUser(tx,req,i)

	err = u.bus.Emit(req.Ctx,domain.UserCreatedEvent,event.UserCreatedEventData{
		UseRelation: userRelation,
		LanguageCode: req.LanguageCode,
		Tx: tx,
	})
	return 
}

