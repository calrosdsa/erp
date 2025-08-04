package partyservice

// import (
// 	"context"
// 	"erp/api/common"
// 	"erp/api/dto"
// 	"erp/internal/domain"
// 	"erp/internal/app/domain/repository"
// 	"erp/internal/app/service/helpers"
// 	"erp/pkg/logger"
// 	"erp/pkg/permission"
// 	"time"
// )

// type partyAddressService struct {
// 	timeout                time.Duration
// 	emitLog                logger.EmitLog
// 	partyAddressRepository repository.PartyAddressRepository
// 	permissionService      permission.PermissionService
// }

// func NewPartyAddressService(
// 	timeout time.Duration,
// 	helpers *helpers.Helpers,
// 	repositories *repository.Repositories,
// 	permissionService permission.PermissionService,
// 	logger logger.Logger,
// ) repository.PartyAddressService {
// 	return &partyAddressService{
// 		timeout:                timeout,
// 		partyAddressRepository: repositories.PartyRepositories.PartyAddress,
// 		emitLog:                logger.EmitLog("party-address-service"),
// 		permissionService:      permissionService,
// 	}
// }

// func (s *partyAddressService) GetAddresses(req *common.RequestContext, i *dto.RequestPaginationData) (
// 	res dto.PaginationResult[[]dto.AddressDto], err error) {
// 	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
// 	defer func() {
// 		cancel()
// 		if err != nil {
// 			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetAddresses"))
// 		}
// 	}()
// 	if allow := s.permissionService.CheckPermission(ctx, req, domain.ADDRESS, domain.VIEW); !allow {
// 		return res, domain.ACTION_NOT_ALLOWED
// 	}
// 	// res, err = s.partyAddressRepository.GetAddresses(ctx, req, i)
// 	return res, err
// }
// func (s *partyAddressService) GetAddress(req *common.RequestContext, i *dto.RequestEntity) (
// 	res dto.ResultEntity[dto.AddressDto], err error) {
// 	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
// 	defer func() {
// 		cancel()
// 		if err != nil {
// 			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetAddress"))
// 		}
// 	}()
// 	if allow := s.permissionService.CheckPermission(ctx, req, domain.ADDRESS, domain.VIEW); !allow {
// 		return res, domain.ACTION_NOT_ALLOWED
// 	}
// 	res, err = s.partyAddressRepository.GetAddress(ctx, req, i)
// 	return res, err
// }

// func (s *partyAddressService) CreatePartyAddress(req *common.RequestContext, i *dto.CreatePartyAddressRequest) (err error) {
// 	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
// 	defer func() {
// 		cancel()
// 		if err != nil {
// 			s.emitLog.Err(err, logger.OptionsLog.WithMethod("CreatePartyAddress"))
// 		}
// 	}()
// 	if allow := s.permissionService.CheckPermission(ctx, req, domain.ADDRESS, domain.CREATE); !allow {
// 		return domain.ACTION_NOT_ALLOWED
// 	}
// 	err = s.partyAddressRepository.CreatePartyAddress(ctx, req, i)
// 	return err
// }

// func (s *partyAddressService) GetAllowedPartiesForAddress(req *common.RequestContext) dto.ResultEntity[[]dto.PartyTypeDto] {
// 	res := s.partyAddressRepository.GetAllowedPartiesForAddress(req)
// 	return res
// }
