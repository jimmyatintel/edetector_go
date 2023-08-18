package dbparser

import "encoding/json"

type AppResourceUsageMonitor struct {
	UUID                         string `json:"uuid"`
	Agent                        string `json:"agent"`
	AgentIP                      string `json:"agentIP"`
	AgentName                    string `json:"agentName"`
	RecordId                     int    `json:"record_id"`
	App_name                     string `json:"app_name"`
	App_id                       int    `json:"app_id"`
	User_name                    string `json:"user_name"`
	User_sid                     string `json:"user_sid"`
	Foregroundcycletime          int64  `json:"foregroundcycletime"`
	Backgroundcycletime          int64  `json:"backgroundcycletime"`
	Facetime                     int64  `json:"facetime"`
	Foregroundbytesread          int64  `json:"foregroundbytesread"`
	Foregroundbyteswritten       int64  `json:"foregroundbyteswritten"`
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
	ItemMain                     string `json:"item_main"`
	DateMain                     int    `json:"date_main"`
	TypeMain                     string `json:"type_main"`
	EtcMain                      string `json:"etc_main"`
}

func (n AppResourceUsageMonitor) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type ARPCache struct {
	UUID            string `json:"uuid"`
	Agent           string `json:"agent"`
	AgentIP         string `json:"agentIP"`
	AgentName       string `json:"agentName"`
	Interface       string `json:"interface"`
	Internetaddress string `json:"internetaddress"`
	Physicaladdress string `json:"physicaladdress"`
	Type            string `json:"type"`
	ItemMain        string `json:"item_main"`
	DateMain        int    `json:"date_main"`
	TypeMain        string `json:"type_main"`
	EtcMain         string `json:"etc_main"`
}

func (n ARPCache) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type BaseService struct {
	UUID              string `json:"uuid"`
	Agent             string `json:"agent"`
	AgentIP           string `json:"agentIP"`
	AgentName         string `json:"agentName"`
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
	ItemMain          string `json:"item_main"`
	DateMain          int    `json:"date_main"`
	TypeMain          string `json:"type_main"`
	EtcMain           string `json:"etc_main"`
}

func (n BaseService) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type ChromeBookmarks struct {
	UUID          string `json:"uuid"`
	Agent         string `json:"agent"`
	AgentIP       string `json:"agentIP"`
	AgentName     string `json:"agentName"`
	Id            int    `json:"id"`
	Parent        int    `json:"parent"`
	Type          string `json:"ttype"`
	Name          string `json:"name"`
	Url           string `json:"url"`
	Guid          string `json:"guid"`
	Date_added    int    `json:"date_added"`
	Date_modified int    `json:"date_modified"`
	ItemMain      string `json:"item_main"`
	DateMain      int    `json:"date_main"`
	TypeMain      string `json:"type_main"`
	EtcMain       string `json:"etc_main"`
}

func (n ChromeBookmarks) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type ChromeCache struct {
	UUID             string `json:"uuid"`
	Agent            string `json:"agent"`
	AgentIP          string `json:"agentIP"`
	AgentName        string `json:"agentName"`
	Id               int    `json:"id"`
	Url              string `json:"url"`
	Web_site         string `json:"web_site"`
	Frame            string `json:"frame"`
	Cache_control    string `json:"cache_control"`
	Content_encoding string `json:"content_encoding"`
	Content_length   int    `json:"content_length"`
	Conent_type      string `json:"content_type"`
	Date             int    `json:"date"`
	Expires          int    `json:"expires"`
	Last_modified    int    `json:"last_modified"`
	Server           string `json:"server"`
	Usage_counter    int    `json:"usage_counter"`
	Reuse_counter    int    `json:"reuse_counter"`
	ItemMain         string `json:"item_main"`
	DateMain         int    `json:"date_main"`
	TypeMain         string `json:"type_main"`
	EtcMain          string `json:"etc_main"`
}

