package chart_rest

type ChartPaths struct {
	Base  string
	Chart string
}

func NewChartPaths(base string) ChartPaths {
	return ChartPaths{
		Base:  base,
		Chart: base + "/{chart}",
	}
}
