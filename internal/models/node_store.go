package models

import "encoding/json"

const NodeStore string = "brokers"

// Node - connected node detail
type Node struct {
	Bind     string `json:"bind"`
	RefId    string `json:"ref_id"`
	MyId     int32  `json:"my_id,omitempty"`
	Priority int32  `json:"priority"`
}

// ToJSON - helps to convert model to bytes or string
func (b *Node) ToJSON() ([]byte, error) {
	return json.Marshal(b)
}

// FromJSON - helps to convert bytes to model object
func (b *Node) FromJSON(data []byte) error {
	return json.Unmarshal(data, &b)
}