func (n ChromeCache) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type ChromeDownload struct {
	UUID             string `json:"uuid"`
	Agent            string `json:"agent"`
	AgentIP          string `json:"agentIP"`
	AgentName        string `json:"agentName"`
	DownloadURL      string `json:"download_url"`
	Guid             string `json:"guid"`
	CurrentPath      string `json:"current_path"`
	TargetPath       string `json:"target_path"`
	ReceivedBytes    int64  `json:"received_bytes"`
	TotalBytes       int64  `json:"total_bytes"`
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
	ItemMain         string `json:"item_main"`
	DateMain         int    `json:"date_main"`
	TypeMain         string `json:"type_main"`
	EtcMain          string `json:"etc_main"`
}

func (n ChromeDownload) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type ChromeHistory struct {
	UUID          string `json:"uuid"`
	Agent         string `json:"agent"`
	AgentIP       string `json:"agentIP"`
	AgentName     string `json:"agentName"`
	URL           string `json:"url"`
	Title         string `json:"title"`
	VisitTime     int    `json:"visit_time"`
	VisitCount    int    `json:"visit_count"`
	LastVisitTime int    `json:"last_visit_time"`
	ItemMain      string `json:"item_main"`
	DateMain      int    `json:"date_main"`
	TypeMain      string `json:"type_main"`
	EtcMain       string `json:"etc_main"`
}

func (n ChromeHistory) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type ChromeKeywordSearch struct {
	UUID      string `json:"uuid"`
	Agent     string `json:"agent"`
	AgentIP   string `json:"agentIP"`
	AgentName string `json:"agentName"`
	Term      string `json:"term"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	ItemMain  string `json:"item_main"`
	DateMain  int    `json:"date_main"`
	TypeMain  string `json:"type_main"`
	EtcMain   string `json:"etc_main"`
}

func (n ChromeKeywordSearch) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type ChromeLogin struct {
	UUID            string `json:"uuid"`
	Agent           string `json:"agent"`
	AgentIP         string `json:"agentIP"`
	AgentName       string `json:"agentName"`
	OriginURL       string `json:"origin_url"`
	ActionURL       string `json:"action_url"`
	UsernameElement string `json:"username_element"`
	UsernameValue   string `json:"username_value"`
	PasswordElement string `json:"password_element"`
	PasswordValue   string `json:"password_value"`
	DateCreated     string `json:"date_created"`
	ItemMain        string `json:"item_main"`
	DateMain        int    `json:"date_main"`
	TypeMain        string `json:"type_main"`
	EtcMain         string `json:"etc_main"`
}

func (n ChromeLogin) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type DNSInfo struct {
	UUID           string `json:"uuid"`
	Agent          string `json:"agent"`
	AgentIP        string `json:"agentIP"`
	AgentName      string `json:"agentName"`
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
	ItemMain       string `json:"item_main"`
	DateMain       int    `json:"date_main"`
	TypeMain       string `json:"type_main"`
	EtcMain        string `json:"etc_main"`
}

func (n DNSInfo) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type EdgeBookmarks struct {
	UUID         string `json:"uuid"`
	Agent        string `json:"agent"`
	AgentIP      string `json:"agentIP"`
	AgentName    string `json:"agentName"`
	ID           int    `json:"id"`
	Parent       int    `json:"parent"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	URL          string `json:"url"`
	Source       string `json:"source"`
	GUID         string `json:"guid"`
	DateAdded    int    `json:"date_added"`
	DateModified int    `json:"date_modified"`
	ItemMain     string `json:"item_main"`
	DateMain     int    `json:"date_main"`
	TypeMain     string `json:"type_main"`
	EtcMain      string `json:"etc_main"`
}

func (n EdgeBookmarks) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type EdgeCache struct {
	UUID            string `json:"uuid"`
	Agent           string `json:"agent"`
	AgentIP         string `json:"agentIP"`
	AgentName       string `json:"agentName"`
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
	ItemMain        string `json:"item_main"`
	DateMain        int    `json:"date_main"`
	TypeMain        string `json:"type_main"`
	EtcMain         string `json:"etc_main"`
}

