package a_company_rest

type CompanyAPaths struct {
	Base    string
	Detail  string
	Modules string
	User    string
}

func NewACompanyAdminPaths(base string) CompanyAPaths {
	return CompanyAPaths{
		Base:    base,
		Detail:  base + "/detail/{id}",
		Modules: base + "/modules",
		User:    base + "/user",
	}
}
