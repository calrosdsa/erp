package serial_no_repo

import (
	"context"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/internal/domain/repository"
	"fmt"
	"time"
)

type SerialNoEventRepository interface {
	OnInvoiceSubmitted(ctx context.Context, payload event.StatusInvoiceEventData) (err error)
	OnReceiptSubmitted(ctx context.Context, payload event.StatusReceiptEventData) (err error)
	OnStockEntrySubmitted(ctx context.Context, payload event.StatusStockEntryEventData) (err error)
	OnInvoiceCancelled(ctx context.Context, payload event.StatusInvoiceEventData) (err error)
	OnReceiptCancelled(ctx context.Context, payload event.StatusReceiptEventData) (err error)
}

type serialNoEventRepo struct {
	generator helpers.Generator
	currency helpers.CurrencyHelper
	accounting repository.AccountingService
}

const SN_TEMPLATE = "SN-#######"
const BATCH_BUNDLE_TEMPLATE = "BB-#######"

func NewSerialEventRepo(
	helpers *helpers.Helpers,
	accounting repository.AccountingService,
) SerialNoEventRepository {
	return &serialNoEventRepo{
		generator: helpers.Generator,
		currency: helpers.Currency,
		accounting: accounting,
	}
}

func (r *serialNoEventRepo) OnReceiptCancelled(ctx context.Context, payload event.StatusReceiptEventData) (err error) {
	switch payload.ReceiptPartyType {
	case proto.PartyType_purchaseReceipt.String():
		err = r.cancellSerialNos(payload.Tx, ctx, payload.Receipt.Code)
	case proto.PartyType_deliveryNote.String():
		err = r.processReturnSale(payload.Tx, ctx, payload.Receipt.Code)
	}
	return
}

func (r *serialNoEventRepo) OnInvoiceCancelled(ctx context.Context, payload event.StatusInvoiceEventData) (err error) {
	if !payload.Invoice.UpdateStock {
		return nil
	}
	switch payload.InvoicePartyType {
	case proto.PartyType_purchaseInvoice.String():
		err = r.cancellSerialNos(payload.Tx, ctx, payload.Invoice.Code)
	case proto.PartyType_saleInvoice.String():
		err = r.processReturnSale(payload.Tx, ctx, payload.Invoice.Code)
	}
	return
}

// func(r *serialNoEventRepo) deleteSerialNoTx(tx *query.QueryTx)(err error){
// 	return err
// }

func (r *serialNoEventRepo) cancellSerialNos(tx *query.QueryTx, ctx context.Context, voucherCode string) (err error) {
	serialQ := tx.SerialNo
	bundleQ := tx.BatchBundle
	batchs, err := bundleQ.Select(bundleQ.ID).WithContext(ctx).Where(
		bundleQ.VoucherCode.Eq(voucherCode),
	).Find()
	if err != nil {
		return
	}
	batchIds := make([]int64, len(batchs))
	for i, batch := range batchs {
		//Delete serialNO transactions
		_, err = tx.SerialNoTransaction.WithContext(ctx).Where(
			tx.SerialNoTransaction.BatchBundleID.Eq(batch.ID),
		).Delete()
		if err != nil {
			return
		}
		//delete serial nos
		_, err = serialQ.WithContext(ctx).Where(
			serialQ.BatchBundleID.Eq(batch.ID),
		).Delete()
		if err != nil {
			return
		}
		batchIds[i] = batch.ID
	}
	_, err = bundleQ.WithContext(ctx).Where(bundleQ.ID.In(batchIds...)).Delete()
	return
}

func (r *serialNoEventRepo) processReturnSale(tx *query.QueryTx, ctx context.Context, voucherCode string) (err error) {
	serialQ := tx.SerialNo
	bundleQ := tx.BatchBundle
	batchs, err := bundleQ.Select(bundleQ.ID).
		Where(bundleQ.VoucherCode.Eq(voucherCode)).Find()
	if err != nil {
		return
	}
	ids := make([]int64, len(batchs))
	for i, batch := range batchs {
		ids[i] = batch.ID
		_, err = tx.SerialNoTransaction.WithContext(ctx).Where(
			tx.SerialNoTransaction.BatchBundleID.Eq(batch.ID),
		).Delete()
		if err != nil {
			return
		}
	}
	_, err = tx.SerialNo.WithContext(ctx).Where(serialQ.BatchBundleID.In(ids...)).UpdateSimple(
		serialQ.Status.Value(proto.State_ACTIVE.String()),
	)
	return
}
func (r *serialNoEventRepo) OnStockEntrySubmitted(ctx context.Context, payload event.StatusStockEntryEventData) (err error) {
	tx := payload.Tx
	stockEntry := payload.StockEntry
	err = r.processStockEntry(tx, ctx, payload.LineItemsData.LineItems, stockEntry.CompanyID, stockEntry.PostingDate, stockEntry.PostingTime,
		proto.PartyType_stockEntry.String(), stockEntry.Code, stockEntry.EntryType)
	return
}

