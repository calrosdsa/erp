package buying

import (
	"erp/internal/app/connection"
	"erp/internal/app/domain/repository"
	"erp/internal/app/service/helpers"
	"erp/pkg/logger"
	"erp/pkg/permission"
	"time"
)

func NewBuyingServices(
	conn *connection.Connection,
	timeout time.Duration,
	helpers *helpers.Helpers,
	permissionService permission.PermissionService,
	repositories *repository.Repositories,
	logger logger.Logger,
) *repository.BuyingServices {
	purchaseService := NewPurchaseService(conn, timeout, helpers, permissionService, repositories,logger)
	return &repository.BuyingServices{
		PurchaseService: purchaseService,
	}
}
