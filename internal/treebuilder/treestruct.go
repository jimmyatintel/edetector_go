package treebuilder

import "encoding/json"

type ExplorerDetails struct {
	FileName          string   `json:"fileName"`
	IsDeleted         bool     `json:"isDeleted"`
	IsDirectory       bool     `json:"isDirectory"`
	CreateTime        int      `json:"createTime"`
	WriteTime         int      `json:"writeTime"`
	AccessTime        int      `json:"accessTime"`
	EntryModifiedTime int      `json:"entryModifiedTime"`
	Datalen           int64    `json:"dataLen"`
	Path              string   `json:"path"`
	Disk              string   `json:"disk"`
	MD5_Sig           string   `json:"md5_sig"`
	YaraRuleHitCount  int      `json:"yaraRuleHitCount"`
	YaraRuleHit       string   `json:"yaraRuleHit"`
	IsRoot            bool     `json:"isRoot"`
	Child             []string `json:"child"`
	UUID              string   `json:"uuid"`
	Agent             string   `json:"agent"`
	AgentIP           string   `json:"agentIP"`
	AgentName         string   `json:"agentName"`
	ItemMain          string   `json:"item_main"`
	DateMain          int      `json:"date_main"`
	TypeMain          string   `json:"type_main"`
	EtcMain           string   `json:"etc_main"`
	Task_id           string   `json:"task_id"`
}

func (n ExplorerDetails) Elastical() ([]byte, error) {
	return json.Marshal(n)
}
