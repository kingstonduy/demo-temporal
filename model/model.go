package model

type Account struct {
	Cif     string  `json:"cif"`
	Balance float64 `json:"balance"`
	IsSms   bool    `json:"isSms`
	IsEmail bool    `json:"isEmail"`
}

type ParallelWorkflowInput struct {
	Cif1 string
	Cif2 string
}
