package court_rest

type CourtPaths struct {
	Base string 
	Detail string
}

func NewCourtPaths(base string)CourtPaths{
	return CourtPaths{
		Base: base,
		Detail: base + "/detail/{id}",
	}
}

type CourtRatePaths struct {
	Base string 
	Court string 
}

func NewCourtRatePaths(base string) CourtRatePaths {
	return CourtRatePaths{
		Base: base,
		Court: base + "/{id}",
	}
}