package resolvers

import (
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nyaruka/phonenumbers"
)

type PhoneNumber struct {
	Number      string `json:"number"`
	CountryCode string `json:"countryCode"`
}

func (i PhoneNumber) Resolve(ctx huma.Context, prefix *huma.PathBuffer) []error {
	num, err := phonenumbers.Parse(i.Number, i.CountryCode)
	validation := phonenumbers.IsPossibleNumber(num)
	fmt.Println("VALIDATION",validation)
	if err != nil {	
		fmt.Println(err)
		return []error{&huma.ErrorDetail{
			Location: prefix.String(),
			Message:  "Invalid Phone Number",
			Value:    i,
		}}
	}
	
	fmt.Println(*num.CountryCode,*num.NationalNumber,num.ProtoReflect().IsValid())
	return nil
}