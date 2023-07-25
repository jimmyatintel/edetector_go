package parsedb

import "encoding/json"

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

type ChromeLogin struct {
	UUID            string `json:"uuid"`
	Agent           string `json:"agent"`
	OriginURL       string `json:"origin_url"`
	ActionURL       string `json:"action_url"`
	UsernameElement string `json:"username_element"`
	UsernameValue   string `json:"username_value"`
	PasswordElement string `json:"password_element"`
	PasswordValue   string `json:"password_value"`
	DateCreated     string `json:"date_created"`
}

func (n ChromeLogin) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type DNSInfo struct {
	UUID           string `json:"uuid"`
	Agent          string `json:"agent"`
	TimeToLive     int    `json:"timetolive"`
	PsComputerName string `json:"pscomputername"`
	Caption        string `json:"caption"`
	Description    string `json:"description"`
	ElementName    string `json:"elementname"`
	InstanceID     int    `json:"instanceid"`
	Data           string `json:"data"`
	DataLength     int    `json:"datalength"`
	Entry          string `json:"entry"`
	Name           string `json:"name"`
	Section        int    `json:"section"`
	Status         int    `json:"status"`
	Type           int    `json:"type"`
}

func (n DNSInfo) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type EdgeBookmarks struct {
	UUID         string `json:"uuid"`
	Agent        string `json:"agent"`
	ID           int    `json:"id"`
	Parent       int    `json:"parent"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	URL          string `json:"url"`
	Source       string `json:"source"`
	GUID         string `json:"guid"`
	DateAdded    int    `json:"date_added"`
	DateModified int    `json:"date_modified"`
}

func (n EdgeBookmarks) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type EdgeCache struct {
	UUID            string `json:"uuid"`
	Agent           string `json:"agent"`
	ID              int    `json:"id"`
	URL             string `json:"url"`
	WebSite         string `json:"web_site"`
	Frame           string `json:"frame"`
	CacheControl    string `json:"cache_control"`
	ContentEncoding string `json:"content_encoding"`
	ContentLength   int    `json:"content_length"`
	ContentType     string `json:"content_type"`
	Date            int    `json:"date"`
	Expires         int    `json:"expires"`
	LastModified    int    `json:"last_modified"`
	Server          string `json:"server"`
	UsageCounter    int    `json:"usage_counter"`
	ReuseCounter    int    `json:"reuse_counter"`
}

func (n EdgeCache) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type EdgeCookies struct {
	UUID           string `json:"uuid"`
	Agent          string `json:"agent"`
	ID             int    `json:"id"`
	CreationUTC    int    `json:"creation_utc"`
	HostKey        string `json:"host_key"`
	Name           string `json:"name"`
	Value          string `json:"value"`
	EncryptedValue string `json:"encrypted_value"`
	ExpiresUTC     int    `json:"expires_utc"`
	LastAccessUTC  int    `json:"last_access_utc"`
	SourcePort     int    `json:"source_port"`
}

func (n EdgeCookies) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type EdgeHistory struct {
	UUID          string `json:"uuid"`
	Agent         string `json:"agent"`
	ID            int    `json:"id"`
	URL           string `json:"url"`
	Title         string `json:"title"`
	VisitTime     int    `json:"visit_time"`
	VisitCount    int    `json:"visit_count"`
	LastVisitTime int    `json:"last_visit_time"`
}

func (n EdgeHistory) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type EdgeLogin struct {
	UUID            string `json:"uuid"`
	Agent           string `json:"agent"`
	ID              int    `json:"id"`
	OriginURL       string `json:"origin_url"`
	ActionURL       string `json:"action_url"`
	UsernameElement string `json:"username_element"`
	UsernameValue   string `json:"username_value"`
	PasswordElement string `json:"password_element"`
	PasswordValue   string `json:"password_value"`
	DateCreated     int    `json:"date_created"`
}

func (n EdgeLogin) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type EventApplication struct {
	UUID                         string `json:"uuid"`
	Agent                        string `json:"agent"`
	EventRecordID                int    `json:"eventrecordid"`
	ProviderName                 string `json:"providername"`
	ProviderGUID                 string `json:"providerguid"`
	EventID                      int    `json:"eventid"`
	Version                      int    `json:"version"`
	Level                        string `json:"level"`
	Keywords                     string `json:"keywords"`
	Task                         int    `json:"task"`
	Opcode                       int    `json:"opcode"`
	CreatedSystemTime            int    `json:"createdsystemtime"`
	CorrelationActivityID        string `json:"correlationactivityid"`
	CorrelationRelatedActivityID string `json:"correlationrelatedactivityid"`
	ExecutionProcessID           int    `json:"executionprocessid"`
	ExecutionThreadID            int    `json:"executionthreadid"`
	Channel                      string `json:"channel"`
	Computer                     string `json:"computer"`
	SecurityUserID               string `json:"securityuserid"`
	EvtRenderData                string `json:"evtrenderdata"`
}

func (n EventApplication) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type EventSecurity struct {
	UUID                         string `json:"uuid"`
	Agent                        string `json:"agent"`
	EventRecordID                int    `json:"eventrecordid"`
	ProviderName                 string `json:"providername"`
	ProviderGUID                 string `json:"providerguid"`
	EventID                      int    `json:"eventid"`
	Version                      int    `json:"version"`
	Level                        string `json:"level"`
	Keywords                     string `json:"keywords"`
	Task                         int    `json:"task"`
	Opcode                       int    `json:"opcode"`
	CreatedSystemTime            int    `json:"createdsystemtime"`
	CorrelationActivityID        string `json:"correlationactivityid"`
	CorrelationRelatedActivityID string `json:"correlationrelatedactivityid"`
	ExecutionProcessID           int    `json:"executionprocessid"`
	ExecutionThreadID            int    `json:"executionthreadid"`
	Channel                      string `json:"channel"`
	Computer                     string `json:"computer"`
	SecurityUserID               string `json:"securityuserid"`
	EvtRenderData                string `json:"evtrenderdata"`
}

func (n EventSecurity) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type EventSystem struct {
	UUID                         string `json:"uuid"`
	Agent                        string `json:"agent"`
	EventRecordID                int    `json:"eventrecordid"`
	ProviderName                 string `json:"providername"`
	ProviderGUID                 string `json:"providerguid"`
	EventID                      int    `json:"eventid"`
	Version                      int    `json:"version"`
	Level                        string `json:"level"`
	Keywords                     string `json:"keywords"`
	Task                         int    `json:"task"`
	Opcode                       int    `json:"opcode"`
	CreatedSystemTime            int    `json:"createdsystemtime"`
	CorrelationActivityID        string `json:"correlationactivityid"`
	CorrelationRelatedActivityID string `json:"correlationrelatedactivityid"`
	ExecutionProcessID           int    `json:"executionprocessid"`
	ExecutionThreadID            int    `json:"executionthreadid"`
	Channel                      string `json:"channel"`
	Computer                     string `json:"computer"`
	SecurityUserID               string `json:"securityuserid"`
	EvtRenderData                string `json:"evtrenderdata"`
}

func (n EventSystem) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type FirefoxBookmarks struct {
	UUID             string `json:"uuid"`
	Agent            string `json:"agent"`
	ID               int    `json:"id"`
	Parent           int    `json:"parent"`
	Type             int    `json:"type"`
	Title            string `json:"title"`
	DateAdded        int    `json:"dateadded"`
	LastModified     int    `json:"lastmodified"`
	GUID             string `json:"guid"`
	SyncStatus       int    `json:"syncstatus"`
	URL              string `json:"url"`
	URLTitle         string `json:"url_title"`
	RevHost          string `json:"rev_host"`
	URLVisitCount    int    `json:"url_visit_count"`
	URLLastVisitDate int    `json:"url_last_visit_date"`
	URLGUID          string `json:"url_guid"`
	URLDescription   string `json:"url_description"`
	PreviewImageURL  string `json:"preview_image_url"`
}

func (n FirefoxBookmarks) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