func (n EdgeCache) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type EdgeCookies struct {
	UUID           string `json:"uuid"`
	Agent          string `json:"agent"`
	AgentIP        string `json:"agentIP"`
	AgentName      string `json:"agentName"`
	ID             int    `json:"id"`
	CreationUTC    int    `json:"creation_utc"`
	HostKey        string `json:"host_key"`
	Name           string `json:"name"`
	Value          string `json:"value"`
	EncryptedValue string `json:"encrypted_value"`
	ExpiresUTC     int    `json:"expires_utc"`
	LastAccessUTC  int    `json:"last_access_utc"`
	SourcePort     int    `json:"source_port"`
	ItemMain       string `json:"item_main"`
	DateMain       int    `json:"date_main"`
	TypeMain       string `json:"type_main"`
	EtcMain        string `json:"etc_main"`
}

func (n EdgeCookies) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type EdgeHistory struct {
	UUID          string `json:"uuid"`
	Agent         string `json:"agent"`
	AgentIP       string `json:"agentIP"`
	AgentName     string `json:"agentName"`
	ID            int    `json:"id"`
	URL           string `json:"url"`
	Title         string `json:"title"`
	VisitTime     int    `json:"visit_time"`
	VisitCount    int    `json:"visit_count"`
	LastVisitTime int    `json:"last_visit_time"`
	ItemMain      string `json:"item_main"`
	DateMain      int    `json:"date_main"`
	TypeMain      string `json:"type_main"`
	EtcMain       string `json:"etc_main"`
}

func (n EdgeHistory) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type EdgeLogin struct {
	UUID            string `json:"uuid"`
	Agent           string `json:"agent"`
	AgentIP         string `json:"agentIP"`
	AgentName       string `json:"agentName"`
	ID              int    `json:"id"`
	OriginURL       string `json:"origin_url"`
	ActionURL       string `json:"action_url"`
	UsernameElement string `json:"username_element"`
	UsernameValue   string `json:"username_value"`
	PasswordElement string `json:"password_element"`
	PasswordValue   string `json:"password_value"`
	DateCreated     int    `json:"date_created"`
	ItemMain        string `json:"item_main"`
	DateMain        int    `json:"date_main"`
	TypeMain        string `json:"type_main"`
	EtcMain         string `json:"etc_main"`
}

func (n EdgeLogin) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type EventApplication struct {
	UUID                         string `json:"uuid"`
	Agent                        string `json:"agent"`
	AgentIP                      string `json:"agentIP"`
	AgentName                    string `json:"agentName"`
	EventRecordID                int    `json:"eventrecordid"`
	ProviderName                 string `json:"providername"`
	ProviderGUID                 string `json:"providerguid"`
	EventID                      int64  `json:"eventid"`
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
	ItemMain                     string `json:"item_main"`
	DateMain                     int    `json:"date_main"`
	TypeMain                     string `json:"type_main"`
	EtcMain                      string `json:"etc_main"`
}

func (n EventApplication) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type EventSecurity struct {
	UUID                         string `json:"uuid"`
	Agent                        string `json:"agent"`
	AgentIP                      string `json:"agentIP"`
	AgentName                    string `json:"agentName"`
	EventRecordID                int    `json:"eventrecordid"`
	ProviderName                 string `json:"providername"`
	ProviderGUID                 string `json:"providerguid"`
	EventID                      int64  `json:"eventid"`
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
	ItemMain                     string `json:"item_main"`
	DateMain                     int    `json:"date_main"`
	TypeMain                     string `json:"type_main"`
	EtcMain                      string `json:"etc_main"`
}

