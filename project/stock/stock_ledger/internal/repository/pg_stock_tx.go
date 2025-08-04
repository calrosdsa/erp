package stock_ledger_repo

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"
	"fmt"
	"strings"

	"gorm.io/gen/helper"
	"gorm.io/gorm"
)

type StockLedgerTxRepository interface {
	GetBalanceItem(ctx context.Context, tx *query.QueryTx, args ...interface{}) (
		res StockBalance, err error)
	GetStockLedgerReport(req *common.RequestContext, i *dto.RequestStockLedger) (
		res []dto.StockLedgerEntryDto, err error)
	GetStockBalanceReport(req *common.RequestContext, i *dto.RequestStockBalance) (
		res []dto.StockBalanceEntryDto, err error)
	// GetItemLineReceipt(ctx context.Context, tx *query.QueryTx, receiptID int64) (
	// 	res []ItemLineReceipt, err error)
	OnCancelled(tx *query.QueryTx, ctx context.Context, voucherCode string) (err error)
}

type stockLedgerTxRepo struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
}

func NewStockLedgeTX(
	db db.Connection,
	helpers *helpers.Helpers,
) StockLedgerTxRepository {
	return &stockLedgerTxRepo{
		Q: db.GetQ(),

		convertor: helpers.Convertor,
	}
}

func (r *stockLedgerTxRepo) OnCancelled(tx *query.QueryTx, ctx context.Context, voucherCode string) (err error) {
	stockTxQ := tx.StockTransaction
	_, err = stockTxQ.WithContext(ctx).Where(
		stockTxQ.VoucherNo.Eq(voucherCode),
	).Delete()
	if err != nil {
		return
	}
	return
}

func (r *stockLedgerTxRepo) GetStockBalanceReport(req *common.RequestContext, i *dto.RequestStockBalance) (
	res []dto.StockBalanceEntryDto, err error) {
	var params []interface{}
	var generateSQL strings.Builder
	generateSQL.WriteString(`
	SELECT
		tx.posting_date AS date,
		it.name AS item_name,
		it.uuid AS item_uuid,
		it.id AS item_id,
		it.code AS item_code,
		uom.code AS stock_uom,
		t.in_qty,
		t.out_qty,
		tx.balance_qty,
		w.name AS warehouse_name,
		w.uuid AS warehouse_uuid,
		tx.average_rate,
		tx.balance_value,
		tx.currency
	FROM stock_transactions AS tx
	LEFT JOIN (
		SELECT
			item_id,
			warehouse_id,
			SUM(in_qty) AS in_qty,
			SUM(out_qty) AS out_qty
		FROM stock_transactions
		GROUP BY item_id, warehouse_id
	) AS t ON tx.item_id = t.item_id AND tx.warehouse_id = t.warehouse_id
	INNER JOIN items AS it ON it.id = tx.item_id AND it.company_id = ?
	INNER JOIN unit_of_measures AS uom ON uom.id = tx.uom_id
	INNER JOIN ware_houses AS w ON w.id = tx.warehouse_id
	JOIN (
		SELECT
			item_id,
			warehouse_id,
			MAX(id) AS latest_id
		FROM stock_transactions
		GROUP BY item_id, warehouse_id
	) AS latest_tx ON tx.id = latest_tx.latest_id
	`)
	var whereSQL0 strings.Builder
	params = append(params, req.ActiveCompany.ID)
	if i.FromDate != "" && i.ToDate != "" {
		params = append(params, i.FromDate)
		params = append(params, i.ToDate)
		whereSQL0.WriteString("tx.posting_date between ? and ? ")
	}
	if i.ItemID != "" {
		itemID := r.convertor.StrtoInt(i.ItemID)
		params = append(params, itemID)
		whereSQL0.WriteString("and tx.item_id = ? ")
	}
	if i.WarehouseID != "" {
		warehouseID := r.convertor.StrtoInt(i.WarehouseID)
		params = append(params, warehouseID)
		whereSQL0.WriteString("and tx.warehouse_id = ? ")
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL0)

	err = r.Q.StockTransaction.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}

