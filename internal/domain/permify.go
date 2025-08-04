package domain

type PermifyEntity string

const (
	COMPANY_ENTITY  PermifyEntity = "company"
	ROLE_ENTITY     PermifyEntity = "role"
	USER_ENTITY     PermifyEntity = "user"
	TEMPLATE_ENTITY PermifyEntity = "template"
)

type Permission string

const (
	//Company
	CREATE_USER Permission = "create_user"
)
