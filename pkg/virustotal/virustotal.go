package virustotal

import (
	"edetector_go/pkg/logger"
	mq "edetector_go/pkg/mariadb/query"
	rq "edetector_go/pkg/redis/query"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func CheckAPIKey(apikey string) (bool, error) {
	url := "https://www.virustotal.com/api/v3/ip_addresses/8.8.8.8"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("x-apikey", apikey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		msg := fmt.Sprintf("VirusTotal API key check returned with status code %d", res.StatusCode)
		logger.Warn(msg)
		return false, nil
	}
	return true, nil
}

func APIKeyExist() (bool, error) {
	key, err := mq.LoadVTKey()
	if err != nil {
		return false, err
	}
	if key == "null" {
		return false, nil
	}
	return true, nil
}

func ScanIP(ip string) (int, int, error) {
	// check if api key exists
	apikey, err := mq.LoadVTKey()
	if err != nil {
		return 0, 0, err
	}
	if apikey == "null" {
		return -1, -1, nil
	}
	// check cache
	maliciousCache, totalCache, err := loadCache(ip)
	if err != nil {
		return 0, 0, err
	}
	if maliciousCache != -1 && totalCache != -1 {
		return maliciousCache, totalCache, nil
	}
	// scan ip
	url := "https://www.virustotal.com/api/v3/ip_addresses/" + ip
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("x-apikey", apikey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, 0, err
	}
	if res.StatusCode != 200 {
		return 0, 0, fmt.Errorf("returned with status code %d", res.StatusCode)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	// fmt.Println(string(body))

	var report map[string]interface{}
	err = json.Unmarshal([]byte(body), &report)
	if err != nil {
		return 0, 0, err
	}
	stats, ok := report["data"].(map[string]interface{})["attributes"].(map[string]interface{})["last_analysis_stats"].(map[string]interface{})
	if !ok {
		return 0, 0, fmt.Errorf("last_analysis_stats not found in response")
	}
	harmless, ok := stats["harmless"].(float64)
	if !ok {
		return 0, 0, fmt.Errorf("harmless not found in response")
	}
	malicious, ok := stats["malicious"].(float64)
	if !ok {
		return 0, 0, fmt.Errorf("malicious not found in response")
	}
	suspicious, ok := stats["suspicious"].(float64)
	if !ok {
		return 0, 0, fmt.Errorf("suspicious not found in response")
	}
	undetected, ok := stats["undetected"].(float64)
	if !ok {
		return 0, 0, fmt.Errorf("undetected not found in response")
	}
	timeout, ok := stats["timeout"].(float64)
	if !ok {
		return 0, 0, fmt.Errorf("timeout not found in response")
	}

	total := int(harmless + malicious + suspicious + undetected + timeout)

	// update cache
	err = updateCache(ip, int(malicious), total)
	if err != nil {
		return 0, 0, err
	}

	return int(malicious), total, nil
}

func loadCache(ip string) (int, int, error) {
	result, err := rq.GetVTCache(ip)
	if err != nil {
		return 0, 0, err
	}
	if result == "null" {
		return -1, -1, nil
	}
	// "2|88|2021-07-07T15:00:00Z"
	var malicious, total, timestamp int
	_, err = fmt.Sscanf(result, "%d|%d|%d", &malicious, &total, &timestamp)
	if err != nil {
		return 0, 0, err
	}
	// check if cache is expired (1 day)
	if timestamp+86400 < int(time.Now().Unix()) {
		return -1, -1, nil
	}
	return malicious, total, nil
}

func updateCache(ip string, malicious int, total int) error {
	timestamp := int(time.Now().Unix())
	result := fmt.Sprintf("%d|%d|%d", malicious, total, timestamp)
	err := rq.SetVTCache(ip, result)
	if err != nil {
		return err
	}
	return nil
}
