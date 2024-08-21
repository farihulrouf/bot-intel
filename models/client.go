package models

type ClientInfo struct {
	ID     string `json:"id"`
	Number string `json:"number,omitempty"`
	QR     string `json:"qr,omitempty"`
	Status string `json:"status"`
	Name   string `json:"name"`
	Busy   bool   `json:"busy,omitempty"`
}
