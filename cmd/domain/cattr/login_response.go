package cattr

type LoginResponse struct {
	Data DataLoginResponse `json:"data"`
}

type DataLoginResponse struct {
	AccessToken string `json:"access_token"`
}
