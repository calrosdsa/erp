package documentinfo_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"

)

type DocumentInfoRepository interface {
	EditAddressAndContact(req *common.RequestContext, d dto.AddressAndContactData) (err error)
	GetAddressAndContact(req *common.RequestContext, d dto.RequestEntityWithParty) (res dto.AddressAndContactDto, err error)
	EditDocTerms(req *common.RequestContext, d dto.DocTermsData) (err error)
	GetDocTerms(req *common.RequestContext, d dto.RequestEntityWithParty) (res dto.DocTermsDto, err error)
	EditDocAccounts(req *common.RequestContext, d dto.DocAccountingData) (err error)
	GetDocAccounts(req *common.RequestContext, d dto.RequestEntityWithParty) (res dto.DocAccountingDto, err error)
}

type documentInfoRepository struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
}

func NewDocumentInfoRepository(
	conn db.Connection,
	helpers *helpers.Helpers,
) DocumentInfoRepository {
	return &documentInfoRepository{
		Q:         conn.GetQ(),
		convertor: helpers.Convertor,
	}
}
func (r *documentInfoRepository) EditDocTerms(req *common.RequestContext, d dto.DocTermsData) (err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	data, err := r.convertor.DataMap(d.Fields)
	if err != nil {
		return
	}
	err = tx.DocTerm.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.DocTerm{DocID: d.DocID},
	).Updates(data).Error
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}
func (r *documentInfoRepository) GetDocTerms(req *common.RequestContext, d dto.RequestEntityWithParty) (
	res dto.DocTermsDto, err error) {
	docTermQ := r.Q.DocTerm
	tacQ := r.Q.TermsAndCondition
	pttQ := r.Q.PaymentTermsTemplate
	err = docTermQ.Select(
		docTermQ.TermsAndConditionID,tacQ.Name.As("terms_and_condition"),tacQ.UUID.As("terms_and_condition_uuid"),
		docTermQ.PaymentTermTemplateID,pttQ.Name.As("payment_term_template"),pttQ.UUID.As("payment_term_template_uuid"),
	).
	LeftJoin(tacQ,tacQ.ID.EqCol(docTermQ.TermsAndConditionID)).
	LeftJoin(pttQ,pttQ.ID.EqCol(docTermQ.PaymentTermTemplateID)).
	Where(docTermQ.DocID.Eq(r.convertor.StrtoInt(d.ID))).Scan(&res)
	return
}
func (r *documentInfoRepository) EditDocAccounts(req *common.RequestContext, d dto.DocAccountingData) (err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	data, err := r.convertor.DataMap(d.Fields)
	if err != nil {
		return
	}
	err = tx.DocAccount.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.DocAccount{DocID: d.DocID},
	).Updates(data).Error
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}
func (r *documentInfoRepository) GetDocAccounts(req *common.RequestContext, d dto.RequestEntityWithParty) (
	res dto.DocAccountingDto, err error) {
	dcQ := r.Q.DocAccount
	creditQ := r.Q.Ledger.As("credit")
	debitQ := r.Q.Ledger.As("debit")
	err = dcQ.Select(
		dcQ.CreditAccountID,creditQ.Name.As("credit_account"),creditQ.UUID.As("credit_account_uuid"),
		dcQ.DebitAccountID,debitQ.Name.As("debit_account"),debitQ.UUID.As("debit_account_uuid"),
	).
	LeftJoin(creditQ,creditQ.ID.EqCol(dcQ.CreditAccountID)).
	LeftJoin(debitQ,debitQ.ID.EqCol(dcQ.DebitAccountID)).
	Where(
		dcQ.DocID.Eq(r.convertor.StrtoInt(d.ID)),
	).Scan(&res)
	return
}
func (r *documentInfoRepository) GetAddressAndContact(req *common.RequestContext, d dto.RequestEntityWithParty) (
	res dto.AddressAndContactDto, err error) {
	aacQ := r.Q.AddressAndContact
	partyAddrQ := r.Q.Address.As("party_addr")
	billingAddrQ := r.Q.Address.As("billing_addr")
	shippingAddrQ := r.Q.Address.As("shipping_addr")
	contactQ := r.Q.Contact

	err = aacQ.Select(
		aacQ.PartyAddressID,partyAddrQ.UUID.As("party_address_uuid"), partyAddrQ.Title.As("party_address"),
		aacQ.BillingAddressID, billingAddrQ.UUID.As("billing_address_uuid"),billingAddrQ.Title.As("billing_address"),
		aacQ.ShippingAddressID, shippingAddrQ.UUID.As("shipping_address_uuid"),shippingAddrQ.Title.As("shipping_address"),
		aacQ.ContactID, contactQ.UUID.As("contact_uuid"),contactQ.Name.As("contact"),
	).
		LeftJoin(partyAddrQ, partyAddrQ.ID.EqCol(aacQ.PartyAddressID)).
		LeftJoin(billingAddrQ, billingAddrQ.ID.EqCol(aacQ.BillingAddressID)).
		LeftJoin(shippingAddrQ, shippingAddrQ.ID.EqCol(aacQ.ShippingAddressID)).
		LeftJoin(contactQ, contactQ.ID.EqCol(aacQ.ContactID)).
		Where(
			aacQ.DocID.Eq(r.convertor.StrtoInt(d.ID)),
		).Scan(&res)
	return
}

func (r *documentInfoRepository) EditAddressAndContact(req *common.RequestContext, d dto.AddressAndContactData) (err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	data, err := r.convertor.DataMap(d.Fields)
	if err != nil {
		return
	}
	err = tx.AddressAndContact.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.AddressAndContact{DocID: d.ID},
	).Updates(data).Error
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}
