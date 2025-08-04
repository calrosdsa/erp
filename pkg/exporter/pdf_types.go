package exporter

import "github.com/johnfercher/maroto/v2/pkg/props"

type ColData struct {
	Size      int
	Value     string
	TextProps *props.Text
}


type Colors struct {
	Gray *props.Color
	DarkGray *props.Color
}