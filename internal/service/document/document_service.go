package document_service

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/repository"
	"erp/pkg/db"
	"erp/pkg/exporter"
	"fmt"

	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/signature"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontfamily"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"gorm.io/gen/field"
)

type documentService struct {
	Q         *query.Query
	core      repository.CoreService
	currency  helpers.CurrencyHelper
	dFontSize float64
	dPad      float64
}

func NewDocumentService(
	conn db.Connection,
	currency helpers.CurrencyHelper,
	core repository.CoreService,
) repository.DocumentService {
	return &documentService{
		currency:  currency,
		dFontSize: 8.3,
		dPad:      1.0,
		core:      core,
		Q:         conn.GetQ(),
	}
}


func (s *documentService) GetTaxLines(q *query.Query, ctx context.Context, id int64) (res dto.TaxLinesData, err error) {
	ledgerQ := q.Ledger
	tacQ := q.TaxAndChargeLine
	err = tacQ.WithContext(ctx).
		Select(tacQ.TaxRate, tacQ.IsDeducted, tacQ.AccountHead, tacQ.Amount,
			ledgerQ.AccountRootType.As("acct_head_root_type"), ledgerQ.IsOffsetAccount).
		Join(ledgerQ, ledgerQ.ID.EqCol(tacQ.AccountHead)).
		Where(
			tacQ.DocPartyID.Eq(id),
		).Scan(&res.TaxLines)
	var (
		totalAmount int64
	)
	for _, line := range res.TaxLines {
		totalAmount += line.Amount
	}
	res.TotalAmount = totalAmount
	return
}

func (s *documentService) GetLineItems(q *query.Query, req *common.RequestContext,
	id int64, opts ...repository.OptionStock) (res dto.LineItemsData, err error) {
	fmt.Println("OPTIONS", opts)
	o := repository.OptionsStock.Apply(opts...)
	var columns []field.Expr
	itemLineQ := q.ItemLine

	columns = append(columns, itemLineQ.ID, itemLineQ.Rate, itemLineQ.Quantity,
		itemLineQ.ItemID, itemLineQ.UnitOfMeasureID)

	builder := q.ItemLine.WithContext(req.Ctx)
	if o.LoadItemInLine {
		itemQ := q.Item
		builder = builder.Join(itemQ, itemQ.ID.EqCol(itemLineQ.ItemID))
		columns = append(columns, itemQ.MaintainStock,itemQ.Name.As("item_name"),
		itemQ.Description.As("item_description"))
	}
	if o.LoadDeliveryLineItem {
		fmt.Println("LOAD DELIVERY LINE ITEM...")
		deliveryLIQ := q.DeliveryLineItem
		builder = builder.Join(deliveryLIQ, deliveryLIQ.ItemLineID.EqCol(itemLineQ.ID))
		columns = append(columns, deliveryLIQ.SourceWarehouseID)
	}
	if o.LoadReceiptLineItem {
		fmt.Println("LOAD LINE RECEIPT...")
		receiptLIQ := q.ItemLineReceipt
		builder = builder.Join(receiptLIQ, receiptLIQ.ItemLine.EqCol(itemLineQ.ID))
		columns = append(columns, receiptLIQ.AcceptedWarehouse, receiptLIQ.AcceptedQuantity)
	}
	if o.LoadLineStockEntry {
		lineStockEntry := q.ItemLineStockEntry
		builder = builder.Join(lineStockEntry, lineStockEntry.ItemLine.EqCol(itemLineQ.ID))
		columns = append(columns, lineStockEntry.SourceWarehouseID, lineStockEntry.TargetWarehouseID)
	}
	err = builder.Select(columns...).Where(
		itemLineQ.PartyID.Eq(id),
	).Scan(&res.LineItems)

	var (
		totalAmount   int64
		totalQuantity int32
	)
	for _, line := range res.LineItems {
		totalAmount += int64(line.Quantity) * line.Rate
		totalQuantity += line.Quantity
	}
	res.TotalAmount = int64(totalAmount)
	res.TotalQuantity = totalQuantity
	return
}


func (s *documentService) GenerateReceiptDocumentPdf(req *common.RequestContext, d dto.ReceiptDto, docPartyType string) (
	res []byte, err error) {
	lineItems, taxLines, err := s.getItemsAndCharges(req, d.ID)
	if err != nil {
		return
	}
	addressAndContact, err := s.GetAddressAndContact(req, d.ID)
	if err != nil {
		return
	}
	e := exporter.NewPdfExporter()
	b := exporter.NewIPdfBuilder()
	e.SetBuilder(b)
	var (
		companyAddress dto.AddressDto
	)
	if req.CompanyDefaults.AddressID != nil {
		companyAddress, _ = s.core.GetAddress(req, *req.CompanyDefaults.AddressID)
	}
	docInfo := DocInfo{
		LogoUrl: req.ActiveCompany.Logo,
		Date:    d.PostingDate,
		DocNo:   d.Code,
		Address: companyAddress,
	}
	addr1 := AddressSection{}
	addr2 := AddressSection{}
	if docPartyType == proto.PartyType_purchaseReceipt.String() {
		addr1.Section = "Vendido a"
		addr1.Address = addressAndContact.ShippingAddress
		docInfo.Title = "Recibo de Compra"
	}
	if docPartyType == proto.PartyType_deliveryNote.String() {
		addr1.Section = "Entregado a"
		addr1.Address = addressAndContact.PartyAddress
		docInfo.Title = "Nota de Entrega"
	}
	s.generateDocInfo(b, docInfo)

	s.generateAddressSection(b, addr1, &addr2)

	if docPartyType == proto.PartyType_purchaseOrder.String() {
		s.generatePurchaseOrder(b)
	}
	s.generateItemsAndChargesSection(b, lineItems, taxLines, d.Currency)

	s.generateSignatureSection(b)
	return e.BuildPdfFile()
}

