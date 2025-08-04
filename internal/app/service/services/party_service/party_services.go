package partyservice

import (
	"erp/internal/app/connection"
	"erp/internal/app/domain/repository"
	"erp/internal/app/service/helpers"
	"erp/pkg/logger"
	"erp/pkg/permission"
	"time"
)

func NewPartyServices(
	conn *connection.Connection,
	timeout time.Duration,
	helpers *helpers.Helpers,
	repositories *repository.Repositories,
	permissionService permission.PermissionService,
	logger logger.Logger,
) *repository.PartyServices{
	// partyAddressService := NewPartyAddressService(
	// 	timeout,helpers,repositories,permissionService,logger,
	// )
	partyService := NewPartyService(
		conn,timeout,helpers,repositories,
	)
	
	return &repository.PartyServices{
		// PartyAddress:partyAddressService,
		PartyService: partyService,
	}
}