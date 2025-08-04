package auth_admin_rest

type AuthAdminPaths struct {
	SignIn string
}

func NewAuthAdminPaths(base string) AuthAdminPaths {
	return AuthAdminPaths{
		SignIn: base + "/sign-in",
	}
}

type RoleTemplatePaths struct {
	Base   string
	Detail string
}

func NewRoleTemplatePaths(base string) RoleTemplatePaths {
	return RoleTemplatePaths{
		Base:   base,
		Detail: base + "/detail/{id}",
	}
}
