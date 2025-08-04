package cash_outflow_pdf

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/repository"
	"erp/pkg/currency"
	"erp/pkg/exporter"
	"fmt"
	"strings"
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

type CashOutflowPdf interface {
	GenerateCashOutflowDocument(req *common.RequestContext, d dto.CashOutflowDto) (res []byte, err error)
}

type cashOutflowPdf struct {
	convertor helpers.ConvertorHelper
	currency  helpers.CurrencyHelper
	locale    helpers.Locale
	document  repository.DocumentService

	defaultFontSize   float64
	dfPad             float64
	docAttrText       props.Text
	defaultFontFamily string
	Q                 *query.Query
}

func NewCashOutflowPdf(
	helpers *helpers.Helpers,
	query *query.Query,
	document repository.DocumentService,
) CashOutflowPdf {
	defaultFontSize := 9.0
	dfPad := 1.0
	defaultFontFamily := fontfamily.Courier
	docAttrText := props.Text{
		Size: defaultFontSize,
		Top:  dfPad, Bottom: dfPad, Right: 1, Left: dfPad,
		Family: defaultFontFamily,
	}
	return &cashOutflowPdf{
		convertor:         helpers.Convertor,
		currency:          helpers.Currency,
		locale:            helpers.Locale,
		defaultFontSize:   defaultFontSize,
		dfPad:             dfPad,
		docAttrText:       docAttrText,
		defaultFontFamily: defaultFontFamily,
		Q:                 query,
		document:          document,
	}
}

func (p *cashOutflowPdf) GenerateCashOutflowDocument(req *common.RequestContext, d dto.CashOutflowDto) (
	res []byte, err error) {
	e := exporter.NewPdfExporter()
	b := exporter.NewIPdfBuilder()

	taxAndCharges, err := p.document.GetTaxLines(p.Q, req.Ctx, d.ID)
	if err != nil {
		return
	}

	e.SetBuilder(b)
	p.generateHeader(b, d, req.ActiveCompany)
	err = p.generateContent(b, d, req.ActiveCompany, req.CompanyDefaults, taxAndCharges)
	if err != nil {
		return
	}
	p.generateSignatureSection(b)

	return e.BuildPdfFile()
}

