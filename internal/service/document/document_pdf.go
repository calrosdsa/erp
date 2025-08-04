package document_service


import (
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/domain/repository"
	"erp/pkg/exporter"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/border"
	"github.com/johnfercher/maroto/v2/pkg/consts/extension"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)



func (s *documentService) generatePurchaseOrder(b exporter.IPdfBuilder) {
	b.SetRow(row.New(3))
	headerText := props.Text{Size: 7, Align: align.Center, Style: fontstyle.Bold,
		Top: 1, Bottom: 1, Right: 1, Left: 1, Color: &props.WhiteColor}

	b.SetRow(
		row.New(5).Add(
			text.NewCol(2, "FECHA DE ENTREGA", headerText),
			text.NewCol(2, "SOLICITANTE", headerText),
			text.NewCol(3, "ENVIADO MEDIANTE", headerText),
			text.NewCol(3, "TERMINOS Y CONDICIONES", headerText),
			text.NewCol(2, "PUNTO F.O.B", headerText),
		).WithStyle(&props.Cell{BackgroundColor: b.Colors().DarkGray}),
	)

	colStyle := &props.Cell{
		BorderType:  border.Full,
		BorderColor: b.Colors().Gray,
	}
	b.SetRow(row.New().Add(
		// col.New(3),
		text.NewCol(2, "2024-05-12", props.Text{Size: 8, Align: align.Center,
			Top: 1, Bottom: 1, Right: 1, Left: 1}).WithStyle(colStyle),
		text.NewCol(2, "", props.Text{Size: 7, Align: align.Center,
			Top: 1, Bottom: 1, Right: 1, Left: 1}).WithStyle(colStyle),
		text.NewCol(3, "", props.Text{Size: 8, Align: align.Center,
			Top: 1, Bottom: 1, Right: 1, Left: 1}).WithStyle(colStyle),
		text.NewCol(3, "", props.Text{Size: 8, Align: align.Center,
			Top: 1, Bottom: 1, Right: 1, Left: 1}).WithStyle(colStyle),
		text.NewCol(2, "", props.Text{Size: 8, Align: align.Center,
			Top: 1, Bottom: 1, Right: 1, Left: 1}).WithStyle(colStyle),
	))

	b.SetRow(row.New(3))
}

type AddressSection struct {
	Address dto.AddressDto
	Section string
}

func (s *documentService) generateAddressSection(b exporter.IPdfBuilder, addr1 AddressSection, addr2 *AddressSection) {
	docAttrText := props.Text{
		Size: s.dFontSize,
		Top:  s.dPad, Bottom: s.dPad, Right: 1, Left: s.dPad,
	}

	row1 := row.New(6).Add(
		col.New(5).Add(
			text.New(addr1.Section, props.Text{
				Size: s.dFontSize,
				Top:  2, Bottom: 2, Right: s.dPad, Left: s.dPad,
				Style: fontstyle.Bold,
			}),
		).WithStyle(&props.Cell{
			BorderType: border.Bottom,
		}),
	)
	if addr2 != nil {
		row1.Add(col.New(1),
			col.New(5).Add(
				text.New(addr2.Section, props.Text{
					Size: s.dFontSize,
					Top:  2, Bottom: 2, Right: s.dPad, Left: s.dPad,
					Style: fontstyle.Bold,
				}),
			).WithStyle(&props.Cell{
				BorderType: border.Bottom,
			}),
		)
	}
	row2 := row.New().Add(
		col.New(4).Add(
			text.New(nullString(addr1.Address.Company), docAttrText),
		),
		col.New(2),
	)
	if addr2 != nil {
		row2.Add(
			col.New(4).Add(text.New(nullString(addr2.Address.Company), docAttrText)),
		)
	}

	row3 := row.New().Add(
		col.New(4).Add(
			text.New(addr1.Address.StreetLine1+addr1.Address.StreetLine2, docAttrText),
		),
		col.New(2),
	)
	if addr2 != nil {
		row3.Add(
			col.New(4).Add(text.New(addr2.Address.StreetLine1+addr2.Address.StreetLine2, docAttrText)))
	}
	row4 := row.New().Add(
		col.New(4).Add(
			text.New(fmt.Sprintf("%s, %s, %s", addr1.Address.City, nullString(addr1.Address.Province),
				nullString(addr1.Address.PostalCode)), docAttrText),
		),
		col.New(2),
	)
	if addr2 != nil {
		row4.Add(
			col.New(4).Add(
				text.New(fmt.Sprintf("%s, %s, %s", addr2.Address.City, nullString(addr2.Address.Province),
					nullString(addr2.Address.PostalCode)), docAttrText),
			),
		)
	}
	row5 := row.New().Add(
		col.New(4).Add(
			text.New(fmt.Sprintf("%s, %s, %s", addr1.Address.City, nullString(addr1.Address.Province),
				nullString(addr1.Address.PostalCode)), docAttrText),
		),
		col.New(2),
	)
	if addr2 != nil {
		row5.Add(col.New(4).Add(
			text.New(fmt.Sprintf("%s, %s, %s", addr2.Address.City, nullString(addr2.Address.Province),
				nullString(addr2.Address.PostalCode)), docAttrText),
		))
	}
	row6 := row.New().Add(
		col.New(4).Add(
			text.New(nullString(addr1.Address.PhoneNumber), docAttrText),
		),
		col.New(2),
	)
	if addr2 != nil {
		row6.Add(col.New(4).Add(
			text.New(nullString(addr2.Address.PhoneNumber), docAttrText),
		))
	}

	row7 := row.New()
	if addr1.Address.Email != nil {
		row7.Add(
			col.New(4).Add(
				text.New(nullString(addr1.Address.Email), docAttrText),
			),
			col.New(2),
		)
	}

	if addr2.Address.Email != nil {
		row7.Add(
			col.New(4).Add(
				text.New(nullString(addr2.Address.PhoneNumber), docAttrText),
			),
		)
	}

	b.SetRows(row1, row2, row3, row4, row5, row6, row7)

}


