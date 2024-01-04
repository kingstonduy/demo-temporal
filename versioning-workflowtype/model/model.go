package model

type Account struct {
	Cif     string  `json:"cif"`
	Balance float64 `json:"balance"`
}

type Otp struct {
	Otp       string `json:"otp"`
	Timestamp int64  `json:"timestamp"`
}

type RequestOtp struct {
	Otp string `json:"otp"`
}

type ResponseOtp struct {
	Check bool `json:"check"`
}
