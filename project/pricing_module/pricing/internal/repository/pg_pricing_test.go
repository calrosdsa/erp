package pricing_repo

import (
	"erp/gen/db/model"
	"erp/internal/domain"
	"fmt"
	"testing"
)

func TestAssignement(t *testing.T){
	m := model.PricingLineItem{}
	// description := "DESCRIPTION"
	// domain.AssignIfNotEmpty(&m.Description,&description)
	var strField *string
	strVal := "Hello"
	fmt.Println("FROM FMT",strField,m)
	domain.AssignIfNotEmpty(&strField, strVal) // Should work
	fmt.Println("FROM FMT",strField)
	// domain.AssignIfNotEmpty(m.Description,"DESCRIPCION")
	t.Log("LOG..")
}