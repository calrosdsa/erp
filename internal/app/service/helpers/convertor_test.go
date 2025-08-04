package helpers_test

import (
	"erp/internal/app/service/helpers"
	"testing"

	"github.com/stretchr/testify/assert"
)


type ExampleStruct struct {
	Name string
	ID int
}

func TestUpdateStructData(t *testing.T) {
	convertor := helpers.NewConvertorHelper()
	exampleStruct := &ExampleStruct{
		ID: 1,
		Name: "Name",
	}
	destStruct := ExampleStruct{ID: 2,}

	t.Log("EXAMPLE STRUCT BEFORE",exampleStruct)
	err :=convertor.UpdateStruct(exampleStruct,destStruct)
	t.Log("EXAMPLE STRUCT After",exampleStruct)
	assert.NoError(t,err)
	assert.Equal(t,exampleStruct.ID,2)
	assert.Equal(t,exampleStruct.Name,"Name")

}