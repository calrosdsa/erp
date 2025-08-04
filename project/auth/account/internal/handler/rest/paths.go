package rest_account

type AccountPaths struct {
	SignIn         string
	Account        string
	Password       string
	Sessions       string
	ResetPassoword string
	ChangePassword string
}

func NewAccountPath(base string) AccountPaths {
	return AccountPaths{
		SignIn:         base + "/sign-in",
		Account:        base,
		Password:       base + "/password",
		Sessions:       base + "/sessions",
		ResetPassoword: base + "/reset/password",
		ChangePassword: base + "/change/password",
	}
}

type RolePaths struct {
	Base              string
	Detail            string
	PermissionActions string
	RoleDefinitions   string
	EntityActions     string
}

func NewRolePaths(base string) RolePaths {
	return RolePaths{
		Base:              base,
		PermissionActions: base + "/permision/actions",
		Detail:            base + "/detail/{id}",
		RoleDefinitions:   base + "/role-definitions",
		EntityActions:     base + "/entity-actions",
	}
}