type DocInfo struct {
	Date    time.Time
	DocNo   string
	Address dto.AddressDto
	LogoUrl *string
	Title   string
}


func (s *documentService) generateItemsAndChargesSection(b exporter.IPdfBuilder, items dto.LineItemsData,
	tacs dto.TaxLinesData, currency string) {
	b.SetRow(row.New(2))

	colStyle := &props.Cell{
		BorderType:  border.Full,
		BorderColor: b.Colors().Gray,
	}
	headerText := props.Text{Size: 7, Align: align.Center, Style: fontstyle.Bold,
		Top: 1, Bottom: 1, Right: 1, Left: 1, Color: &props.WhiteColor}

	b.SetRow(
		row.New().Add(
			text.NewCol(3, "ARTICULO", headerText),
			text.NewCol(4, "DESCRIPCION", headerText),
			text.NewCol(1, "CANTIDAD", headerText),
			text.NewCol(2, "PRECIO UNITARIO", headerText),
			text.NewCol(2, "TOTAL", headerText),
		).WithStyle(&props.Cell{BackgroundColor: b.Colors().DarkGray}),
	)

	tableCellText := props.Text{Size: 8, Align: align.Center,
		Top: 1, Bottom: 1, Right: 1, Left: 1}
	var contentsRow []core.Row
	for _, content := range items.LineItems {
		r := row.New().Add(
			// col.New(3),
			text.NewCol(2, content.ItemName, tableCellText).WithStyle(colStyle),
			text.NewCol(5, content.ItemDescription, props.Text{Size: 7, Align: align.Left,
				Top: 1, Bottom: 1, Right: 1, Left: 1}).WithStyle(colStyle),
			text.NewCol(1, fmt.Sprintf("%d", content.Quantity), props.Text{Size: 8, Align: align.Left,
				Top: 1, Bottom: 1, Right: 1, Left: 1}).WithStyle(colStyle),
			text.NewCol(2, s.currency.FormatCurrencyAmount(content.Rate, currency), props.Text{Size: 8, Align: align.Left,
				Top: 1, Bottom: 1, Right: 1, Left: 1}).WithStyle(colStyle),
			text.NewCol(2, s.currency.FormatCurrencyAmount(content.Rate*int64(content.Quantity), currency), props.Text{Size: 8, Align: align.Left,
				Top: 1, Bottom: 1, Right: 1, Left: 1}).WithStyle(colStyle),
		)
		contentsRow = append(contentsRow, r)
	}
	b.SetRows(contentsRow...)

	totalText := props.Text{
		Size:  8,
		Align: align.Right,
		Top:   4, Bottom: 1, Right: 1, Left: 1,
	}

	b.SetRow(
		row.New().Add(
			text.NewCol(2, "Notas y observaciones"),
			col.New(5),
			text.NewCol(3, "SUBTOTAL", totalText),
			text.NewCol(2, s.currency.FormatCurrencyAmount(items.TotalAmount, currency), totalText).WithStyle(
				&props.Cell{
					BorderType: border.Bottom,
				},
			),
		),
	)
	b.SetRow(
		row.New().Add(
			col.New(7),
			text.NewCol(3, "DESCUENTO", totalText),
			text.NewCol(2, "", totalText).WithStyle(
				&props.Cell{
					BorderType: border.Bottom,
				},
			),
		),
	)

	b.SetRow(
		row.New().Add(
			col.New(7),
			text.NewCol(3, "SUBTOTAL MENOS DESCUENTO", totalText),
			text.NewCol(2, s.currency.FormatCurrencyAmount(items.TotalAmount, currency), totalText).WithStyle(
				&props.Cell{
					BorderType: border.Bottom,
				},
			),
		),
	)
	b.SetRow(
		row.New().Add(
			col.New(7),
			text.NewCol(3, "TOTAL IMPUESTOS Y CARGOS", totalText),
			text.NewCol(2, s.currency.FormatCurrencyAmount(tacs.TotalAmount, currency), totalText).WithStyle(
				&props.Cell{
					BorderType: border.Bottom,
				},
			),
		),
	)

	b.SetRow(
		row.New().Add(
			col.New(7),
			text.NewCol(3, "TOTAL", totalText),
			text.NewCol(2, s.currency.FormatCurrencyAmount(items.TotalAmount+tacs.TotalAmount, currency),
				totalText).WithStyle(
				&props.Cell{
					BorderType: border.Bottom,
				},
			),
		),
	)
}