func (n EventSecurity) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type EventSystem struct {
	UUID                         string `json:"uuid"`
	Agent                        string `json:"agent"`
	AgentIP                      string `json:"agentIP"`
	AgentName                    string `json:"agentName"`
	EventRecordID                int    `json:"eventrecordid"`
	ProviderName                 string `json:"providername"`
	ProviderGUID                 string `json:"providerguid"`
	EventID                      int64  `json:"eventid"`
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
	ItemMain                     string `json:"item_main"`
	DateMain                     int    `json:"date_main"`
	TypeMain                     string `json:"type_main"`
	EtcMain                      string `json:"etc_main"`
}

func (n EventSystem) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type FirefoxBookmarks struct {
	UUID             string `json:"uuid"`
	Agent            string `json:"agent"`
	AgentIP          string `json:"agentIP"`
	AgentName        string `json:"agentName"`
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
	ItemMain         string `json:"item_main"`
	DateMain         int    `json:"date_main"`
	TypeMain         string `json:"type_main"`
	EtcMain          string `json:"etc_main"`
}

func (n FirefoxBookmarks) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type FirefoxCache struct {
	UUID            string `json:"uuid"`
	Agent           string `json:"agent"`
	AgentIP         string `json:"agentIP"`
	AgentName       string `json:"agentName"`
	ID              int    `json:"id"`
	URL             string `json:"url"`
	ServerResponse  string `json:"server_response"`
	ServerName      string `json:"server_name"`
	CacheControl    string `json:"cache_control"`
	ContentEncoding string `json:"content_encoding"`
	ContentLength   int    `json:"content_length"`
	ContentType     string `json:"content_type"`
	FetchCount      int    `json:"fetch_count"`
	LastFetched     int    `json:"last_fetched"`
	LastModified    int    `json:"last_modified"`
	Frequency       int    `json:"frequency"`
	Expiration      int    `json:"expiration"`
	ItemMain        string `json:"item_main"`
	DateMain        int    `json:"date_main"`
	TypeMain        string `json:"type_main"`
	EtcMain         string `json:"etc_main"`
}

func (n FirefoxCache) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type FirefoxCookies struct {
	UUID         string `json:"uuid"`
	Agent        string `json:"agent"`
	AgentIP      string `json:"agentIP"`
	AgentName    string `json:"agentName"`
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Value        string `json:"value"`
	Host         string `json:"host"`
	Path         string `json:"path"`
	LastAccessed int    `json:"lastaccessed"`
	CreationTime int    `json:"creationtime"`
	ItemMain     string `json:"item_main"`
	DateMain     int    `json:"date_main"`
	TypeMain     string `json:"type_main"`
	EtcMain      string `json:"etc_main"`
}

func (n FirefoxCookies) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type FirefoxHistory struct {
	UUID          string `json:"uuid"`
	Agent         string `json:"agent"`
	AgentIP       string `json:"agentIP"`
	AgentName     string `json:"agentName"`
	URL           string `json:"url"`
	Title         string `json:"title"`
	FromURL       string `json:"from_url"`
	RevHost       string `json:"rev_host"`
	GUID          string `json:"guid"`
	Description   string `json:"description"`
	PreviewImgURL string `json:"preview_image_url"`
	VisitCount    int    `json:"visit_count"`
	VisitDate     int    `json:"visit_date"`
	LastVisitDate int    `json:"last_visit_date"`
	ItemMain      string `json:"item_main"`
	DateMain      int    `json:"date_main"`
	TypeMain      string `json:"type_main"`
	EtcMain       string `json:"etc_main"`
}

func (n FirefoxHistory) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type IEHistory struct {
	UUID            string `json:"uuid"`
	Agent           string `json:"agent"`
	AgentIP         string `json:"agentIP"`
	AgentName       string `json:"agentName"`
	URL             string `json:"url"`
	Title           string `json:"title"`
	ExpiresTime     int    `json:"expirestime"`
	LastUpdatedTime int    `json:"lastupdatedtime"`
	VisitedTime     int    `json:"visitedtime"`
	ItemMain        string `json:"item_main"`
	DateMain        int    `json:"date_main"`
	TypeMain        string `json:"type_main"`
	EtcMain         string `json:"etc_main"`
}

