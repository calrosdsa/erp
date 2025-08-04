package buying

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/app/connection"
	"erp/internal/app/domain/repository"
	"erp/internal/app/service/helpers"
	"erp/pkg/logger"
	"erp/pkg/permission"
	"time"
)

type purchaseService struct {
	conn              *connection.Connection
	timeout           time.Duration
	permissionService permission.PermissionService
	emitLog           logger.EmitLog
	purchaseRepository repository.PurchaseRepository
}

func NewPurchaseService(
	conn *connection.Connection,
	timeout time.Duration,
	helpers *helpers.Helpers,
	permissionService permission.PermissionService,
	repositories *repository.Repositories,
	logger logger.Logger,
) repository.PurchaseService {
	return &purchaseService{
		conn:              conn,
		timeout:           timeout,
		permissionService: permissionService,
		emitLog:           logger.EmitLog("purchase-service"),
		purchaseRepository: repositories.BuyingRepository.PurchaseRepository,
	}
}


func (s *purchaseService) CreatePurchaseOrder(req *common.RequestContext, i *dto.CreatePurchaseOrderRequest) (
	dto.ResultEntity[dto.OrderDto], error) {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	var (
		err error
		res dto.ResultEntity[dto.OrderDto]
	)
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("CreatePurchase"))
		}
	}()

	res,err = s.purchaseRepository.CreatePurchaseOrder(ctx,req,i)
	return res, nil
}

