package purchase_record_pdf

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/app/service/helpers"
	"erp/pkg/exporter"
	"fmt"
	"time"

	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/signature"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontfamily"
	"github.com/johnfercher/maroto/v2/pkg/consts/linestyle"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type PurchaseRecordPDF interface{
	GeneratePurchaseRecordDocument(req *common.RequestContext, d dto.PurchaseRecordDto) (res []byte, err error)
}

type purchaseRecordPDF struct {
	convertor helpers.ConvertorHelper
	currency  helpers.CurrencyHelper
	locale    helpers.Locale

	defaultFontSize   float64
	dfPad             float64
	docAttrText       props.Text
	defaultFontFamily string
	Q                 *query.Query
}

func NewPurchaseRecordPdf(
	helpers *helpers.Helpers,
	query *query.Query,
) PurchaseRecordPDF {
	defaultFontSize := 9.0
	dfPad := 1.0
	defaultFontFamily := fontfamily.Courier
	docAttrText := props.Text{
		Size: defaultFontSize,
		Top:  dfPad, Bottom: dfPad, Right: 1, Left: dfPad,
		Family: defaultFontFamily,
	}
	return &purchaseRecordPDF{
		convertor:         helpers.Convertor,
		currency:          helpers.Currency,
		locale:            helpers.Locale,
		defaultFontSize:   defaultFontSize,
		dfPad:             dfPad,
		docAttrText:       docAttrText,
		defaultFontFamily: defaultFontFamily,
		Q:                 query,
	}
}
func (p *purchaseRecordPDF) GeneratePurchaseRecordDocument(req *common.RequestContext, d dto.PurchaseRecordDto) (res []byte, err error) {
	e := exporter.NewPdfExporter()
	b := exporter.NewIPdfBuilder()
	e.SetBuilder(b)
	p.generateHeader(b, d, req.ActiveCompany)
	err = p.generateContent(b, d, req.CompanyDefaults)
	if err != nil {
		return
	}
	p.generateSignatureSection(b)
	return e.BuildPdfFile()
}


func (p *purchaseRecordPDF) generateContent(b exporter.IPdfBuilder, d dto.PurchaseRecordDto,
	companyDefaults model.CompanyDefault) (err error) {
	b.SetRow(
		row.New().Add(
			col.New(2).Add(
				text.New("Acreedor", p.docAttrText),
			),
			col.New(10).Add(
				text.New(fmt.Sprintf(": %s", d.Supplier), p.docAttrText),
			),
		),
	)

	b.SetRows(
		row.New(4),
		row.New().Add(
			col.New(2).Add(
				text.New("Factura", p.docAttrText),
			),
			col.New(3).Add(
				text.New(fmt.Sprintf(": %s", d.InvoiceNo), p.docAttrText),
			),
			col.New(1),
			col.New(2).Add(
				text.New("Total Fact", p.docAttrText),
			),
			col.New(3).Add(
				text.New(fmt.Sprintf(": %s", 
				p.currency.FormatCurrencyAmount(d.TotalPurchaseAmount, companyDefaults.Currency)), p.docAttrText),
			),
		),

		row.New().Add(
			col.New(2).Add(
				text.New("No. Aut.", p.docAttrText),
			),
			col.New(3).Add(
				text.New(fmt.Sprintf(": %s", d.AuthorizationCode), p.docAttrText),
			),
			col.New(1),
			col.New(2).Add(
				text.New("Impt. Dcto.", p.docAttrText),
			),
			col.New(3).Add(
				text.New(fmt.Sprintf(": %s", 
				p.currency.FormatCurrencyAmount(0, companyDefaults.Currency)), p.docAttrText),
			),
		),

		row.New().Add(
			col.New(6),
			col.New(2).Add(
				text.New("Impt. IE.", p.docAttrText),
			),
			col.New(3).Add(
				text.New(fmt.Sprintf(": %s", 
				p.currency.FormatCurrencyAmount(0, companyDefaults.Currency)), p.docAttrText),
			),
		),

		row.New().Add(
			col.New(2).Add(
				text.New("NIT.", p.docAttrText),
			),
			col.New(3).Add(
				text.New(fmt.Sprintf(": %s", d.SupplierNit), p.docAttrText),
			),
			col.New(1),
			col.New(2).Add(
				text.New("Impt. Exento", p.docAttrText),
			),
			col.New(3).Add(
				text.New(fmt.Sprintf(": %s", 
				p.currency.FormatCurrencyAmount(d.ExemptAmounts, companyDefaults.Currency)), p.docAttrText),
			),
		),

		row.New().Add(
			col.New(2).Add(
				text.New("Emision", p.docAttrText),
			),
			col.New(3).Add(
				text.New(fmt.Sprintf(": %s", d.InvoiceDuiDimDate.Format(time.DateOnly)), p.docAttrText),
			),
			col.New(1),
			col.New(2).Add(
				text.New("Impt. Neto", p.docAttrText),
			),
			col.New(3).Add(
				text.New(fmt.Sprintf(": %s", 
				p.currency.FormatCurrencyAmount(d.Subtotal, companyDefaults.Currency)), p.docAttrText),
			),
		),
		row.New().Add(
			col.New(2).Add(
				text.New("Proveedor", p.docAttrText),
			),
			col.New(3).Add(
				text.New(fmt.Sprintf(": %s", d.SupplierBusinessName), p.docAttrText),
			),
		
		),
	)

	return
}

func (p *purchaseRecordPDF) generateSignatureSection(b exporter.IPdfBuilder) {
	signatureProps := props.Signature{FontFamily: p.defaultFontFamily, LineStyle: linestyle.Dashed,
		FontSize: p.defaultFontSize}
	b.SetRow(
		row.New(40).Add(
			signature.NewCol(3, "MCC", signatureProps),
			col.New(1),
			signature.NewCol(3, "Vo. Bo.", signatureProps),
			col.New(1),
			signature.NewCol(3, "Interesado", signatureProps),
		),
	)
}

func (p *purchaseRecordPDF) generateHeader(b exporter.IPdfBuilder, d dto.PurchaseRecordDto,
	company model.Company) {
	b.RegisterHeader(
		row.New().Add(
			col.New(4).Add(
				text.New("REGISTRO DE COMPRA", p.docAttrText),
			),
			col.New(4).Add(
				text.New(company.Name, p.docAttrText),
			),
			col.New(1),
			col.New(3).Add(
				text.New(fmt.Sprintf("No Cta. : %s", "1351"), p.docAttrText),
			),
		),
		row.New().Add(
			col.New(4).Add(
				text.New(time.Now().Format(time.DateOnly), p.docAttrText),
			),
			col.New(4),
			// .Add(
			// 	text.New("REGISTRO CUENTA POR PAGAR", p.docAttrText),
			// ),
			col.New(1),
			col.New(3).Add(
				text.New(fmt.Sprintf("Fecha    : %s", d.CreatedAt.Format(time.DateOnly)), p.docAttrText),
			),
		),
		row.New().Add(
			col.New(4).Add(
				text.New(time.Now().Format(time.TimeOnly), p.docAttrText),
			),
			col.New(5),
		),
		row.New(2),
		row.New().Add(
			line.NewCol(12, props.Line{Style: linestyle.Dashed}),
		),
		row.New(2),
	)
}
