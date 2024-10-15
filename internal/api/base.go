package api

type BaseReq struct {
}

type BaseRsp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
