package models

import "encoding/json"

const AppsStore string = "apps_status"

type AppService struct {
	NodeId   int32  `json:"node_id"`
	Name     string `json:"name"`
	AppId    int32  `json:"app_id"`
	Priority int32  `json:"priority"`
	Persist  int32  `json:"persist"`
	Ok       bool   `json:"ok"`
}

// ToJSON - helps to convert model to bytes or string
func (b *AppService) ToJSON() ([]byte, error) {
	return json.Marshal(b)
}

// FromJSON - helps to convert bytes to model object
func (b *AppService) FromJSON(data []byte) error {
	return json.Unmarshal(data, &b)
}
