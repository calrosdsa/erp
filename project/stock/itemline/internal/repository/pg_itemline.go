package itemline_repo

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"
	"fmt"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

type ItemLineRepository interface {
	EditItemLine(req *common.RequestContext, i *dto.EditLineItemRequest) error
	GetItemLines(req *common.RequestContext, d *dto.RequestLineItems) (
		res []dto.LineItemDto, err error)
	DeleteLineItem(req *common.RequestContext, d *dto.DeleteLineItemRequest) (err error)
	DeleteLineItems(ctx context.Context, tx *query.QueryTx, docPartyID int64) (err error)
	AddLineItem(req *common.RequestContext, d *dto.AddLineItemRequest) (err error)
	CreateLineItem(ctx context.Context, tx *query.QueryTx, line dto.LineItemData, docPartyID int64,
		docPartyType string, updateStock bool) (err error)
	UpsertProductList(req *common.RequestContext, dto dto.ProductListData) (err error)
}

type itemLineRepository struct {
	conn      db.Connection
	Q         *query.Query
	currency  helpers.CurrencyHelper
	convertor helpers.ConvertorHelper
}

func NewItemLineRepository(
	conn db.Connection,
	helpers *helpers.Helpers,
) ItemLineRepository {
	return &itemLineRepository{
		conn:      conn,
		Q:         conn.GetQ(),
		currency:  helpers.Currency,
		convertor: helpers.Convertor,
	}
}