func (r *serialNoEventRepo) OnReceiptSubmitted(ctx context.Context, payload event.StatusReceiptEventData) (err error) {
	tx := payload.Tx
	receipt := payload.Receipt
	companyDefault := payload.CompanyDefault
	exchangeRate, err := r.accounting.GetCurrencyExchangeRate(ctx, tx, companyDefault,payload.Receipt.Currency, false, true)
	if err != nil {
		return
	}
	switch payload.ReceiptPartyType {
	case proto.PartyType_purchaseReceipt.String():
		err = r.receiptStock(tx, ctx, payload.LineItemsData.LineItems, receipt.CompanyID,
			receipt.PostingDate, receipt.PostingTime, payload.ReceiptPartyType, receipt.Code,exchangeRate)
	case proto.PartyType_deliveryNote.String():
		err = r.processOutStock(tx, ctx, payload.LineItemsData.LineItems, receipt.CompanyID,
			receipt.PostingDate, receipt.PostingTime, payload.ReceiptPartyType, receipt.Code)
	}
	return
}

func (r *serialNoEventRepo) OnInvoiceSubmitted(ctx context.Context, payload event.StatusInvoiceEventData) (err error) {
	if !payload.Invoice.UpdateStock {
		return
	}
	tx := payload.Tx
	invoice := payload.Invoice
	purchaseInvoice := proto.PartyType_purchaseInvoice.String()
	saleInvoice := proto.PartyType_saleInvoice.String()
	companyDefault := payload.CompanyDefaults
	exchangeRate, err := r.accounting.GetCurrencyExchangeRate(ctx, tx, companyDefault,invoice.Currency, false, true)
	if err != nil {
		return
	}
	switch payload.InvoicePartyType {
		
	case purchaseInvoice:
		err = r.receiptStock(tx, ctx, payload.LineItemsData.LineItems, invoice.CompanyID,
			invoice.PostingDate, invoice.PostingTime, purchaseInvoice, invoice.Code,exchangeRate)
	case saleInvoice:
		err = r.processOutStock(tx, ctx, payload.LineItemsData.LineItems, invoice.CompanyID,
			invoice.PostingDate, invoice.PostingTime, saleInvoice, invoice.Code)
	}
	return
}
func (r *serialNoEventRepo) processOutStock(tx *query.QueryTx, ctx context.Context, lines []dto.LineItem,
	companyID int64,  postingDate time.Time, postingTime string,
	voucherType string, voucherCode string) (err error) {
		fmt.Println("processOutStock...")
	for _, line := range lines {
		err = r.onOutSerialNos(ctx, tx, line, companyID, line.SourceWarehouseID, postingDate,
			postingTime, voucherType, voucherCode)
		if err != nil {
			return err
		}
	}
	return
}
func (r *serialNoEventRepo) processStockEntry(tx *query.QueryTx, ctx context.Context, lines []dto.LineItem,
	companyID int64, postingDate time.Time, postingTime string,
	voucherType string, voucherCode string, stockEntryType string) (err error) {
	for _, line := range lines {
		if !line.MaintainStock {
			return
		}
		switch stockEntryType {
		case proto.StockEntryType_MATERIAL_RECEIPT.String():
			if line.TargetWarehouseID == 0 {
				return domain.BLANK_VALUE
			}
			err = r.createSerialNos(ctx, tx, line, companyID, line.TargetWarehouseID, postingDate,
				postingTime, line.Rate, voucherType, voucherCode)
		}

		if err != nil {
			return err
		}
	}
	return
}

