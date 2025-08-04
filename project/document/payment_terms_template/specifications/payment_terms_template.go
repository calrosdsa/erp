package ptt_spec

import (
	payment_terms_t_rest "erp/project/document/payment_terms_template/handler/rest"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type GreeterSpec interface {
	Greet(name string) (string, error)
}

func GreetSpecification(t *testing.T,client  *http.Client,baseRoute string) {
	adapter := payment_terms_t_rest.NewDriver(client,baseRoute)
	greetSpecification(t,&adapter)
	// t.Run("Introduce", func(t *testing.T) {
	// 	got, err := greeter.Introduce("Mike")
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, got, "My name is, Mike")
	// })
}

func greetSpecification(t *testing.T,greeter GreeterSpec) {
	t.Run("Greet", func(t *testing.T) {
		got, err := greeter.Greet("Mike")
		assert.NoError(t, err)
		assert.Equal(t, got, "Hello, Mike")
	})
}
