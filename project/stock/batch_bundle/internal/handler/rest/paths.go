package rest_batch_bundle

type BatchBundlePaths struct {
	Base   string
	Detail string
}

func NewBatchBundlePaths(base string) BatchBundlePaths {
	return BatchBundlePaths{
		Base:   base,
		Detail: base + "/detail/{id}",
	}
}
