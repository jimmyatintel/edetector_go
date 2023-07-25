package parsedb

import "encoding/json"

type ARPCache struct {
	UUID            string `json:"uuid"`
	Agent           string `json:"agent"`
	Interface       string `json:"interface"`
	Internetaddress string `json:"internetaddress"`
	Physicaladdress string `json:"physicaladdress"`
	Type            string `json:"type"`
}

func (n ARPCache) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type AppResourceUsageMonitor struct {
	UUID                         string `json:"uuid"`
	Agent                        string `json:"agent"`
	RecordId                     int    `json:"record_id"`
	App_name                     string `json:"app_name"`
	App_id                       int    `json:"app_id"`
	User_name                    string `json:"user_name"`
	User_sid                     string `json:"user_sid"`
	Foregroundcycletime          int    `json:"foregroundcycletime"`
	Backgroundcycletime          int    `json:"backgroundcycletime"`
	Facetime                     int    `json:"facetime"`
	Foregroundbytesread          int    `json:"foregroundbytesread"`
	Foregroundbyteswritten       int    `json:"foregroundbyteswritten"`
	Foregroundnumreadoperations  int    `json:"foregroundnumreadoperations"`
	Foregroundnumwriteoperations int    `json:"foregroundnumwriteoperations"`
	Foregroundnumberofflushes    int    `json:"foregroundnumberofflushes"`
	Backgroundbytesread          int    `json:"backgroundbytesread"`
	Backgroundbyteswritten       int    `json:"backgroundbyteswritten"`
	Backgroundnumreadoperations  int    `json:"backgroundnumreadoperations"`
	Backgroundnumwriteoperations int    `json:"backgroundnumwriteoperations"`
	Backgroundnumberofflushes    int    `json:"backgroundnumberofflushes"`
	Interfaceluid                string `json:"interfaceluid"`
	Timestamp                    int    `json:"timestamp"`
}

func (n AppResourceUsageMonitor) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type BaseService struct {
	UUID              string `json:"uuid"`
	Agent             string `json:"agent"`
	Name              string `json:"name"`
	Caption           string `json:"caption"`
	Description       string `json:"description"`
	Displayname       string `json:"displayname"`
	Errorcontrol      string `json:"errorcontrol"`
	Pathname          string `json:"pathname"`
	Creationclassname string `json:"creationclassname"`
	ServiceType       string `json:"servicetype"`
	Started           string `json:"started"`
	StartMode         string `json:"startmode"`
	Startname         string `json:"startname"`
	State             string `json:"state"`
	Status            string `json:"status"`
	Systemname        string `json:"systemname"`
	Acceptpause       bool   `json:"acceptpause"`
	Acceptstop        bool   `json:"acceptstop"`
}

func (n BaseService) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type ChromeBookmarks struct {
	UUID          string `json:"uuid"`
	Agent         string `json:"agent"`
	Id            int    `json:"id"`
	Parent        int    `json:"parent"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	Url           string `json:"url"`
	Guid          string `json:"guid"`
	Date_added    int    `json:"date_added"`
	Date_modified int    `json:"date_modified"`
}

func (n ChromeBookmarks) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type ChromeCache struct {
	UUID             string `json:"uuid"`
	Agent            string `json:"agent"`
	Id               int    `json:"id"`
	Url              string `json:"url"`
	Web_site         string `json:"web_site"`
	Frame            string `json:"frame"`
	Cache_control    string `json:"cache_control"`
	Content_encoding string `json:"content_encoding"`
	Content_length   int    `json:"content_length"`
	Conent_type      string `json:"ententent_type"`
	Date             int    `json:"date"`
	Expires          int    `json:"expires"`
	Last_modified    int    `json:"last_modified"`
	Server           string `json:"server"`
	Usage_counter    int    `json:"usage_counter"`
	Reuse_counter    int    `json:"reuse_counter"`
}

func (n ChromeCache) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type ChromeDownload struct {
	UUID             string `json:"uuid"`
	Agent            string `json:"agent"`
	DownloadURL      string `json:"download_url"`
	Guid             string `json:"guid"`
	CurrentPath      string `json:"current_path"`
	TargetPath       string `json:"target_path"`
	ReceivedBytes    int    `json:"received_bytes"`
	TotalBytes       int    `json:"total_bytes"`
	StartTime        int    `json:"start_time"`
	EndTime          int    `json:"end_time"`
	LastAccessTime   int    `json:"last_access_time"`
	Danger           string `json:"danger"`
	InterruptReason  string `json:"interrupt_reason"`
	Opened           bool   `json:"opened"`
	Referrer         string `json:"referrer"`
	SiteURL          string `json:"site_url"`
	TabURL           string `json:"tab_url"`
	TabReferrerURL   string `json:"tab_referrer_url"`
	ETag             string `json:"etag"`
	LastModified     int    `json:"last_modified"`
	MimeType         string `json:"mime_type"`
	OriginalMimeType string `json:"original_mime_type"`
}

func (n ChromeDownload) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type ChromeHistory struct {
	UUID          string `json:"uuid"`
	Agent         string `json:"agent"`
	URL           string `json:"url"`
	Title         string `json:"title"`
	VisitTime     int    `json:"visit_time"`
	VisitCount    int    `json:"visit_count"`
	LastVisitTime int    `json:"last_visit_time"`
}

func (n ChromeHistory) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type ChromeKeywordSearch struct {
	Term  string `json:"term"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

func (n ChromeKeywordSearch) Elastical() ([]byte, error) {
	return json.Marshal(n)
}
