package model

type PingRequest struct {
	Data string `json:"data" validate:"required"`
}
type GenerateReq struct {
	Values map[string]interface{} `json:"values"`
	Id     string                 `json:"-"`
}
