package treebuilder

import "encoding/json"

type ExplorerRelation struct {
	Agent  string   `json:"agent"`
	IsRoot bool     `json:"isRoot"`
	Parent string   `json:"parent"`
	Child  []string `json:"child"`
}

func (n ExplorerRelation) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type ExplorerDetails struct {
	UUID              string `json:"uuid"`
	Agent             string `json:"agent"`
	AgentIP           string `json:"agentIP"`
	AgentName         string `json:"agentName"`
	FileName          string `json:"fileName"`
	IsDeleted         bool   `json:"isDeleted"`
	IsDirectory       bool   `json:"isDirectory"`
	CreateTime        int    `json:"createTime"`
	WriteTime         int    `json:"writeTime"`
	AccessTime        int    `json:"accessTime"`
	EntryModifiedTime int    `json:"entryModifiedTime"`
	Datalen           int    `json:"dataLen"`
	ItemMain          string `json:"item_main"`
	DateMain          int    `json:"date_main"`
	TypeMain          string `json:"type_main"`
	EtcMain           string `json:"etc_main"`
}

func (n ExplorerDetails) Elastical() ([]byte, error) {
	return json.Marshal(n)
}
