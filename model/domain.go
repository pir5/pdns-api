package model

import "strings"

type Domain struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Master         string  `json:"master"`
	LastCheck      int     `json:"last_check"`
	Type           string  `json:"type"`
	NotifiedSerial int32   `json:"notified_serial"`
	Account        string  `json:"account"`
	Records        Records `json:"records"`
}

type Domains []Domain

func (ds Domains) FilterOwnerDomains(domains []string) {
	ret := Domains{}
	for _, v := range ds {
		for _, vv := range domains {
			if strings.ToLower(v.Name) == strings.ToLower(vv) {
				ret = append(ret, v)
			}
		}
	}
	ds = ret
}