func (s *documentService) GenerateInvoiceDocumentPdf(req *common.RequestContext, d dto.InvoiceDto, docPartyType string) (
	res []byte, err error) {
	lineItems, taxLines, err := s.getItemsAndCharges(req, d.ID)
	if err != nil {
		return
	}
	addressAndContact, err := s.GetAddressAndContact(req, d.ID)
	if err != nil {
		return
	}
	e := exporter.NewPdfExporter()
	b := exporter.NewIPdfBuilder()
	e.SetBuilder(b)
	var (
		companyAddress dto.AddressDto
	)
	if req.CompanyDefaults.AddressID != nil {
		companyAddress, _ = s.core.GetAddress(req, *req.CompanyDefaults.AddressID)
	}
	docInfo := DocInfo{
		LogoUrl: req.ActiveCompany.Logo,
		Date:    d.PostingDate,
		DocNo:   d.Code,
		Address: companyAddress,
	}
	addr1 := AddressSection{}
	addr2 := AddressSection{}
	addr1.Section = "Facturado desde"
	addr2.Section = "Facturado a"
	if docPartyType == proto.PartyType_purchaseInvoice.String() {
		// addr1.Address = dto.AddressDto{}
		addr1.Address = addressAndContact.PartyAddress
		addr2.Address = addressAndContact.BillingAddress
		docInfo.Title = "Factura de Compra"
	}
	if docPartyType == proto.PartyType_saleInvoice.String() {
		// addr1.Address = dto.AddressDto{}
		addr1.Address = addressAndContact.BillingAddress
		addr2.Address = addressAndContact.PartyAddress
		docInfo.Title = "Factura de Venta"
	}
	s.generateDocInfo(b, docInfo)

	s.generateAddressSection(b, addr1, &addr2)

	if docPartyType == proto.PartyType_purchaseOrder.String() {
		s.generatePurchaseOrder(b)
	}
	s.generateItemsAndChargesSection(b, lineItems, taxLines, d.Currency)

	s.generateSignatureSection(b)
	return e.BuildPdfFile()
}

func (s *documentService) GenerateOrderDocumentPdf(req *common.RequestContext, d dto.OrderDto,
	orderPartyType string) (
	res []byte, err error) {
	lineItems, taxLines, err := s.getItemsAndCharges(req, d.ID)
	if err != nil {
		return
	}
	addressAndContact, err := s.GetAddressAndContact(req, d.ID)
	if err != nil {
		return
	}

	e := exporter.NewPdfExporter()
	b := exporter.NewIPdfBuilder()
	e.SetBuilder(b)
	var (
		companyAddress dto.AddressDto
	)
	if req.CompanyDefaults.AddressID != nil {
		companyAddress, _ = s.core.GetAddress(req, *req.CompanyDefaults.AddressID)
	}
	docInfo := DocInfo{
		LogoUrl: req.ActiveCompany.Logo,
		Date:    d.PostingDate,
		DocNo:   d.Code,
		Address: companyAddress,
	}
	addr1 := AddressSection{}
	addr2 := AddressSection{}

	if orderPartyType == proto.PartyType_purchaseOrder.String() {
		addr1.Section = "Vendedor"
		addr2.Section = "Enviar a"
		// addr1.Address = dto.AddressDto{}
		addr1.Address = addressAndContact.PartyAddress
		addr2.Address = addressAndContact.ShippingAddress
		docInfo.Title = "Orden de Compra"
	}

	if orderPartyType == proto.PartyType_saleOrder.String() {
		addr1.Section = "Facturar a"
		addr2.Section = "Enviar a"
		addr1.Address = addressAndContact.BillingAddress
		addr2.Address = addressAndContact.ShippingAddress
		docInfo.Title = "Orden de Venta"
	}
	s.generateDocInfo(b, docInfo)

	s.generateAddressSection(b, addr1, &addr2)

	if orderPartyType == proto.PartyType_purchaseOrder.String() {
		s.generatePurchaseOrder(b)
	}
	s.generateItemsAndChargesSection(b, lineItems, taxLines, d.Currency)

	s.generateSignatureSection(b)
	return e.BuildPdfFile()
}

func (s *documentService) generateSignatureSection(b exporter.IPdfBuilder) {
	b.SetRow(
		row.New(22).Add(
			col.New(6),
			signature.NewCol(3, "Authorizado por", props.Signature{FontFamily: fontfamily.Courier}),
			signature.NewCol(3, "Fecha", props.Signature{FontFamily: fontfamily.Courier}),
		),
	)
}




func (s *documentService) GetAddressAndContact(req *common.RequestContext, id int64) (
	res dto.AddressAndContact, err error) {
	aacQ := s.Q.AddressAndContact
	aac, err := aacQ.WithContext(req.Ctx).Where(
		aacQ.DocID.Eq(id),
	).First()
	if err != nil {
		return
	}
	if aac.ShippingAddressID != nil {
		res.ShippingAddress, err = s.core.GetAddress(req, *aac.ShippingAddressID)
		if err != nil {
			return
		}
	}
	if aac.BillingAddressID != nil {
		res.BillingAddress, err = s.core.GetAddress(req, *aac.BillingAddressID)
		if err != nil {
			return
		}
	}
	if aac.PartyAddressID != nil {
		res.PartyAddress, err = s.core.GetAddress(req, *aac.PartyAddressID)
		if err != nil {
			return
		}
	}
	if aac.ContactID != nil {
		res.PartyContact, err = s.core.GetContact(req, *aac.ContactID)
		if err != nil {
			return
		}
	}
	return
}

