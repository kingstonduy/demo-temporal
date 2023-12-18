package model

type Account struct {
	Cif     string  `json:"cif"`
	Balance float64 `json:"balance"`
	IsSms   bool    `json:"isSms`
	IsEmail bool    `json:"isEmail"`
}

type PaymentDetail struct {
	Id              string  `json:"id"`
	SourceAccountId string  `json:"sourceAccount"`
	TargetAccountId string  `json:"targetAccount"`
	Amount          float64 `json:"amount"`
}

type saferResponse struct {
	Response string `json:"response"`
}
