package domain

type Request struct {
	Trace TraceRequest `json:"trace"`
	Data  interface{}  `json:"data"`
}

type TraceRequest struct {
	Frm string `json:"frm"`
	To  string `json:"to"`
	Cid string `json:"cid"`
	Sid string `json:"sid"`
	Cts int64  `json:"cts"`
	Sts int64  `json:"sts"`
}
