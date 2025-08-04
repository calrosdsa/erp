package payment_pdf

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/app/service/helpers"
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
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/border"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontfamily"
	"github.com/johnfercher/maroto/v2/pkg/consts/linestyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)


type PaymentPDF  interface {
	GeneratePaymentDocument(req *common.RequestContext, d dto.PaymentDetailDto) (res []byte, err error)
}

type paymentPDF struct {
	convertor helpers.ConvertorHelper
	currency  helpers.CurrencyHelper
	locale    helpers.Locale

	defaultFontSize   float64
	dfPad             float64
	docAttrText       props.Text
	defaultFontFamily string
	Q                 *query.Query
}

func NewPaymentPdf(
	helpers *helpers.Helpers,
	query *query.Query,
) PaymentPDF {
	defaultFontSize := 9.0
	dfPad := 1.0
	defaultFontFamily := fontfamily.Courier
	docAttrText := props.Text{
		Size: defaultFontSize,
		Top:  dfPad, Bottom: dfPad, Right: 1, Left: dfPad,
		Family: defaultFontFamily,
	}
	return &paymentPDF{
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

func (p *paymentPDF) GeneratePaymentDocument(req *common.RequestContext, d dto.PaymentDetailDto) (res []byte, err error) {
	e := exporter.NewPdfExporter()
	b := exporter.NewIPdfBuilder()
	e.SetBuilder(b)
	p.generateHeader(b, d, req.ActiveCompany)
	err = p.generateContent(b, d, req.ActiveCompany)
	if err != nil {
		return
	}
	p.generateRfcsSection(b, d.PaymentReferences)
	p.generateSignatureSection(b)

	return e.BuildPdfFile()
}

func (p *paymentPDF) generateContent(b exporter.IPdfBuilder, d dto.PaymentDetailDto, company model.Company) (err error) {
	t := p.locale.Translate("es")	
	b.SetRow(
		row.New().Add(
			col.New(2).Add(
				text.New(t(fmt.Sprintf("Party.%s",d.PartyType)), p.docAttrText),
			),
			col.New(10).Add(
				text.New(fmt.Sprintf(": %s", d.PartyName), p.docAttrText),
			),
		))
	if d.CompanyBankAccountID != nil {
		cba, err := p.getBankAccount(*d.CompanyBankAccountID, company.ID)
		if err != nil {
			return err
		}
		b.SetRow(
			row.New().Add(
				col.New(2).Add(
					text.New("No. Cta.", p.docAttrText),
				),
				col.New(8).Add(
					text.New(fmt.Sprintf(": %s - %s", *cba.BankAccountNumber, cba.Bank), p.docAttrText),
				),
				col.New(2).Add(
					text.New(p.currency.FormatCurrencyAmount(int64(d.Amount), d.PaidFromCurrency), p.docAttrText),
				),
			))
	}
	// b.SetRow(
	// 	row.New().Add(
	// 		col.New(2).Add(
	// 			text.New("No. Tran.", p.docAttrText),
	// 		),
	// 		col.New(8).Add(
	// 			text.New(": 397", p.docAttrText),
	// 		),
	// 	))
	// b.SetRow(
	// 	row.New().Add(
	// 		col.New(2).Add(
	// 			text.New("U/N", p.docAttrText),
	// 		),
	// 		col.New(8).Add(
	// 			text.New(": 2 AGEN. OFICINA CENTRAL", p.docAttrText),
	// 		),
	// ))
	amount := p.currency.IntToFloat(d.Amount)

	b.SetRow(
		row.New().Add(
			col.New(2).Add(
				text.New("Son", p.docAttrText),
			),
			col.New(8).Add(
				text.New(fmt.Sprintf(": %s %s",
					strings.ToUpper(currency.AmountToWords(amount)),
					d.PaidFromCurrency), p.docAttrText),
			),
		))

	if d.ChequeReferenceNo != nil {
		b.SetRow(
			row.New().Add(
				col.New(2).Add(
					text.New("C/NO. DE REF.", p.docAttrText),
				),
				col.New(8).Add(
					text.New(fmt.Sprintf(": %s", *d.ChequeReferenceNo), p.docAttrText),
				),
			))
	}
	if d.ChequeReferenceDate != nil {
		b.SetRow(
			row.New().Add(
				col.New(2).Add(
					text.New("C/FCH. DE REF.", p.docAttrText),
				),
				col.New(8).Add(
					text.New(fmt.Sprintf(": %s", d.ChequeReferenceDate.Format(time.DateOnly)), p.docAttrText),
				),
			))

	}
	return
}

func (p *paymentPDF) generateRfcsSection(b exporter.IPdfBuilder, references []dto.PaymentReferenceDto) (err error) {
	if len(references) == 0 {
		return
	}
	t := p.locale.Translate("es")
	headerText := props.Text{Size: 8, Align: align.Center,
		Top: 1, Bottom: 1, Right: 1, Left: 1,
		Family: p.defaultFontFamily,
	}

	colStyle := &props.Cell{
		BorderType: border.Top,
		LineStyle:  linestyle.Dashed,
	}

	b.SetRow(row.New(3))
	b.SetRow(
		row.New(5).Add(
			text.NewCol(3, "Tipo", headerText).WithStyle(colStyle),
			text.NewCol(3, "No. Doc.", headerText).WithStyle(colStyle),
			text.NewCol(3, "Detalle", headerText).WithStyle(colStyle),
			text.NewCol(3, "Importe", headerText).WithStyle(colStyle),
		),
	)
	var contentsRow []core.Row

	tableCellText := props.Text{Size: 8, Align: align.Center,
		Family: p.defaultFontFamily,
		Top:    1, Bottom: 1, Right: 1, Left: 1}

	for _, reference := range references {
		r := row.New().Add(
			// col.New(3),
			text.NewCol(3, t(fmt.Sprintf("Party.%s",reference.PartyType)), tableCellText).WithStyle(colStyle),
			text.NewCol(3, reference.PartyCode, tableCellText).WithStyle(colStyle),
			text.NewCol(3, "", tableCellText).WithStyle(colStyle),
			text.NewCol(3, p.currency.FormatCurrencyAmount(reference.Total, reference.Currency), tableCellText).WithStyle(colStyle),
		)
		// if i%2 == 0 {
		// r.WithStyle(&props.Cell{BorderColor: gray,BorderType: border.Full})
		// }

		contentsRow = append(contentsRow, r)
	}
	b.SetRows(contentsRow...)
	return
}

func (p *paymentPDF) getBankAccount(bankAccountID int64, companyID int64) (res dto.BankAccountDto, err error) {
	bankQ := p.Q.Bank
	bankAQ := p.Q.BankAccount
	err = p.Q.BankAccount.Select(
		bankAQ.BankAccountNumber,
		bankQ.Name.As("bank"),
	).Join(
		bankQ, bankQ.ID.EqCol(bankAQ.BankID),
	).Where(
		bankAQ.ID.Eq(bankAccountID),
		bankAQ.CompanyID.Eq(companyID),
	).Scan(&res)
	return
}

func (p *paymentPDF) generateSignatureSection(b exporter.IPdfBuilder) {
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
	b.SetRow(row.New(2))
	b.SetRow(
		row.New().Add(
			col.New(7),
			col.New(2).Add(
				text.New("Nombre: ", props.Text{Size: p.defaultFontSize,
					Top: p.dfPad, Bottom: p.dfPad, Right: 1, Left: p.dfPad,
					Family: p.defaultFontFamily,
					Align:  align.Right,
				}),
			),
			col.New(3).Add(
				text.New(""),
			).WithStyle(&props.Cell{BorderType: border.Bottom}),
		),
	)
	b.SetRow(
		row.New().Add(
			col.New(7),
			col.New(2).Add(
				text.New("No.Doc: ", props.Text{Size: p.defaultFontSize,
					Top: p.dfPad, Bottom: p.dfPad, Right: 1, Left: p.dfPad,
					Family: p.defaultFontFamily,
					Align:  align.Right,
				}),
			),
			col.New(3).Add(
				text.New(""),
			).WithStyle(&props.Cell{BorderType: border.Bottom}),
		),
	)
}

func (p *paymentPDF) generateHeader(b exporter.IPdfBuilder, d dto.PaymentDetailDto,
	company model.Company) {

	b.RegisterHeader(
		row.New().Add(
			col.New(4).Add(
				text.New("ENTRADA DE PAGO", p.docAttrText),
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
				text.New(time.Now().Format(time.DateOnly), p.docAttrText),
			),
			col.New(4),
			// .Add(
			// 	text.New("REGISTRO CUENTA POR PAGAR", p.docAttrText),
			// ),
			col.New(1),
			col.New(3).Add(
				text.New(fmt.Sprintf("Fecha    : %s", d.PostingDate.Format(time.DateOnly)), p.docAttrText),
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
