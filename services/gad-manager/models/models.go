package models

type CalcRequest struct {
	Exprs []string `json:"exprs"`
}

type CaclResponse struct {
	Answers []string `json:"answers"`
	Status  string   `json:"status"`
}

type CalcErrorResponse struct {
	Error string `json:"error"`
}
