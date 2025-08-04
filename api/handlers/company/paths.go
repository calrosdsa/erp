package company

type CompanyPaths struct {
	Base                 string
	ValidParentCompanies string
	Detail               string
	Uuid                 string
	User                 string
}

func NewCompanyPath(base string) CompanyPaths {
	return CompanyPaths{
		Base:                 base,
		ValidParentCompanies: base + "/valid/parent/companies",
		Detail:               base + "/detail/{id}",
		Uuid:                 base + "/{uuid}",
		User:                 base + "/user/companies",
	}
}