func (n IEHistory) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type InstalledSoftware struct {
	UUID                      string `json:"uuid"`
	Agent                     string `json:"agent"`
	AgentIP                   string `json:"agentIP"`
	AgentName                 string `json:"agentName"`
	DisplayName               string `json:"displayname"`
	RegistryName              string `json:"registryname"`
	DisplayVersion            string `json:"displayversion"`
	InstallDate               int    `json:"installdate"`
	InstalledFor              string `json:"installedfor"`
	InstallLocation           string `json:"installlocation"`
	Publisher                 string `json:"publisher"`
	UninstallString           string `json:"uninstallstring"`
	ModifyPath                string `json:"modifypath"`
	Comments                  string `json:"comments"`
	URLInfoAbout              string `json:"urlinfoabout"`
	URLUpdateInfo             string `json:"urlupdateinfo"`
	HelpLink                  string `json:"helplink"`
	InstallSource             string `json:"installsource"`
	ReleaseType               string `json:"releasetype"`
	DisplayIcon               string `json:"displayicon"`
	EstimatedSize             int    `json:"estimatedsize"`
	RegistryTime              int    `json:"registrytime"`
	InstallFolderCreatedTime  int    `json:"installfoldercreatedtime"`
	InstallFolderModifiedTime int    `json:"installfoldermodifiedtime"`
	ItemMain                  string `json:"item_main"`
	DateMain                  int    `json:"date_main"`
	TypeMain                  string `json:"type_main"`
	EtcMain                   string `json:"etc_main"`
}

func (n InstalledSoftware) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type JumpList struct {
	UUID          string `json:"uuid"`
	Agent         string `json:"agent"`
	AgentIP       string `json:"agentIP"`
	AgentName     string `json:"agentName"`
	FullPath      string `json:"fullpath"`
	ApplicationID string `json:"application_id"`
	ComputerName  string `json:"computer_name"`
	FileSize      int64  `json:"filesize"`
	EntryID       int    `json:"entry_id"`
	RecordTime    int    `json:"recordtime"`
	CreateTime    int    `json:"createtime"`
	AccessTime    int    `json:"accesstime"`
	ModifiedTime  int    `json:"modifiedtime"`
	ItemMain      string `json:"item_main"`
	DateMain      int    `json:"date_main"`
	TypeMain      string `json:"type_main"`
	EtcMain       string `json:"etc_main"`
}

func (n JumpList) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type MUICache struct {
	UUID            string `json:"uuid"`
	Agent           string `json:"agent"`
	AgentIP         string `json:"agentIP"`
	AgentName       string `json:"agentName"`
	ApplicationPath string `json:"applicationpath"`
	ApplicationName string `json:"applicationname"`
	ItemMain        string `json:"item_main"`
	DateMain        int    `json:"date_main"`
	TypeMain        string `json:"type_main"`
	EtcMain         string `json:"etc_main"`
}

func (n MUICache) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type Network struct {
	UUID           string `json:"uuid"`
	Agent          string `json:"agent"`
	AgentIP        string `json:"agentIP"`
	AgentName      string `json:"agentName"`
	ProcessID      int    `json:"processid"`
	ProcessName    string `json:"processname"`
	LocalAddress   string `json:"localaddress"`
	LocalPort      int    `json:"localport"`
	RemoteAddress  string `json:"remoteaddress"`
	RemotePort     int    `json:"remoteport"`
	State          string `json:"state"`
	RemoteHostname string `json:"remotehostname"`
	Protocol       string `json:"protocol"`
	ItemMain       string `json:"item_main"`
	DateMain       int    `json:"date_main"`
	TypeMain       string `json:"type_main"`
	EtcMain        string `json:"etc_main"`
}