func (r *serialNoEventRepo) receiptStock(tx *query.QueryTx, ctx context.Context, lines []dto.LineItem,
	companyID int64, postingDate time.Time, postingTime string,
	voucherType string, voucherCode string,exchangeRate int32) (err error) {
	
	for _, line := range lines {
		if !line.MaintainStock {
			return
		}
		rate := r.currency.CurrencyExchange(line.Rate, exchangeRate)

		err = r.createSerialNos(ctx, tx, line, companyID, line.AcceptedWarehouse, postingDate,
			postingTime, rate, voucherType, voucherCode)
		if err != nil {
			return err
		}
	}
	return
}
func (r *serialNoEventRepo) onOutSerialNos(ctx context.Context, tx *query.QueryTx, line dto.LineItem,
	companyId int64, warehouseID int64, postingDate time.Time, postingTime string,
	voucherType string, voucherCode string) (err error) {
	batchBundle, err := r.createBatchBundle(ctx, tx, companyId, warehouseID, line.ItemID, postingDate, postingTime,
		voucherType, voucherCode)
	if err != nil {
		return
	}
	batchBundleQ := tx.BatchBundle
	serialNoQ := tx.SerialNo
	snTransactions := make([]*model.SerialNoTransaction, line.Quantity)
	for i := 0; i < int(line.Quantity); i++ {
		serialNo, err := tx.SerialNo.WithContext(ctx).
			Join(batchBundleQ,
				batchBundleQ.ID.EqCol(serialNoQ.BatchBundleID),
				batchBundleQ.WarehouseID.Eq(warehouseID),
			).Where(
			serialNoQ.Status.Eq(proto.State_ACTIVE.String()),
			serialNoQ.ItemID.Eq(line.ItemID),
		).First()
		if err != nil {
			return domain.ERROR_OUT_OF_STOCK
		}
		_, err = tx.SerialNo.WithContext(ctx).Where(serialNoQ.ID.Eq(serialNo.ID)).UpdateSimple(
			serialNoQ.Status.Value(proto.State_DELIVERED.String()),
			serialNoQ.BatchBundleID.Value(batchBundle.ID),
		)
		if err != nil {
			return err
		}
		snTransaction := r.serialNoTransaction(serialNo.ID, -1, proto.State_DELIVERED.String(),
			batchBundle.ID)
		snTransactions[i] = &snTransaction
	}
	err = tx.SerialNoTransaction.WithContext(ctx).CreateInBatches(snTransactions, len(snTransactions))
	if err != nil {
		return err
	}
	return
}
func (r *serialNoEventRepo) createSerialNos(ctx context.Context, tx *query.QueryTx, line dto.LineItem,
	companyId int64, warehouseID int64, postingDate time.Time, postingTime string,
	rate int64, voucherType string, voucherCode string) (err error) {
	// snReference *model.SerialNoReference,serialNos []*model.SerialNo,
	fmt.Println("LINE ITEM", line)
	batchBundle, err := r.createBatchBundle(ctx, tx, companyId, warehouseID, line.ItemID, postingDate, postingTime,
		voucherType, voucherCode)
	if err != nil {
		return
	}
	serialNos := make([]*model.SerialNo, line.Quantity)
	snTransactions := make([]*model.SerialNoTransaction, line.Quantity)
	snList, err := r.countItemSerialNumbers(ctx, tx, line.ItemID, line.Quantity)
	if err != nil {
		return err
	}
	for i := 0; i < int(line.Quantity); i++ {
		serialNo := &model.SerialNo{}
		serialNoID, err := tx.SerialNo.InsertParty(proto.PartyType_serialNo.String())
		if err != nil {
			return err
		}
		serialNo.BatchBundleID = batchBundle.ID
		serialNo.SerialNo = snList[i]
		serialNo.Status = proto.State_ACTIVE.String()
		serialNo.ID = serialNoID
		serialNo.ItemID = line.ItemID
		serialNo.ValuationRate = rate
		serialNos[i] = serialNo

		//Insert serial no transactionEntry
		snTransaction := r.serialNoTransaction(serialNoID, 1, proto.State_ACTIVE.String(),
			batchBundle.ID)
		snTransactions[i] = &snTransaction
	}
	err = tx.SerialNo.WithContext(ctx).CreateInBatches(serialNos, len(serialNos))
	if err != nil {
		return err
	}
	err = tx.SerialNoTransaction.WithContext(ctx).CreateInBatches(snTransactions, len(snTransactions))
	return
}

func (e *serialNoEventRepo) countItemSerialNumbers(ctx context.Context, tx *query.QueryTx,
	itemID int64, quantity int32) (snList []string, err error) {
	count, err := tx.SerialNo.WithContext(ctx).Unscoped().Count()
	// count, err := tx.SerialNo.WithContext(ctx).Where(
	// 	tx.SerialNo.ItemID.Eq(itemID),
	// ).Count()
	if err != nil {
		return
	}
	fmt.Println("SERIAL NO COUNT", count, itemID)
	snList, err = e.generator.GenerateSN(SN_TEMPLATE, int(count), int(quantity))
	if err != nil {
		return
	}
	return
}

func (r *serialNoEventRepo) serialNoTransaction(snID int64, qty int16, status string, batchBundleID int64) model.SerialNoTransaction {
	return model.SerialNoTransaction{
		SerialNoID:    snID,
		Qty:           qty,
		Status:        status,
		BatchBundleID: batchBundleID,
	}
}

func (r *serialNoEventRepo) createBatchBundle(ctx context.Context, tx *query.QueryTx,
	companyId int64, warehouseID int64, itemID int64, postingDate time.Time, postingTime string,
	voucherType string, voucherCode string) (batchBundle model.BatchBundle, err error) {
	count, err := tx.BatchBundle.WithContext(ctx).Where(
		tx.BatchBundle.CompanyID.Eq(companyId),
	).Count()
	if err != nil {
		return
	}
	snList, err := r.generator.GenerateSN(BATCH_BUNDLE_TEMPLATE, int(count), 1)
	if err != nil {
		return
	}
	if len(snList) == 0 {
		return batchBundle, domain.ERROR_EMPTY_ARRAY
	}
	batchBundleID, err := tx.BatchBundle.InsertParty(proto.PartyType_batchBundle.String())
	if err != nil {
		return
	}
	batchBundle.ID = batchBundleID
	batchBundle.CompanyID = companyId
	batchBundle.WarehouseID = warehouseID
	batchBundle.ItemID = itemID
	batchBundle.BatchBundleNo = snList[0]
	batchBundle.PostingDate = postingDate
	batchBundle.PostingTime = postingTime
	batchBundle.VoucherCode = voucherCode
	batchBundle.VoucherType = voucherType
	err = tx.BatchBundle.WithContext(ctx).Save(&batchBundle)
	if err != nil {
		return
	}
	return
}

// func (r *serialNoEventRepo)OnInvoice
