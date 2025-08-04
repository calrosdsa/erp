package domain

type permissionOpts struct {
	EntityType      PermifyEntity
	EntityID        string
	Permission      Permission
	SubjectType     PermifyEntity
	SubjectID       string
	SubjectRelation string
}

var PermissionOpts permissionOpts

type PermissionOpt func(opts *permissionOpts)

func (*permissionOpts) WithEntityType(entityType PermifyEntity) PermissionOpt {
	return func(opts *permissionOpts) {
		opts.EntityType = entityType
	}
}
func (*permissionOpts) WithEntityID(entityId string) PermissionOpt {
	return func(opts *permissionOpts) {
		opts.EntityID = entityId
	}
}
func (*permissionOpts) WithPermission(permission Permission) PermissionOpt {
	return func(opts *permissionOpts) {
		opts.Permission = permission
	}
}
func (*permissionOpts) WithSubjectType(subjectType PermifyEntity) PermissionOpt {
	return func(opts *permissionOpts) {
		opts.SubjectType = subjectType
	}
}
func (*permissionOpts) WithSubjectId(subjectId string) PermissionOpt {
	return func(opts *permissionOpts) {
		opts.SubjectID = subjectId
	}
}
func (*permissionOpts) WithSubjectRelation(subjectRelation string) PermissionOpt {
	return func(opts *permissionOpts) {
		opts.SubjectRelation = subjectRelation
	}
}

func (*permissionOpts) Apply(opts ...PermissionOpt) permissionOpts {
	options := permissionOpts{}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}
