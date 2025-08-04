package repository

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
)

type PartyContactRepository interface {
	CreateContact(ctx context.Context, req *common.RequestContext, i *dto.ContactData) (dto.ContactDto, error)
	GetContacts(ctx context.Context, req *common.RequestContext, i *dto.RequestPaginationData) (
		dto.PaginationResult[[]dto.ContactDto], error)
	GetContact(ctx context.Context, req *common.RequestContext, i *dto.RequestEntity) (dto.ContactDto, error)

}

type PartyContactService interface {
	CreateContact(req *common.RequestContext, i *dto.ContactDataRequest) (dto.ContactDto, error)
	GetContacts(req *common.RequestContext, i *dto.RequestPaginationData) (
		dto.PaginationResult[[]dto.ContactDto], error)
	GetContact(req *common.RequestContext, i *dto.RequestEntity) (dto.ContactDto, error)

}