func (n Network) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type NetworkDataUsageMonitor struct {
	UUID               string `json:"uuid"`
	Agent              string `json:"agent"`
	AgentIP            string `json:"agentIP"`
	AgentName          string `json:"agentName"`
	RecordID           int    `json:"record_id"`
	AppName            string `json:"app_name"`
	AppID              int    `json:"app_id"`
	UserName           string `json:"user_name"`
	UserSID            string `json:"user_sid"`
	BytesSent          int    `json:"bytes_sent"`
	BytesReceived      int64  `json:"bytes_recvd"`
	NetworkAdapter     string `json:"network_adapter"`
	NetworkAdapterGUID string `json:"network_adapter_guid"`
	InterfaceLUID      string `json:"interfaceluid"`
	Timestamp          int    `json:"timestamp"`
	ItemMain           string `json:"item_main"`
	DateMain           int    `json:"date_main"`
	TypeMain           string `json:"type_main"`
	EtcMain            string `json:"etc_main"`
}

func (n NetworkDataUsageMonitor) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type NetworkResources struct {
	UUID          string `json:"uuid"`
	Agent         string `json:"agent"`
	AgentIP       string `json:"agentIP"`
	AgentName     string `json:"agentName"`
	ResourcesName string `json:"resourcesname"`
	ResourcesType string `json:"resourcestype"`
	Comment       string `json:"comment"`
	LocalPath     string `json:"localpath"`
	Provider      string `json:"provider"`
	Scope         string `json:"scope"`
	DisplayType   string `json:"displaytype"`
	Usage         string `json:"usage"`
	IPAddress     string `json:"ipaddress"`
	ItemMain      string `json:"item_main"`
	DateMain      int    `json:"date_main"`
	TypeMain      string `json:"type_main"`
	EtcMain       string `json:"etc_main"`
}

func (n NetworkResources) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type OpenedFiles struct {
	UUID        string `json:"uuid"`
	Agent       string `json:"agent"`
	AgentIP     string `json:"agentIP"`
	AgentName   string `json:"agentName"`
	ProcessID   int    `json:"processid"`
	ProcessName string `json:"processname"`
	Type        string `json:"type"`
	ObjectName  string `json:"objectname"`
	ItemMain    string `json:"item_main"`
	DateMain    int    `json:"date_main"`
	TypeMain    string `json:"type_main"`
	EtcMain     string `json:"etc_main"`
}

func (n OpenedFiles) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type Prefetch struct {
	UUID               string `json:"uuid"`
	Agent              string `json:"agent"`
	AgentIP            string `json:"agentIP"`
	AgentName          string `json:"agentName"`
	FileName           string `json:"filename"`
	ProcessName        string `json:"processname"`
	LastRunTime        int    `json:"lastruntime"`
	ProcessPath        string `json:"processpath"`
	RunCount           int    `json:"runcount"`
	FileSize           int64  `json:"filesize"`
	FolderCreatedTime  int    `json:"foldercreatedtime"`
	FolderModifiedTime int    `json:"foldermodifiedtime"`
	ItemMain           string `json:"item_main"`
	DateMain           int    `json:"date_main"`
	TypeMain           string `json:"type_main"`
	EtcMain            string `json:"etc_main"`
}

func (n Prefetch) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type Process struct {
	UUID              string `json:"uuid"`
	Agent             string `json:"agent"`
	AgentIP           string `json:"agentIP"`
	AgentName         string `json:"agentName"`
	PID               int    `json:"pid"`
	ProcessName       string `json:"process_name"`
	ParentPID         int    `json:"parent_pid"`
	ProcessCreateTime int    `json:"processcreatetime"`
	ProcessPath       string `json:"process_path"`
	ProcessCommand    string `json:"process_command"`
	UserName          string `json:"user_name"`
	DigitalSignature  string `json:"digitalsignature"`
	ProductName       string `json:"productname"`
	FileVersion       string `json:"fileversion"`
	FileDescription   string `json:"filedescription"`
	CompanyName       string `json:"companyname"`
	Priority          string `json:"priority"`
	ItemMain          string `json:"item_main"`
	DateMain          int    `json:"date_main"`
	TypeMain          string `json:"type_main"`
	EtcMain           string `json:"etc_main"`
}

