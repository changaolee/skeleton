package model

import "encoding/json"

type AuthzResponse struct {
	Allowed bool   `json:"allowed"`
	Denied  bool   `json:"denied,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (resp *AuthzResponse) ToString() string {
	data, _ := json.Marshal(resp)

	return string(data)
}
