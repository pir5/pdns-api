package model

type Domain struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Master         string `json:"master"`
	LastCheck      int    `json:"last_check"`
	Type           string `json:"type"`
	NotifiedSerial int32  `json:"notified_serial"`
	Account        string `json:"account"`
}

type Domains []Domain
