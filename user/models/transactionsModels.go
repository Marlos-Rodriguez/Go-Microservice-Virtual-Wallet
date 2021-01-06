package models

//TransactionResponse struct
type TransactionResponse struct {
	TsID      string  `json:"tsID"`
	FromUser  string  `json:"fromId"`
	FromName  string  `json:"fromName"`
	ToUser    string  `json:"toId"`
	ToName    string  `json:"toName"`
	Amount    float32 `json:"amount"`
	Message   string  `json:"message,omitempty"`
	CreatedAt string  `json:"createAt"`
}
