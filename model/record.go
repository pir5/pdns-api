package model

type Record struct {
	ID        int    `json:"id"`
	DomainID  int    `json:"domain_id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	TTL       int    `json:"ttl"`
	Prio      int    `json:"prio"`
	Disabled  bool   `json:"disabled"`
	OlderName string `json:"older_name"`
	Auth      bool   `json:"auth"`
}

type Records []Record
