package contact_repo

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/domain"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"
	"fmt"
	"strings"

	"gorm.io/gen/helper"
)

type ContactRepository interface {
	CreateContact(req *common.RequestContext, i dto.ContactData) (
		res model.Contact, err error)
	OnCreateContact(ctx context.Context, tx *query.QueryTx,
		company model.Company, profile model.Profile, i dto.ContactData) (
		res model.Contact, err error)
	GetContacts(req *common.RequestContext, i dto.ContactsRequest) (
		res []dto.ContactDto, err error)
	GetContact(req *common.RequestContext, i *dto.RequestEntity) (dto.ResultEntity[dto.ContactDto], error)
	EditContact(req *common.RequestContext, d dto.ContactData) (err error)
	ContactBulk(req *common.RequestContext, d dto.ContactBulkData) (err error)
	ContactBulkTx(tx *query.QueryTx, req *common.RequestContext, d dto.ContactBulkData) (err error)
}

type partyContactRepository struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
	query     helpers.QueryHelper
	dbHelper  db.DbHelper
}

func NewContactRepository(
	conn db.Connection,
	helpers *helpers.Helpers,
) ContactRepository {
	return &partyContactRepository{
		convertor: helpers.Convertor,
		Q:         conn.GetQ(),
		query:     helpers.Query,
		dbHelper:  conn.GetDbHelper(),
	}
}

func (r *partyContactRepository) ContactBulk(req *common.RequestContext, d dto.ContactBulkData) (err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	err = r.ContactBulkTx(tx, req, d)
	return tx.Commit()
}

func (r *partyContactRepository) ContactBulkTx(tx *query.QueryTx, req *common.RequestContext, d dto.ContactBulkData) (err error) {
	for _, contact := range d.Contacts {
		switch contact.Action {
		case string(domain.EDIT):
			fmt.Println("CONTACT TO EDIT",contact)
			err = r.editContact(tx, req, contact)
			if err != nil {
				return
			}
			partYReference := model.PartyReference{
				ReferenceID: d.PartyID,
				PartyID:     contact.ID,
			}
			err = tx.PartyReference.Save(&partYReference)
		case string(domain.CREATE):
			contact.ReferenceID = &d.PartyID
			_, err = r.createContact(tx, req, contact)
		case string(domain.DELETE):
			_, err = tx.PartyReference.Where(
				tx.PartyReference.PartyID.Eq(contact.ID),
				tx.PartyReference.ReferenceID.Eq(*contact.ReferenceID),
			).Delete()
		}

	}
	return
}

func (r *partyContactRepository) EditContact(req *common.RequestContext, d dto.ContactData) (err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}

	}()
	err = r.editContact(tx, req, d)
	if err != nil {
		return
	}
	err = tx.Commit()

	return
}

func (r *partyContactRepository) editContact(tx *query.QueryTx, req *common.RequestContext, d dto.ContactData) (err error) {
	data, err := r.convertor.DataMap(d.Fields)
	if err != nil {
		return
	}
	err = tx.Contact.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.Contact{ID: d.ID},
	).Updates(data).Error
	if err != nil {
		return
	}

	err = tx.Contact.InsertActivity(d.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	return
}

func (r *partyContactRepository) CreateContact(req *common.RequestContext, i dto.ContactData) (
	res model.Contact, err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	res, err = r.createContact(tx, req, i)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}
func (r *partyContactRepository) createContact(tx *query.QueryTx, req *common.RequestContext, i dto.ContactData) (
	res model.Contact, err error) {
	res, err = r.OnCreateContact(req.Ctx, tx, req.ActiveCompany, req.Profile, i)
	return
}

func (r *partyContactRepository) OnCreateContact(ctx context.Context, tx *query.QueryTx,
	company model.Company, profile model.Profile, i dto.ContactData) (
	res model.Contact, err error) {
	partyId, err := tx.Contact.InsertParty(domain.PARTY_CONTACT)
	if err != nil {
		return res, err
	}
	res.ID = partyId
	res.Name = i.Fields.Name
	res.Gender = i.Fields.Gender
	res.PhoneNumber = i.Fields.PhoneNumber
	res.Email = i.Fields.Email
	res.CompanyID = company.ID
	err = tx.Contact.Save(&res)
	if i.ReferenceID != nil && *i.ReferenceID != 0 {
		var partyReference model.PartyReference
		partyReference.PartyID = res.ID
		partyReference.ReferenceID = *i.ReferenceID
		err := tx.PartyReference.WithContext(ctx).Save(&partyReference)
		if err != nil {
			return res, err
		}
	}
	return res, err
}

func (r *partyContactRepository) GetContacts(req *common.RequestContext, d dto.ContactsRequest) (
	res []dto.ContactDto, err error) {
	var (
		generateSQL strings.Builder
	)
	builder := r.Q.WithContext(req.Ctx).Deal
	queryData := r.convertor.GenerateQueryMap(d)
	params := r.contactQuery(req, queryData, &generateSQL)
	fmt.Println("WHERE SQL", generateSQL.String(),d.PartyID)
	err = builder.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res).Error
	return res, err
}

func (r *partyContactRepository) contactQuery(req *common.RequestContext, d map[string]string, generateSQL *strings.Builder,
) (params []interface{}) {	
	var (
		whereSQL strings.Builder
	)
	generateSQL.WriteString(`SELECT 
				e.id,e.uuid,e.created_at,e.name,e.gender,
				e.email,e.phone_number
				from contacts as e 
			`)
	whereSQL.WriteString(` e.deleted_at is null and e.company_id = ? `)
	params = append(params, req.ActiveCompany.ID)
	columnFilters := []string{
		"name",
	}
	r.query.FilterBuilder(&whereSQL, &params, d, columnFilters...)
	r.query.PartyFilterBuilder(generateSQL, &whereSQL, &params, d, "party_id")


	helper.JoinWhereBuilder(generateSQL, whereSQL)
	
	r.query.OrderAndLimitBuilder(generateSQL, d)
	return
}

func (r *partyContactRepository) GetContact(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[dto.ContactDto], err error) {
	contactQ := r.Q.Contact
	id := r.convertor.StrtoInt(i.ID)
	err = r.Q.Contact.WithContext(req.Ctx).Select(
		contactQ.Email, contactQ.Name, contactQ.Gender,
		contactQ.PhoneNumber, contactQ.ID, contactQ.UUID,
	).Where(
		r.Q.Contact.CompanyID.Eq(req.ActiveCompany.ID),
		r.Q.Contact.ID.Eq(id),
	).Limit(1).Scan(&res.Entity)
	if err != nil {
		return res, err
	}
	return res, err
}
