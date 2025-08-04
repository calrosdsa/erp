package helpers 

type PartyHelper interface {

}

type partyHelper struct {
}

func NewHelperParty() PartyHelper{
	return  &partyHelper{}
}
