package parsedb

import "encoding/json"

type KeyColumn struct {
	Item string `json:"item"`
	Date string `json:"date"`
	Type string `json:"type"`
	Etc  string `json:"etc"`
}

type ARPCache struct {
	Interface       string `json:"interface"`
	Internetaddress string `json:"internetaddress"`
	Physicaladdress string `json:"physicaladdress"`
	Type            string `json:"type"`
}

func (n ARPCache) Elastical() ([]byte, error) {
	return json.Marshal(n)
}
