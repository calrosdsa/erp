package cuatropf

type CuatropfPath struct {
	Base         string
	Subscription string
}

func NewCuatropfPath(base string) CuatropfPath {
	return CuatropfPath{
		Base:         base,
		Subscription: base + "/subscription/{companyUuid}",
	}
}
