package model

type GenerateReq struct {
	Values map[string]interface{} `json:"values"`
	Id     string                 `json:"-"`
}
