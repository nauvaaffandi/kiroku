package models

type Settings map[string]string

type SettingReq struct {
	Value string `json:"value"`
}

type SettingResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
