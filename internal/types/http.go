package types

type ErrorResponse struct {
	Error interface{} `json:"error"`
}

type SucessResponse struct {
	Sucess bool        `json:"sucess"`
	Data   interface{} `json:"data"`
}
