package exporter

import (
	"fmt"

	"github.com/johnfercher/maroto/v2"
	// "github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontfamily"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/core/entity"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type pdfDo struct {
	m      core.Maroto
	colors Colors
}

type IPdfBuilder interface {
	SetRow(row core.Row) IPdfBuilder
	Build() ([]byte, error)
	Save(fileName string)
	SetRows(rows ...core.Row) IPdfBuilder
	Colors() Colors
	RegisterHeader(row ...core.Row)
	RegisterFooter(row ...core.Row)
	AddPages(pages ...core.Page)
	// ExportInvoiceDoc(m core.Maroto, headers []interface{}, data [][]interface{}) ([]byte, error)
}

func NewIPdfBuilder(args ...interface{}) IPdfBuilder {
	var (
		cfg *entity.Config
	)
	pageNumber := props.PageNumber{
		Pattern: "Página {current} de {total}",
		Place:   props.Bottom,
		Family:  fontfamily.Courier,
		Style:   fontstyle.Bold,
		Size:    9,
	}
	if len(args) == 1 {
		if val, ok := args[0].(*entity.Config); ok {
			cfg = val
		}
	} else {
		cfg = config.NewBuilder().
			WithLeftMargin(10).
			WithTopMargin(20).
			WithPageNumber(pageNumber).
			WithRightMargin(10).
			WithBottomMargin(20).
			Build()
	}
	// darkGrayColor := getDarkGrayColor()
	mrt := maroto.New(cfg)
	m := maroto.NewMetricsDecorator(mrt)
	b := &pdfDo{
		m:      m,
		colors: NewColors(),
	}
	return b
}

func (p *pdfDo) AddPages(pages ...core.Page) {
	p.m.AddPages(pages...)
}

func (p *pdfDo) RegisterHeader(r ...core.Row) {

	p.m.RegisterHeader(r...)
}

func (p *pdfDo) RegisterFooter(r ...core.Row) {
	p.m.RegisterFooter(r...)
}

func (p *pdfDo) Colors() Colors {
	return p.colors
}

func (p *pdfDo) Build() (res []byte, err error) {
	document, err := p.m.Generate()
	if err != nil {
		fmt.Println("Fail to generate pdf")
		return
	}
	return document.GetBytes(), err
}

func (p *pdfDo) Save(fileName string) {
	document, err := p.m.Generate()
	if err != nil {
		fmt.Println("Fail to generate pdf")
	}
	err = document.Save(fileName)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (p *pdfDo) SetHeader() {
}

func (p *pdfDo) SetRow(row core.Row) IPdfBuilder {
	rows := []core.Row{
		row,
	}
	p.m.AddRows(rows...)
	return p
}

// func (p *pdfDo) SetColSpan(row core.Row) IPdfBuilder {

// 	p.m.AddRows(rows...)
// 	return p
// }

func (p *pdfDo) SetRows(rows ...core.Row) IPdfBuilder {
	// for i, row := range rows {
	// if i %2 == 0 {
	// gray := getGrayColor()
	// row.WithStyle(&props.Cell{BackgroundColor: gray})
	// }
	// }
	p.m.AddRows(rows...)
	return p
}

func NewColors() Colors {
	return Colors{
		Gray: &props.Color{
			Red:   200,
			Green: 200,
			Blue:  200,
		},
		DarkGray: &props.Color{
			Red:   55,
			Green: 55,
			Blue:  55,
		},
	}
}
