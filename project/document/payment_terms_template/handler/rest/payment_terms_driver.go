package payment_terms_t_rest

import (
	"erp/internal/domain"
	"io"
	"net/http"
)

type Driver struct {
	client  *http.Client
	paths Paths
}

func NewDriver(
	client *http.Client,
	base string,
) Driver {
	paths := NewPaths(base + domain.PAYMENT_TERMS_TEMPLATE_ROUTE)
	return Driver{
		client: client,
		paths: paths,
	}
}

func (d *Driver) Greet(name string) (string, error) {
	return d.getAndReadFrom(d.paths.Greet, name)
}


func (d *Driver) getAndReadFrom(path string, name string) (string, error) {
	res, err := d.client.Get(path + "?name=" + name)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	greeting, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(greeting), nil
}
