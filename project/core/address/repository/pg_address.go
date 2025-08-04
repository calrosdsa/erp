package address_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/pkg/db"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

type AddressRepository interface {
	CreateAddress(req *common.RequestContext, d dto.AddressData) (dto.AddressDto, error)
	EditAddress(req *common.RequestContext, d dto.AddressData) error
	GetAddresses(req *common.RequestContext, d dto.RequestAddresses) (dto.ResponseDataList[[]dto.AddressDto], error)
	GetAddress(req *common.RequestContext, d *dto.RequestEntity) (dto.ResultEntity[dto.AddressDto], error)
	GetAllowedPartiesForAddress(req *common.RequestContext) dto.ResultEntity[[]dto.PartyTypeDto]
	UpdateStatus(req *common.RequestContext, d dto.UpdateStatusWithEvent, nextState string) (
		err error)
}

type addressRepository struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
	locale    helpers.Locale
	dbHelper  db.DbHelper
	
}

func NewAddressRepository(
	db db.Connection,
	helpers *helpers.Helpers,
	
) AddressRepository {
	return &addressRepository{
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
		locale:    helpers.Locale,
		dbHelper:  db.GetDbHelper(),
		
	}
}

func (r *addressRepository) UpdateStatus( req *common.RequestContext, d dto.UpdateStatusWithEvent, nextState string) (
	err error) {
	e := r.Q.Address
	id:= r.convertor.StrtoInt(d.Body.PartyID)
	_, err = r.Q.Address.WithContext(req.Ctx).Where(
		e.CompanyID.Eq(req.ActiveCompany.ID),
		e.Status.Eq(d.Body.CurrentState),
		e.ID.Eq(id),
	).UpdateSimple(
		e.Status.Value(nextState),
	)
	return
}

func (r *addressRepository) GetAddress(req *common.RequestContext, d *dto.RequestEntity) (
	dto.ResultEntity[dto.AddressDto], error) {
	id := r.convertor.StrtoInt(d.ID)
	var (
		res dto.ResultEntity[dto.AddressDto]
		err error
	)
	address, err := r.Q.Address.WithContext(req.Ctx).Where(
		r.Q.Address.CompanyID.Eq(req.ActiveCompany.ID),
		r.Q.Address.ID.Eq(id),
	).First()
	if err != nil {
		return res, err
	}
	res.Entity = dto.AddressDtoFromModel(address)


	return res, err
}

func (r *addressRepository) GetAddresses(req *common.RequestContext, d dto.RequestAddresses) (
	res dto.ResponseDataList[[]dto.AddressDto], err error) {
	var (
		conds []gen.Condition
		order field.Expr
	)
	addressQ := r.Q.Address
	builder := r.Q.WithContext(req.Ctx).Address

	//ADDING CONDITIONS
	conds = append(conds, addressQ.CompanyID.Eq(req.ActiveCompany.ID))
	if d.Title != "" {
		conds = append(conds, addressQ.Title.Eq(d.Title))
	}
	builder = builder.Where(conds...)

	orderCol, ok := r.Q.Address.GetFieldByName(d.OrderColumn) // maybe orderColStr == "id"
	if ok {
		if d.Order == "ASC" {
			order = orderCol.Asc()
		} else {
			order = orderCol.Desc()
		}
		builder = builder.Order(order)
		// User doesn't contains orderColStr
	}

	err = builder.Limit(int(r.convertor.StrtoInt(d.Size))).Scan(&res.Body.Result)
	return res, err
}

func (r *addressRepository) CreateAddress(req *common.RequestContext, d dto.AddressData) (
	res dto.AddressDto, err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	res, err = r.createAddress(tx, req, d)
	if err != nil {
		return
	}
	if d.ReferenceID != nil {
		var references []*int64
		references = append(references, d.ReferenceID)
		err = r.dbHelper.InsertReferences(req.Ctx, tx, res.ID, references)
		if err != nil {
			return
		}
	}
	err = tx.Commit()
	return
}

func (r *addressRepository) EditAddress(req *common.RequestContext, d dto.AddressData) (err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	err = r.editAddress(tx, req, d)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (r *addressRepository) editAddress(tx *query.QueryTx, req *common.RequestContext, d dto.AddressData) (err error) {
	data, err := r.convertor.DataMap(d.Fields)
	if err != nil {
		return
	}
	err = tx.Address.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.Address{ID: d.ID},
	).Updates(data).Error
	if err != nil {
		return
	}

	err = tx.Address.InsertActivity(d.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	return
}

func (r *addressRepository) createAddress(tx *query.QueryTx, req *common.RequestContext, d dto.AddressData) (res dto.AddressDto, err error) {
	var address model.Address
	id, err := tx.Address.InsertParty(proto.PartyType_address.String())
	if err != nil {
		return
	}
	fields := d.Fields
	address.ID = id
	address.CompanyID = req.ActiveCompany.ID
	if err = r.convertor.CopyStructData(fields, &address); err != nil {
		return
	}

	err = tx.WithContext(req.Ctx).Address.Save(&address)
	if err != nil {
		return
	}
	if err = tx.Address.InsertActivity(id, req.Profile.ID, proto.ActivityType_CREATE.String(), nil); err != nil {
		return
	}
	res = dto.AddressDtoFromModel(&address)
	return
}


func (r *addressRepository) GetAllowedPartiesForAddress(req *common.RequestContext) dto.ResultEntity[[]dto.PartyTypeDto] {
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
	// err := s.Db.WithContext(ctx).Find(&res.Entity).Error
	// if err != nil {
	// 	s.emitLog.Err(err,logger.OptionsLog.WithMethod("GetUserPartyTypes"))
	// }

	return res
}
