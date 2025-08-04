package party

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/internal/app/connection"
	"erp/internal/app/domain"
	"erp/internal/app/domain/repository"
	"erp/internal/app/service/helpers"
	"fmt"

	"gorm.io/gen/field"
	"gorm.io/gorm/schema"
)

type partyRepository struct {
	conn   *connection.Connection
	locale helpers.Locale
	convertor helpers.ConvertorHelper
}

func NewPartyRepository(
	conn *connection.Connection,
	helpers *helpers.Helpers,
) repository.PartyRepository {
	return &partyRepository{
		conn:   conn,
		locale: helpers.Locale,
		convertor: helpers.Convertor,
	}
}

func (r *partyRepository) GetPartyReferences(ctx context.Context,req *common.RequestContext,i *dto.RequestPaginationData)(
	dto.PaginationResult[[]dto.PartyReferenceDto],error){
	var (
		res dto.PaginationResult[[]dto.PartyReferenceDto]
	)	
	warehouseQ := r.conn.Q.WareHouse 
	customerQ := r.conn.Q.Customer
	supplierQ := r.conn.Q.Supplier
	partyQ := r.conn.Q.Party
	partyReferenceQ := r.conn.Q.PartyReference
	limit,offset := r.convertor.ToPaginationParams(i.Page,i.Size)
	partyID := r.convertor.StrtoInt(i.FilterID)
	err := r.conn.Q.PartyReference.WithContext(ctx).Select(
		warehouseQ.Name,warehouseQ.UUID,
		customerQ.Name,customerQ.UUID,
		supplierQ.Name,supplierQ.UUID,
		partyQ.PartyTypeCode.As("party_type"),
	).
	Join(partyQ,partyReferenceQ.ReferenceID.EqCol(partyQ.ID)).
	LeftJoin(customerQ,partyReferenceQ.ReferenceID.EqCol(customerQ.ID)).
	LeftJoin(supplierQ,partyReferenceQ.ReferenceID.EqCol(supplierQ.ID)).
	LeftJoin(warehouseQ,partyReferenceQ.ReferenceID.EqCol(warehouseQ.ID)).
	Where(partyReferenceQ.PartyID.Eq(partyID)).
	Limit(limit).Offset(offset).Scan(&res.Results)
	for i,reference := range res.Results {
		reference.Reference = r.locale.MustLocalize(
			helpers.OptionsLocale.WithID(fmt.Sprintf("Party.%s",reference.PartyType)),
			helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
		)
		res.Results[i] = reference
	}
	return res,err
}
// func (r *partyRepository) GetPartyReferences(ctx context.Context,)
func (r *partyRepository) AddPartyReference(ctx context.Context, req *common.RequestContext, i *dto.RequestAddPartyReference) error {
	var partyReference model.PartyReference
	partyReference.PartyID = i.Body.PartyID
	partyReference.ReferenceID = i.Body.ReferenceID
	err := r.conn.Q.PartyReference.Save(&partyReference)
	return err
}

func (r *partyRepository) GetPartiesByReference(ctx context.Context, req *common.RequestContext, i *dto.RequestPartyReference) (
	dto.ResultEntity[[]dto.PartyDto], error,
) {
	warehouseQ := r.conn.Q.WareHouse
	customerQ := r.conn.Q.Customer
	supplierQ := r.conn.Q.Supplier
	partyQ := r.conn.Q.Party
	var (
		res       dto.ResultEntity[[]dto.PartyDto]
		sColumns  []field.Expr
		joinExpr  []field.Expr
		joinTable schema.Tabler
	)
	builder := r.conn.Q.WithContext(ctx)
	switch i.PartyType {
	case domain.PARTY_CUSTOMER:
		joinTable = customerQ
		sColumns = append(sColumns, customerQ.ID.As("id"), customerQ.Name.As("name"),customerQ.UUID.As("uuid"))
		joinExpr = append(joinExpr, partyQ.ID.EqCol(customerQ.ID), customerQ.CompanyID.Eq(req.ActiveCompany.ID))
		if i.Query != "" {
			joinExpr = append(joinExpr, customerQ.Name.Like("%"+i.Query+"%"))
		}
	case domain.PARTY_SUPPLIER:
		joinTable = supplierQ
		sColumns = append(sColumns, supplierQ.ID.As("id"), supplierQ.Name.As("name"),supplierQ.UUID.As("uuid"))
		joinExpr = append(joinExpr, partyQ.ID.EqCol(supplierQ.ID), supplierQ.CompanyID.Eq(req.ActiveCompany.ID))
		if i.Query != "" {
			joinExpr = append(joinExpr, supplierQ.Name.Like("%"+i.Query+"%"))
		}
	case domain.PARTY_WAREHOUSE:
		joinTable = warehouseQ
		sColumns = append(sColumns, warehouseQ.ID.As("id"), warehouseQ.Name.As("name"),warehouseQ.UUID.As("uuid"))
		joinExpr = append(joinExpr, partyQ.ID.EqCol(warehouseQ.ID), warehouseQ.CompanyID.Eq(req.ActiveCompany.ID))
		if i.Query != "" {
			joinExpr = append(joinExpr, warehouseQ.Name.Like("%"+i.Query+"%"))
		}
	default:
		return res, domain.TYPE_NOT_FOUND
	}
	err := builder.Party.Select(sColumns...).Join(joinTable, joinExpr...).Limit(domain.DEFAULT_LIMIT).Scan(&res.Entity)
	return res, err
}
func (r *partyRepository) GetUserPartyTypes(req *common.RequestContext) dto.ResultEntity[[]dto.PartyTypeDto] {
	var res dto.ResultEntity[[]dto.PartyTypeDto]
	userPartyTypes := []dto.PartyTypeDto{
		{Name: r.locale.MustLocalize(
			helpers.OptionsLocale.WithID("Party.employee"),
			helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
		), Code: domain.PARTY_EMPLOYEE},
		{Name: r.locale.MustLocalize(
			helpers.OptionsLocale.WithID("Party.admin"),
			helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
		), Code: domain.PARTY_ADMIN},
		{Name: r.locale.MustLocalize(
			helpers.OptionsLocale.WithID("Party.client"),
			helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
		), Code: domain.PARTY_CLIENT},
	}
	res.Entity = userPartyTypes

	return res
}


func (r *partyRepository)GetAllowedPartiesForReferences(req *common.RequestContext)([]dto.PartyTypeDto) {
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
	return userPartyTypes
}
