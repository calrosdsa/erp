package repository

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
)

type CoreService interface {
	GerActivitiesByPartyID(req *common.RequestContext, partyID int64) (
		res []dto.ActivityDto)
	GetPartyContacts(req *common.RequestContext, partyID int64) []dto.ContactDto
	GetPartyAddresses(req *common.RequestContext, partyID int64) []dto.AddressDto
	GetPartyAccountingDimension(req *common.RequestContext, partyID int64) (dto.AccountingDimensionDto, error)
	GetPartyIDByType(ctx context.Context, partyID int64, partyType string) (*int64, error)

	GetContact(req *common.RequestContext, id int64) (res dto.ContactDto, err error)
	GetAddress(req *common.RequestContext, id int64) (res dto.AddressDto, err error)
}
