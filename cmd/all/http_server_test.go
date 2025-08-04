package main_test

import (
	// "erp/cmd/all/internal"
	ptt_spec "erp/project/document/payment_terms_template/specifications"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	var (
		port   = "9090"
		client = &http.Client{
			Timeout: 1 * time.Second,
		}
		baseURL = fmt.Sprintf("http://localhost:%s", port)
	)

	// internal.StartDockerServer(t, port, "httpserver")

	ptt_spec.GreetSpecification(t, client, baseURL)
}
