package company_admin_specifications

type CompanyAdminSpecifications interface {
	CreateCompany(name string) (string, error)
}



