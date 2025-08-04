package rest_project

type ProjectPaths struct {
	Base   string
	Detail string
	TestRequest string
}

func NewProjectPaths(base string) ProjectPaths {
	return ProjectPaths{
		Base:   base,
		Detail: base + "/detail/{id}",
		TestRequest: "/delete-collectors",
	}
}