func (r *stockLedgerTxRepo) GetStockLedgerReport(req *common.RequestContext, i *dto.RequestStockLedger) (
	res []dto.StockLedgerEntryDto, err error) {
	var params []interface{}
	var generateSQL strings.Builder
	generateSQL.WriteString(`
	select 
		tx.posting_date as date,
		it.name as item_name,it.uuid as item_uuid,
		it.id as item_id,it.code AS item_code,
		uom.code as stock_uom,
		tx.in_qty,tx.out_qty,tx.balance_qty,
		w.name as warehouse_name,w.uuid as warehouse_uuid,
		tx.incoming_rate,tx.average_rate,tx.valuation_rate,tx.balance_value,
		tx.voucher_type,tx.voucher_no,tx.currency
		from stock_transactions as tx
		inner join items as it on it.id = tx.item_id and it.company_id = ?
		inner join unit_of_measures as uom on uom.id = tx.uom_id
		inner join ware_houses as w on w.id = tx.warehouse_id
	`)
	var whereSQL0 strings.Builder
	params = append(params, req.ActiveCompany.ID)
	params = append(params, i.FromDate)
	params = append(params, i.ToDate)
	whereSQL0.WriteString("tx.posting_date between ? and ? ")

	if i.ItemID != "" {
		itemID := r.convertor.StrtoInt(i.ItemID)
		params = append(params, itemID)
		whereSQL0.WriteString("and tx.item_id = ? ")
	}
	if i.WarehouseID != "" {
		warehouseID := r.convertor.StrtoInt(i.WarehouseID)
		params = append(params, warehouseID)
		whereSQL0.WriteString("and tx.warehouse_id = ? ")
	}
	if i.VoucherNo != "" {
		params = append(params, i.VoucherNo)
		whereSQL0.WriteString("and tx.voucher_no = ? ")
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL0)
	err = r.Q.StockTransaction.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}

func (r *stockLedgerTxRepo) GetBalanceItem(ctx context.Context, tx *query.QueryTx, args ...interface{}) (res StockBalance, err error) {
	itemID := args[0].(int64)
	warehouseID := args[1].(int64)
	incomingRate := args[2].(int64)
	in_qty := args[3].(int32)
	stockTx := tx.StockTransaction
	builder := stockTx.WithContext(ctx).Where(
		stockTx.ItemID.Eq(itemID),
		stockTx.WarehouseID.Eq(warehouseID),
	)
	count, err := builder.Count()
	if err != nil {
		fmt.Println("GET BALANCE ITEM", err)
	}
	//get last register
	stockTransaction, err := builder.
		Select(stockTx.BalanceQty, stockTx.BalanceValue, stockTx.AverageRate).
		Last()
	if err == gorm.ErrRecordNotFound {
		err = nil
		stockTransaction = &model.StockTransaction{}
	}
	if err != nil {
		return
	}
	fmt.Println("COUNT", count)
	fmt.Println("STOCK TX", stockTransaction)
	res.AvgRate = ((stockTransaction.AverageRate * count) + incomingRate) / (count + 1)
	res.BalanceQuantity = stockTransaction.BalanceQty + in_qty
	res.BalanceRate = stockTransaction.BalanceValue + (int64(in_qty) * incomingRate)

	return
}

// func (r *stockLedgerTxRepo) GetItemLineReceipt(ctx context.Context, tx *query.QueryTx, receiptID int64) (
// 	res []ItemLineReceipt, err error) {
// 	itemLineQ := tx.ItemLine
// 	itemLineReceiptQ := tx.ItemLineReceipt
// 	err = tx.ItemLine.WithContext(ctx).Select(
// 		itemLineQ.Rate, itemLineQ.ItemPriceID,
// 		itemLineReceiptQ.AcceptedQuantity, itemLineReceiptQ.RejectedQuantity,
// 		itemLineReceiptQ.AcceptedWarehouse,
// 	).Where(
// 		itemLineQ.PartyID.Eq(receiptID),
// 	).Join(
// 		itemLineReceiptQ, itemLineQ.ID.EqCol(itemLineReceiptQ.ItemLine),
// 	).Scan(&res)
// 	return
// }

type StockBalance struct {
	BalanceQuantity int32 `json:"balance_quantity"` // Etiquetas JSON, por si se usan en una API REST
	BalanceRate     int64 `json:"balance_rate"`
	AvgRate         int64 `json:"avg_rate"`
}
