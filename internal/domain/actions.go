package domain
type ActionType string 
const (
	VIEW ActionType ="view"
	CREATE ActionType ="create"
	DELETE ActionType ="delete"
	EDIT ActionType ="edit"
)