func (s *documentService) getItemsAndCharges(req *common.RequestContext, docID int64) (items dto.LineItemsData,
	tacs dto.TaxLinesData, err error) {
	items, err = s.GetLineItems(s.Q, req, docID,
		repository.OptionsStock.WithLoadItemInLine(true))
	if err != nil {
		return
	}
	tacs, err = s.GetTaxLines(s.Q, req.Ctx, docID)
	if err != nil {
		return
	}
	return
}

func nullString(s *string) string {
	if s == nil {
		return ""
	} else {
		return *s
	}
}

func getImageBytes(url string) ([]byte, error) {
	// Send GET request to the image URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the image data into a byte slice
	imageBytes, err := io.ReadAll(resp.Body)
	fmt.Println("ERR",err)
	if err != nil {
		return nil, err
	}

	return imageBytes, nil
}


func (s *documentService) generateDocInfo(b exporter.IPdfBuilder, d DocInfo) {
	docAttrText := props.Text{
		Size: s.dFontSize,
		Top:  s.dPad, Bottom: s.dPad, Right: 1, Left: s.dPad,
	}

	headerInfo := row.New(10)
	if d.LogoUrl != nil {
		bytes, err := getImageBytes(*d.LogoUrl)
		if err != nil {
			fmt.Println(err)
		}
		headerInfo.Add(
			col.New(2).Add(image.NewFromBytes(bytes, extension.Png)),
		)
	} else {
		headerInfo.Add(
			col.New(2),
		)
	}
	headerInfo.Add(
		col.New(6),
		col.New(4).Add(
			text.New(d.Title, props.Text{
				Size:  16,
				Align: align.Right,
			}),
		))
	b.SetRow(headerInfo)

	b.SetRow(row.New().Add(
		// col.New(2).Add(
		// 	text.New("Empresa", docAttrText),
		// ),
		col.New(4).Add(
			text.New(nullString(d.Address.Company), props.Text{
				Size:  9.0,
				Style: fontstyle.Bold,
				Top:   s.dPad, Bottom: s.dPad, Right: 1, Left: s.dPad,
			}),
		),
		col.New(3),
		col.New(2).Add(
			text.New("Fecha", docAttrText),
		).WithStyle(
			&props.Cell{
				BorderType: border.Full,
			},
		),
		col.New(3).Add(
			text.New(d.Date.Format(time.DateOnly), docAttrText),
		).WithStyle(
			&props.Cell{
				BorderType: border.Full,
			},
		),
	))

	b.SetRow(row.New().Add(

		col.New(4).Add(
			text.New(d.Address.StreetLine1+d.Address.StreetLine2, docAttrText),
		),
		col.New(3),
		col.New(2).Add(
			text.New("No de Orden", docAttrText),
		).WithStyle(
			&props.Cell{
				BorderType: border.Full,
			},
		),
		col.New(3).Add(
			text.New(d.DocNo, docAttrText),
		).WithStyle(
			&props.Cell{
				BorderType: border.Full,
			},
		),
	))
	b.SetRow(row.New().Add(
		col.New(3).Add(
			text.New(d.Address.City, docAttrText),
		),
	))
}
