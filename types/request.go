package types

type LoginRequest struct {
	Email  string `json:"email"`
	Passwd string `json:"passwd"`
}

type Result struct {
}

type LoginResp struct {
	Ret   int    `json:"ret"`
	Token string `json:"token"`
}