func (r *itemLineRepository) UpsertProductList(req *common.RequestContext, d dto.ProductListData) (err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	err = r.DeleteLineItems(req.Ctx, tx, d.PartyID)
	if err != nil {
		return
	}
	for _, line := range d.Lines {
		err = r.CreateLineItem(req.Ctx, tx, line, d.PartyID, d.PartyType, false)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()

	return
}

func (r *itemLineRepository) GetItemLines(req *common.RequestContext, d *dto.RequestLineItems) (
	res []dto.LineItemDto, err error) {
	var (
		conds   []gen.Condition
		columns []field.Expr
	)
	res = []dto.LineItemDto{}
	itemLineQ := r.Q.ItemLine
	itemQ := r.Q.Item
	warehouseQ := r.Q.WareHouse.As("w1")
	warehouse2Q := r.Q.WareHouse.As("w2")
	uomQ := r.Q.UnitOfMeasure
	builder := r.Q.WithContext(req.Ctx).ItemLine

	//Default Columns
	columns = append(columns, itemLineQ.ID, itemLineQ.Rate, itemLineQ.Quantity,
		itemLineQ.ItemLineReferenceID,
		itemLineQ.ItemID, itemLineQ.UnitOfMeasureID, itemLineQ.Type.As("line_type"),
		itemQ.Name.As("item_name"), itemQ.Code.As("item_code"),itemQ.Description.As("item_description"),
		uomQ.Code.As("uom"))

	conds = append(conds, itemLineQ.PartyID.Eq(r.convertor.StrtoInt(d.ID)))
	builder = builder.
		Join(itemQ, itemQ.ID.EqCol(itemLineQ.ItemID)).
		Join(uomQ, uomQ.ID.EqCol(itemLineQ.UnitOfMeasureID))
	if proto.ItemLineType_ITEM_LINE_RECEIPT.String() == d.LineType ||
		(r.convertor.StrToBool(d.UpdateStock) && proto.PartyType_purchaseInvoice.String() == d.PartyType) {
		ilReceiptQ := r.Q.ItemLineReceipt
		columns = append(columns, ilReceiptQ.AcceptedQuantity, ilReceiptQ.RejectedQuantity,
			warehouseQ.Name.As("accepted_warehouse"), ilReceiptQ.AcceptedWarehouse.As("accepted_warehouse_id"),
			warehouse2Q.Name.As("rejected_warehouse"), ilReceiptQ.RejectedWarehouse.As("rejected_warehouse_id"),
		)
		builder = builder.Join(ilReceiptQ, ilReceiptQ.ItemLine.EqCol(itemLineQ.ID)).
			Join(warehouseQ, warehouseQ.ID.EqCol(ilReceiptQ.AcceptedWarehouse)).
			LeftJoin(warehouse2Q, warehouse2Q.ID.EqCol(ilReceiptQ.RejectedWarehouse))
	}

	if proto.ItemLineType_DELIVERY_LINE_ITEM.String() == d.LineType ||
		(r.convertor.StrToBool(d.UpdateStock) && proto.PartyType_saleInvoice.String() == d.PartyType) {
		diLineQ := r.Q.DeliveryLineItem
		columns = append(columns,
			warehouseQ.Name.As("source_warehouse"), warehouseQ.ID.As("source_warehouse_id"),
		)
		builder = builder.Join(diLineQ, diLineQ.ItemLineID.EqCol(itemLineQ.ID)).
			Join(warehouseQ, warehouseQ.ID.EqCol(diLineQ.SourceWarehouseID))
	}

	err = builder.Select(columns...).Where(conds...).Scan(&res)
	return
}

func (r *itemLineRepository) AddLineItem(req *common.RequestContext, d *dto.AddLineItemRequest) (err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	err = r.CreateLineItem(req.Ctx, tx, d.Body.LineItemData, d.Body.DocPartyID, d.Body.DocPartyType, d.Body.UpdateStock)
	if err != nil {
		return
	}
	err = r.updateDocTotalAmount(req, tx, d.Body.DocPartyID, d.Body.DocPartyType, d.Body.TotalAmount, d.Body.TotalItems)
	if err != nil {
		return
	}
	err = r.adjustTaxAndChargeAmount(req.Ctx, tx, d.Body.TotalAmountItems, d.Body.Charges)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (r *itemLineRepository) EditItemLine(req *common.RequestContext, d *dto.EditLineItemRequest) (err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	lineItemQ := tx.ItemLine
	lineItemData := d.Body.LineItemData
	var columns []field.AssignExpr
	columns = append(columns, lineItemQ.Rate.Value(r.currency.FloatToInt(lineItemData.Rate)),
		lineItemQ.Quantity.Value(lineItemData.Quantity), lineItemQ.ItemID.Value(lineItemData.ItemID))
	_, err = tx.ItemLine.WithContext(req.Ctx).Where(
		lineItemQ.ID.Eq(int32(d.Body.ID)),
	).UpdateSimple(columns...)
	if err != nil {
		return
	}
	fmt.Println("TOTAL ITEM", d.Body.TotalItems)
	err = r.updateDocTotalAmount(req, tx, d.Body.DocPartyID, d.Body.DocPartyType, d.Body.TotalAmount, d.Body.TotalItems)
	if err != nil {
		return
	}
	err = r.adjustTaxAndChargeAmount(req.Ctx, tx, d.Body.TotalAmountItems, d.Body.Charges)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (r *itemLineRepository) DeleteLineItem(req *common.RequestContext, d *dto.DeleteLineItemRequest) (err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	_, err = r.Q.ItemLine.WithContext(req.Ctx).Where(
		r.Q.ItemLine.ID.Eq(d.Body.ID),
	).Delete()
	err = r.updateDocTotalAmount(req, tx, d.Body.DocPartyID, d.Body.DocPartyType, d.Body.TotalAmount, d.Body.TotalItems)
	if err != nil {
		return
	}
	err = r.adjustTaxAndChargeAmount(req.Ctx, tx, d.Body.TotalAmountItems, d.Body.Charges)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (r *itemLineRepository) adjustTaxAndChargeAmount(
	ctx context.Context,
	tx *query.QueryTx,
	totalAmountItems float64,
	charges []dto.TaxAndChargeLineDto,
) (err error) {
	fmt.Println("TOTAL AMOUNT:", totalAmountItems)

	// Iterate through each charge and apply adjustments
	for _, charge := range charges {
		if charge.Type == proto.TaxChargeLineType_ON_NET_TOTAL.String() {
			// Calculate the amount based on the tax rate and total amount of items
			amount := totalAmountItems * float64(charge.TaxRate) / 100
			// Update the tax and charge line amount in the database
			_, err = tx.WithContext(ctx).TaxAndChargeLine.
				Where(tx.TaxAndChargeLine.ID.Eq(int32(charge.ID))).
				UpdateSimple(tx.TaxAndChargeLine.Amount.Value(r.currency.FloatToInt(amount)))
			if err != nil {
				return
			}
		}
	}

	return nil
}
func (r *itemLineRepository) updateDocTotalAmount(req *common.RequestContext,
	tx *query.QueryTx, docPartyID int64, docPartyType string, totalAmount float64, totalItems int32) (err error) {
	switch docPartyType {
	case proto.PartyType_purchaseOrder.String(), proto.PartyType_saleOrder.String():
		_, err = tx.ProgressOrder.WithContext(req.Ctx).Where(
			tx.ProgressOrder.OrderID.Eq(docPartyID),
		).UpdateSimple(
			tx.ProgressOrder.TotalAmount.Value(r.currency.FloatToInt(totalAmount)),
			tx.ProgressOrder.TotalItems.Value(totalItems),
		)
	case proto.PartyType_purchaseInvoice.String(), proto.PartyType_saleInvoice.String():
		_, err = tx.ProgressInvoice.WithContext(req.Ctx).Where(
			tx.ProgressInvoice.InvoiceID.Eq(docPartyID),
		).UpdateSimple(
			tx.ProgressInvoice.TotalAmount.Value(r.currency.FloatToInt(totalAmount)),
		)
	}
	return err
}

func (r *itemLineRepository) DeleteLineItems(ctx context.Context, tx *query.QueryTx, docPartyID int64) (err error) {
	_, err = tx.WithContext(ctx).ItemLine.Unscoped().Where(
		tx.ItemLine.PartyID.Eq(docPartyID),
	).Delete()
	// tx.Line
	return
}

func (r *itemLineRepository) CreateLineItem(ctx context.Context, tx *query.QueryTx, line dto.LineItemData, docPartyID int64,
	docPartyType string, updateStock bool) (err error) {
	itemLine := model.ItemLine{}
	itemLine.ItemID = line.ItemID
	itemLine.PartyID = &docPartyID
	itemLine.Quantity = line.Quantity
	itemLine.UnitOfMeasureID = line.UnitOfMeasureID
	itemLine.Rate = r.currency.FloatToInt(line.Rate)
	itemLine.ItemLineReferenceID = line.ItemLineReferenceID
	itemLine.Type = line.LineType

	err = tx.ItemLine.WithContext(ctx).Save(&itemLine)
	if err != nil {
		return err
	}

	//Creating invoiced itemline
	if line.LineType == proto.ItemLineType_ITEM_LINE_INVOICE.String() {
		invoicedItemLine := &model.InvoicedItemLine{}
		invoicedItemLine.ItemLine = itemLine.ID
		invoicedItemLine.InvoiceID = docPartyID
		err = tx.InvoicedItemLine.WithContext(ctx).Save(invoicedItemLine)
		if err != nil {
			return err
		}
	}

	fmt.Println("ITEM LINE DATA", updateStock)
	fmt.Println("DOC PARTY", docPartyType)
	//Creating Recipt itemLine
	if line.LineType == proto.ItemLineType_ITEM_LINE_RECEIPT.String() ||
		(updateStock && docPartyType == proto.PartyType_purchaseInvoice.String()) {
		fmt.Println("CREATING LINE RECEIPT...")
		itemLineReceipt := &model.ItemLineReceipt{}
		itemLineReceipt.ItemLine = itemLine.ID
		itemLineReceipt.AcceptedQuantity = line.Quantity
		// itemLineReceipt.AcceptedQuantity = line.LineItemReceipt.AcceptedQuantity
		// itemLineReceipt.RejectedQuantity = line.LineItemReceipt.RejectedQuantity
		itemLineReceipt.AcceptedWarehouse = line.LineItemReceipt.AcceptedWarehouse
		if line.LineItemReceipt.RejectedWarehouse != 0 {
			itemLineReceipt.RejectedWarehouse = &line.LineItemReceipt.RejectedWarehouse
		}
		err = tx.ItemLineReceipt.WithContext(ctx).Save(itemLineReceipt)
		if err != nil {
			return err
		}
	}

	if line.LineType == proto.ItemLineType_DELIVERY_LINE_ITEM.String() ||
		(updateStock && docPartyType == proto.PartyType_saleInvoice.String()) {
		deliveryLineItem := &model.DeliveryLineItem{}
		deliveryLineItem.ItemLineID = itemLine.ID
		deliveryLineItem.SourceWarehouseID = line.DeliveryLineItem.SourceWarehouse
		err = tx.DeliveryLineItem.WithContext(ctx).Save(deliveryLineItem)
		if err != nil {
			return err
		}
	}

	if line.LineType == proto.ItemLineType_ITEM_LINE_STOCK_ENTRY.String() {
		itemLineStockEntry := &model.ItemLineStockEntry{}
		itemLineStockEntry.ItemLine = itemLine.ID
		if line.LineItemStockEntry.SourceWarehouse != 0 {
			itemLineStockEntry.SourceWarehouseID = &line.LineItemStockEntry.SourceWarehouse
		}
		if line.LineItemStockEntry.TargetWarehouse != 0 {
			itemLineStockEntry.TargetWarehouseID = &line.LineItemStockEntry.TargetWarehouse
		}
		err = tx.ItemLineStockEntry.WithContext(ctx).Save(itemLineStockEntry)
		if err != nil {
			return err
		}
	}
	return
}
