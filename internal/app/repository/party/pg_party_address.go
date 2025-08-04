package party

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/app/connection"
	"erp/internal/app/domain"
	"erp/internal/app/domain/repository"
	"erp/internal/app/service/helpers"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

type partyAddressRepository struct {
	conn      *connection.Connection
	convertor repository.ConvertorHelper
	locale    helpers.Locale
}

func NewPartyAddressRepository(
	conn *connection.Connection,
	helpers *helpers.Helpers,
) repository.PartyAddressRepository {
	return &partyAddressRepository{
		conn:      conn,
		convertor: helpers.Convertor,
		locale:    helpers.Locale,
	}
}

func (r *partyAddressRepository) GetAddress(ctx context.Context, req *common.RequestContext, i *dto.RequestEntity) (
	dto.ResultEntity[dto.AddressDto], error) {
	var (
		res dto.ResultEntity[dto.AddressDto]
		err error
	)
	address, err := r.conn.Q.Address.WithContext(ctx).Where(
		r.conn.Q.Address.CompanyID.Eq(req.ActiveCompany.ID),
		r.conn.Q.Address.UUID.Eq(i.ID),
	).First()
	if err != nil {
		return res, err
	}
	res.Entity = dto.AddressDtoFromModel(address)
	return res, err
}

func (r *partyAddressRepository) GetAddresses(ctx context.Context, req *common.RequestContext, i dto.RequestAddresses) (
	res dto.ResponseDataList[[]dto.AddressDto], err error) {
	var (
		conds []gen.Condition
		order field.Expr
	)
	addressQ := r.conn.Q.Address
	builder := r.conn.Q.WithContext(ctx).Address

	//ADDING CONDITIONS
	conds = append(conds, addressQ.CompanyID.Eq(req.ActiveCompany.ID))
	if i.Title != "" {
		conds = append(conds, addressQ.Title.Eq(i.Title))
	}
	builder = builder.Where(conds...)

	orderCol, ok := r.conn.Q.Address.GetFieldByName(i.OrderColumn) // maybe orderColStr == "id"
	if ok {
		if i.Order == "ASC" {
			order = orderCol.Asc()
		} else {
			order = orderCol.Desc()
		}
		builder = builder.Order(order)
		// User doesn't contains orderColStr
	}

	err = builder.Limit(int(r.convertor.StrtoInt(i.Size))).Scan(&res.Body.Result)
	return res, err
}

func (r *partyAddressRepository) CreatePartyAddress(ctx context.Context, req *common.RequestContext, i *dto.AddressDataRequest) error {
	tx := r.conn.Q.Begin()
	var err error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	// address, err := r.insertAddress(ctx, tx, &i.Body.Address, req.ActiveCompany.ID)
	// if err != nil {
	// 	return err
	// }
	// if i.Body.ReferenceID != nil && *i.Body.ReferenceID != 0 {
	// 	var partyReference model.PartyReference
	// 	partyReference.PartyID = address.ID
	// 	partyReference.ReferenceID = *i.Body.ReferenceID
	// 	err := tx.PartyReference.WithContext(ctx).Save(&partyReference)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return tx.Commit()
}

func (r *partyAddressRepository) insertAddress(ctx context.Context, tx *query.QueryTx, i *dto.AddressRequestData, companyID int64) (model.Address, error) {
	var (
		address model.Address
	)
	partyId, err := tx.Address.WithContext(ctx).InsertParty(domain.PARTY_ADDRESS)
	if err != nil {
		return address, err
	}
	address.ID = partyId
	address.Title = i.Title
	address.City = i.City
	address.Company = &i.Company
	address.CountryCode = &i.CountryCode
	address.Email = &i.Email
	address.StreetLine1 = i.StreetLine1
	address.StreetLine2 = i.StreetLine2
	address.IdentificationNumber = &i.IdentificationNumber
	address.PhoneNumber = &i.PhoneNumber
	address.PostalCode = &i.PostalCode
	address.Province = &i.Province
	address.CompanyID = companyID
	address.IsBillingAddress = i.IsBillingAddress
	address.IsShippingAddress = i.IsShippingAddress
	address.Enabled = i.Enabled

	err = tx.Address.WithContext(ctx).Save(&address)
	return address, err
}

func (r *partyAddressRepository) GetAllowedPartiesForAddress(req *common.RequestContext) dto.ResultEntity[[]dto.PartyTypeDto] {
	var res dto.ResultEntity[[]dto.PartyTypeDto]
	userPartyTypes := []dto.PartyTypeDto{
		{Name: r.locale.MustLocalize(
			helpers.OptionsLocale.WithID("Party.customer"),
			helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
		), Code: domain.PARTY_CUSTOMER},
		{Name: r.locale.MustLocalize(
			helpers.OptionsLocale.WithID("Party.supplier"),
			helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
		), Code: domain.PARTY_SUPPLIER},
		{Name: r.locale.MustLocalize(
			helpers.OptionsLocale.WithID("Party.warehouse"),
			helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
		), Code: domain.PARTY_WAREHOUSE},
		// {Name: r.locale.MustLocalize(
		// 	helpers.OptionsLocale.WithID("Party.Company"),
		// 	helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
		// ), Code: domain.PARTY_COMPANY},
	}
	res.Entity = userPartyTypes
	// err := s.conn.Db.WithContext(ctx).Find(&res.Entity).Error
	// if err != nil {
	// 	s.emitLog.Err(err,logger.OptionsLog.WithMethod("GetUserPartyTypes"))
	// }

	return res
}
