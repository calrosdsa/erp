package payment_terms_repo

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"
)

type PaymentTermsLineRepo interface {
	DeletePaymentTermsLines(tx *query.QueryTx, ctx context.Context, documentID int64) (err error)
	CreatePaymentTermLines(tx *query.QueryTx, ctx context.Context, d []dto.PaymentTermsLineData, documentID int64) (err error)
	UpdatePaymentTermsLines(tx *query.QueryTx, ctx context.Context, d []dto.PaymentTermsLineData,
		documentID int64) (err error)
	GetPaymentTermLines(req *common.RequestContext, d dto.RequestEntity) (res []dto.PaymentTermsLineDto, err error)
}
type paymentTermsLineRepo struct {
	Q *query.Query
	convertor helpers.ConvertorHelper
}

func NewPaymentTermsLineRepo(
	db db.Connection,
	helpers *helpers.Helpers,
) PaymentTermsLineRepo {
	return &paymentTermsLineRepo{
		Q: db.GetQ(),
		convertor: helpers.Convertor,
	}
}
func (r *paymentTermsLineRepo) GetPaymentTermLines(req *common.RequestContext, d dto.RequestEntity) (
	res []dto.PaymentTermsLineDto, err error) {
	lineQ := r.Q.PaymentTermsLine
	paymentTermQ := r.Q.PaymentTerm
	err =lineQ.WithContext(req.Ctx).Select(
		lineQ.ID,lineQ.InvoicePortion,lineQ.CreditDays,lineQ.DueDateBaseOn,lineQ.Description,
		lineQ.PaymentTermsID,paymentTermQ.Name.As("payment_term"),
	).
	Join(paymentTermQ,paymentTermQ.ID.EqCol(lineQ.PaymentTermsID)).
	Where(
		lineQ.DocumentID.Eq(r.convertor.StrtoInt(d.ID)),
	).Scan(&res)
	return 
}

func (r *paymentTermsLineRepo) UpdatePaymentTermsLines(tx *query.QueryTx, ctx context.Context,
	d []dto.PaymentTermsLineData, documentID int64) (err error) {
	if err = r.DeletePaymentTermsLines(tx, ctx, documentID); err != nil {
		return
	}
	err = r.CreatePaymentTermLines(tx, ctx, d, documentID)
	return
}

func (r *paymentTermsLineRepo) DeletePaymentTermsLines(tx *query.QueryTx, ctx context.Context, documentID int64) (
	err error) {
	_, err = tx.PaymentTermsLine.WithContext(ctx).Unscoped().Where(
		tx.PaymentTermsLine.DocumentID.Eq(documentID),
	).Delete()
	return
}

func (r *paymentTermsLineRepo) CreatePaymentTermLines(tx *query.QueryTx, ctx context.Context,
	d []dto.PaymentTermsLineData, documentID int64) (err error) {
	paymentTermsLines := make([]*model.PaymentTermsLine, len(d))
	for i, line := range d {
		paymentTermLine := &model.PaymentTermsLine{}
		paymentTermLine.DocumentID = documentID
		paymentTermLine.PaymentTermsID = line.PaymentTermsID
		paymentTermLine.InvoicePortion = line.InvoicePortion
		paymentTermLine.Description = line.Description
		paymentTermLine.CreditDays = line.CreditDays
		paymentTermLine.DueDateBaseOn = line.DueDateBaseOn
		paymentTermsLines[i] = paymentTermLine
	}
	err = tx.PaymentTermsLine.WithContext(ctx).CreateInBatches(paymentTermsLines, len(paymentTermsLines))
	return
}
