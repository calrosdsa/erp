package repository

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
)

type PartyRepositories struct {
	PartyAddress    PartyAddressRepository
	PartyRepository PartyRepository
	PartyContact    PartyContactRepository
}

type PartyServices struct {
	PartyAddress PartyAddressService
	PartyContact PartyContactService
	PartyService PartyService
}

type PartyRepository interface {
	GetUserPartyTypes(req *common.RequestContext) dto.ResultEntity[[]dto.PartyTypeDto]
	GetPartiesByReference(ctx context.Context, req *common.RequestContext, i *dto.RequestPartyReference) (
		dto.ResultEntity[[]dto.PartyDto], error,
	)
	AddPartyReference(ctx context.Context, req *common.RequestContext, i *dto.RequestAddPartyReference) error
	GetPartyReferences(ctx context.Context, req *common.RequestContext, i *dto.RequestPaginationData) (
		dto.PaginationResult[[]dto.PartyReferenceDto], error)

	GetAllowedPartiesForReferences(req *common.RequestContext)([]dto.PartyTypeDto)

}

type PartyService interface {
	GetUserPartyTypes(req *common.RequestContext) dto.ResultEntity[[]dto.PartyTypeDto]
	GetPartiesByReference(req *common.RequestContext, i *dto.RequestPartyReference) (
		dto.ResultEntity[[]dto.PartyDto], error,
	)
	AddPartyReference(req *common.RequestContext, i *dto.RequestAddPartyReference) error
	GetPartyReferences(req *common.RequestContext, i *dto.RequestPaginationData) (
		dto.PaginationResult[[]dto.PartyReferenceDto], error)

	GetAllowedPartiesForReferences(req *common.RequestContext)([]dto.PartyTypeDto)
}
