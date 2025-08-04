package rest_serial_no

type SerialNoPaths struct {
	Base   string
	Detail string
	SerialNoTransactions string
}

func NewSerialNoPaths(base string) SerialNoPaths {
	return SerialNoPaths{
		Base:   base,
		Detail: base + "/detail/{id}",
		SerialNoTransactions: base + "/transactions",
	}
}