func (n Process) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type Service struct {
	UUID         string `json:"uuid"`
	Agent        string `json:"agent"`
	AgentIP      string `json:"agentIP"`
	AgentName    string `json:"agentName"`
	Name         string `json:"name"`
	Caption      string `json:"caption"`
	Description  string `json:"description"`
	DisplayName  string `json:"displayname"`
	ErrorControl string `json:"errorcontrol"`
	PathName     string `json:"pathname"`
	ProcessID    int    `json:"processid"`
	ServiceType  string `json:"servicetype"`
	Started      string `json:"started"`
	StartMode    string `json:"startmode"`
	StartName    string `json:"startname"`
	State        string `json:"state"`
	Status       string `json:"status"`
	SystemName   string `json:"systemname"`
	ItemMain     string `json:"item_main"`
	DateMain     int    `json:"date_main"`
	TypeMain     string `json:"type_main"`
	EtcMain      string `json:"etc_main"`
}

func (n Service) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type Shortcuts struct {
	UUID             string `json:"uuid"`
	Agent            string `json:"agent"`
	AgentIP          string `json:"agentIP"`
	AgentName        string `json:"agentName"`
	ShortcutName     string `json:"shortcutname"`
	LinkPath         string `json:"linkpath"`
	LinkTo           string `json:"linkto"`
	Arguments        string `json:"arguments"`
	Description      string `json:"description"`
	WorkingDirectory string `json:"workingdirectory"`
	IconLocation     string `json:"iconlocation"`
	BrokenShortcut   bool   `json:"brokenshortcut"`
	Hotkey           string `json:"hotkey"`
	ShowCmd          string `json:"showcmd"`
	ModifyTime       int    `json:"modifytime"`
	ItemMain         string `json:"item_main"`
	DateMain         int    `json:"date_main"`
	TypeMain         string `json:"type_main"`
	EtcMain          string `json:"etc_main"`
}

func (n Shortcuts) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type StartRun struct {
	UUID        string `json:"uuid"`
	Agent       string `json:"agent"`
	AgentIP     string `json:"agentIP"`
	AgentName   string `json:"agentName"`
	Name        string `json:"name"`
	Command     string `json:"command"`
	User        string `json:"user"`
	Location    string `json:"location"`
	Caption     string `json:"caption"`
	Description string `json:"description"`
	UserSID     string `json:"usersid"`
	ItemMain    string `json:"item_main"`
	DateMain    int    `json:"date_main"`
	TypeMain    string `json:"type_main"`
	EtcMain     string `json:"etc_main"`
}

func (n StartRun) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type TaskSchedule struct {
	UUID          string `json:"uuid"`
	Agent         string `json:"agent"`
	AgentIP       string `json:"agentIP"`
	AgentName     string `json:"agentName"`
	Name          string `json:"name"`
	Path          string `json:"path"`
	Command       string `json:"command"`
	LastRunTime   int    `json:"lastruntime"`
	NextRunTime   int    `json:"nextruntime"`
	StartBoundary int64  `json:"startboundary"`
	EndBoundary   int    `json:"endboundary"`
	ItemMain      string `json:"item_main"`
	DateMain      int    `json:"date_main"`
	TypeMain      string `json:"type_main"`
	EtcMain       string `json:"etc_main"`
}

