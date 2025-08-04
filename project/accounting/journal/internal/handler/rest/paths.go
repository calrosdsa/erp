package rest_journal

type JournalPaths struct {
	Base   string
	Detail string
	UpdateStatus string 
}

func NewJournalPaths(base string) JournalPaths {
	return JournalPaths{
		Base:   base,
		Detail: base + "/detail/{id}",
		UpdateStatus: base + "/update-status",
	}
}
