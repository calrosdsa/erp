package repository

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
)

type PartyAddressRepository interface {
	CreatePartyAddress(ctx context.Context, req *common.RequestContext, i *dto.AddressDataRequest) error
	GetAddresses(ctx context.Context, req *common.RequestContext, i dto.RequestAddresses) (dto.ResponseDataList[[]dto.AddressDto], error)
	GetAddress(ctx context.Context, req *common.RequestContext, i *dto.RequestEntity) (dto.ResultEntity[dto.AddressDto], error)

	GetAllowedPartiesForAddress(req *common.RequestContext) dto.ResultEntity[[]dto.PartyTypeDto]
}

type PartyAddressService interface {
	CreatePartyAddress(req *common.RequestContext, i *dto.AddressDataRequest) error
	GetAddresses(req *common.RequestContext, i dto.RequestAddresses) (dto.ResponseDataList[[]dto.AddressDto], error)
	GetAddress(req *common.RequestContext, i *dto.RequestEntity) (dto.ResultEntity[dto.AddressDto], error)
	GetAllowedPartiesForAddress(req *common.RequestContext) dto.ResultEntity[[]dto.PartyTypeDto]
}