func (n TaskSchedule) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type USBdevices struct {
	UUID                           string `json:"uuid"`
	Agent                          string `json:"agent"`
	AgentIP                        string `json:"agentIP"`
	AgentName                      string `json:"agentName"`
	DeviceInstanceID               string `json:"device_instance_id"`
	DeviceDescription              string `json:"device_description"`
	HardwareIDs                    string `json:"hardware_ids"`
	BusReportedDeviceDescription   string `json:"bus_reported_device_description"`
	DeviceManufacturer             string `json:"device_manufacturer"`
	DeviceFriendlyName             string `json:"device_friendly_name"`
	DeviceLocationInfo             string `json:"device_location_info"`
	DeviceSecurityDescriptorString string `json:"device_security_descriptor_string"`
	ContainerID                    string `json:"containerid"`
	DeviceDisplayCategory          string `json:"device_display_category"`
	DeviceLetter                   string `json:"device_letter"`
	Enumerator                     string `json:"enumerator"`
	InstallDate                    int    `json:"install_date"`
	FirstInstallDate               int    `json:"first_install_date"`
	LastArrivalDate                int    `json:"last_arrival_date"`
	LastRemovalDate                int    `json:"last_removal_date"`
	ItemMain                       string `json:"item_main"`
	DateMain                       int    `json:"date_main"`
	TypeMain                       string `json:"type_main"`
	EtcMain                        string `json:"etc_main"`
}

func (n USBdevices) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type UserAssist struct {
	UUID             string `json:"uuid"`
	Agent            string `json:"agent"`
	AgentIP          string `json:"agentIP"`
	AgentName        string `json:"agentName"`
	Name             string `json:"name"`
	ClassID          string `json:"classid"`
	OfTimesExecuted  int    `json:"of_times_executed"`
	FocusCount       int    `json:"focus_count"`
	FocusTimeSeconds int    `json:"focus_time(s)"`
	ModifiedTime     int    `json:"modifiedtime"`
	ItemMain         string `json:"item_main"`
	DateMain         int    `json:"date_main"`
	TypeMain         string `json:"type_main"`
	EtcMain          string `json:"etc_main"`
}

func (n UserAssist) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type UserProfiles struct {
	UUID               string `json:"uuid"`
	Agent              string `json:"agent"`
	AgentIP            string `json:"agentIP"`
	AgentName          string `json:"agentName"`
	Username           string `json:"username"`
	ProfilePath        string `json:"profilepath"`
	UserSID            string `json:"usersid"`
	RegistryLoaded     bool   `json:"registryloaded"`
	FolderCreatedTime  int    `json:"foldercreatedtime"`
	FolderModifiedTime int    `json:"foldermodifiedtime"`
	LastLoginTime      int    `json:"lastlogontime"`
	PrivilegeLevel     string `json:"privileglevel"`
	ItemMain           string `json:"item_main"`
	DateMain           int    `json:"date_main"`
	TypeMain           string `json:"type_main"`
	EtcMain            string `json:"etc_main"`
}

func (n UserProfiles) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type WindowsActivity struct {
	UUID                 string `json:"uuid"`
	Agent                string `json:"agent"`
	AgentIP              string `json:"agentIP"`
	AgentName            string `json:"agentName"`
	UserName             string `json:"user_name"`
	AppID                string `json:"app_id"`
	AppActivityID        string `json:"app_activity_id"`
	ActivityType         string `json:"activity_type"`
	ActivityStatus       string `json:"activity_status"`
	Tag                  string `json:"tag"`
	Group                string `json:"group"`
	Priority             int    `json:"priority"`
	IsLocalOnly          bool   `json:"is_local_only"`
	ETag                 int    `json:"etag"`
	CreatedInCloud       int    `json:"created_in_cloud"`
	LastModifiedTime     int    `json:"last_modified_time"`
	ExpirationTime       int    `json:"expiration_time"`
	StartTime            int    `json:"start_time"`
	EndTime              int    `json:"end_time"`
	LastModifiedOnClient int    `json:"last_modified_on_client"`
	ItemMain             string `json:"item_main"`
	DateMain             int    `json:"date_main"`
	TypeMain             string `json:"type_main"`
	EtcMain              string `json:"etc_main"`
}

func (n WindowsActivity) Elastical() ([]byte, error) {
	return json.Marshal(n)
}
