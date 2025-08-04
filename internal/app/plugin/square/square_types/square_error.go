package squaretypes

type SquareErrors struct {
	Errors []SquareError `json:"errors"`
}

type SquareError struct {
	Code     string `json:"code"`
	Detail   string `json:"detail"`
	Category string `json:"category"`
}