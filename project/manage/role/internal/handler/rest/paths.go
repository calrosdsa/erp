package role_rest


type RolePaths struct {
	Base string
	Detail string
	PermissionActions string
	RoleDefinitions string
	EntityActions string
}

func NewRolePaths(base string) RolePaths {
	return RolePaths{
		Base: base,
		PermissionActions: base + "/permision/actions",
		Detail: base + "/detail/{id}",
		RoleDefinitions: base + "/role-definitions",
		EntityActions: base + "/entity-actions",
	}
}
