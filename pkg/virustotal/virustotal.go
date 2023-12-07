package virustotal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ScanIP(ip string, apikey string) (int, int, error) {
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

	return int(malicious), total, nil
}
