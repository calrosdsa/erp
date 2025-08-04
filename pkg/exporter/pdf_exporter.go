package exporter

type PdfExporter interface {
	SetBuilder(b IPdfBuilder)
	Save(fileName string)
	BuildPdfFile()([]byte, error)
}

type pdfExporter struct {
	b IPdfBuilder
}

func NewPdfExporter() PdfExporter {
	return &pdfExporter{}
}

// func (e *pdfExporter) Export(m core.Maroto, headers []interface{}, data [][]interface{}) (res []byte, err error) {
// 	return
// }
func(e *pdfExporter)Save(fileName string) {
	e.b.Save(fileName)
}

func (e *pdfExporter) BuildPdfFile() (res []byte, err error) {
	return e.b.Build()
}

func (e *pdfExporter) SetBuilder(b IPdfBuilder) {
	e.b = b
}