func (p *cashOutflowPdf) generateContent(b exporter.IPdfBuilder, d dto.CashOutflowDto, company model.Company,
	companyDefaults model.CompanyDefault, taxAndCharges dto.TaxLinesData) (err error) {
	amount := p.currency.IntToFloat(d.Amount)
	b.SetRows(
		row.New().Add(
			col.New(2).Add(
				text.New("Entregue a", p.docAttrText),
			),
			col.New(10).Add(
				text.New(fmt.Sprintf(": %s", d.Party), p.docAttrText),
			),
		),
		row.New().Add(
			col.New(2).Add(
				text.New("La suma de", p.docAttrText),
			),
			col.New(8).Add(
				text.New(fmt.Sprintf(": %s %s",
					strings.ToUpper(currency.AmountToWords(amount)),
					companyDefaults.Currency), p.docAttrText),
			),
		),
		row.New().Add(
			col.New(2).Add(
				text.New("Concepto", p.docAttrText),
			),
			col.New(8).Add(
				text.New(fmt.Sprintf(": %s", *d.Concept), p.docAttrText),
			),
		),

		row.New().Add(
			col.New(2).Add(
				text.New("T. Compra", p.docAttrText),
			),
			col.New(8).Add(
				text.New(fmt.Sprintf(": %s", *d.CashOutflowType), p.docAttrText),
			),
		),

		row.New().Add(
			col.New(2).Add(
				text.New("Factura", p.docAttrText),
			),
			col.New(3).Add(
				text.New(fmt.Sprintf(": %s", *d.InvoiceNo), p.docAttrText),
			),
			col.New(1),
			col.New(2).Add(
				text.New("Total Fact.", p.docAttrText),
			),
			col.New(3).Add(
				text.New(fmt.Sprintf(": %s", p.currency.FormatCurrencyAmount(
					taxAndCharges.TotalAmount+d.Amount,
					companyDefaults.Currency,
				)), p.docAttrText),
			),
		),

		row.New().Add(
			col.New(2).Add(
				text.New("No. Aut.", p.docAttrText),
			),
			col.New(3).Add(
				text.New(fmt.Sprintf(": %s", *d.AuthCode), p.docAttrText),
			),
			col.New(1),
			col.New(2).Add(
				text.New("Imp. Total.", p.docAttrText),
			),
			col.New(3).Add(
				text.New(fmt.Sprintf(": %.2f", p.currency.IntToFloat(taxAndCharges.TotalAmount)), p.docAttrText),
			),
		),

		row.New().Add(
			col.New(2).Add(
				text.New("Cdgo. Ctrl", p.docAttrText),
			),
			col.New(3).Add(
				text.New(fmt.Sprintf(": %s", *d.CtrlCode), p.docAttrText),
			),
			col.New(1),
			col.New(2).Add(
				text.New("Total Desc.", p.docAttrText),
			),
			col.New(3).Add(
				text.New(fmt.Sprintf(": %.2f", p.currency.IntToFloat(0)), p.docAttrText),
			),
		),
		row.New().Add(
			col.New(2).Add(
				text.New("Nit", p.docAttrText),
			),
			col.New(3).Add(
				text.New(fmt.Sprintf(": %s", *d.Nit), p.docAttrText),
			),
			col.New(1),
			col.New(2).Add(
				text.New("Impt. Neto", p.docAttrText),
			),
			col.New(3).Add(
				text.New(fmt.Sprintf(": %.2f", p.currency.IntToFloat(d.Amount)), p.docAttrText),
			),
		),
	)

	if d.EmisionDate != nil {
		b.SetRow(
			row.New().Add(
				col.New(2).Add(
					text.New("Emision", p.docAttrText),
				),
				col.New(3).Add(
					text.New(fmt.Sprintf(": %s", d.EmisionDate.Format(time.DateOnly)), p.docAttrText),
				),
			),
		)
	}
	b.SetRow(
		row.New().Add(
			col.New(2).Add(
				text.New("Proveedor", p.docAttrText),
			),
			col.New(3).Add(
				text.New(fmt.Sprintf(": %s", d.Party), p.docAttrText),
			),
		),
	)
	return
}

func (p *cashOutflowPdf) generateSignatureSection(b exporter.IPdfBuilder) {
	signatureProps := props.Signature{FontFamily: p.defaultFontFamily, LineStyle: linestyle.Dashed,
		FontSize: p.defaultFontSize}
	b.SetRow(
		row.New(40).Add(
			signature.NewCol(3, "PHA - Caja 200 \nU/N 2", signatureProps),
			col.New(1),
			signature.NewCol(3, "Vo. Bo.", signatureProps),
			col.New(1),
			signature.NewCol(3, "Interesado", signatureProps),
		),
	)
}

func (p *cashOutflowPdf) generateHeader(b exporter.IPdfBuilder, d dto.CashOutflowDto,
	company model.Company) {

	b.RegisterHeader(
		row.New().Add(
			col.New(4).Add(
				text.New("CAJA", p.docAttrText),
			),
			col.New(4).Add(
				text.New(company.Name, p.docAttrText),
			),
			col.New(1),
			col.New(3).Add(
				text.New(fmt.Sprintf("No Tran. : %s", d.Code), p.docAttrText),
			),
		),
		row.New().Add(
			col.New(4).Add(
				text.New(d.PostingDate.Format(time.DateOnly), p.docAttrText),
			),
			col.New(4).Add(
				text.New("EGRESO DE CAJA", p.docAttrText),
			),
			col.New(1),
			col.New(3).Add(
				text.New(fmt.Sprintf("Fecha    : %s", d.PostingDate.Format(time.DateOnly)), p.docAttrText),
			),
		),
		row.New().Add(
			col.New(4).Add(
				text.New(d.PostingTime, p.docAttrText),
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
