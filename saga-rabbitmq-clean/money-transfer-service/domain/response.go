package domain

type Response struct {
	Result Result        `json:"result"`
	Trace  TraceResponse `json:"trace"`
	Data   interface{}   `json:"data"`
}

type TraceResponse struct {
	Frm string `json:"frm"`
	To  string `json:"to"`
	Cid string `json:"cid"`
	Sid string `json:"sid"`
	Cts int64  `json:"cts"`
	Sts int64  `json:"sts"`
	Dur int64  `json:"dur"`
}

type Result struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}
