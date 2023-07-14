package fflag

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	growthbook "github.com/growthbook/growthbook-golang"
)

// Features API response
type GrowthBookApiResp struct {
	Features json.RawMessage
	Status   int
}

var GB *growthbook.GrowthBook

func GetFeatureMap() []byte {
	// Fetch features JSON from api
	resp, err := http.Get("https://cdn.growthbook.io/api/features/sdk-Mnq4fT9aYWX1ecP")
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	// Just return the features map from the API response
	apiResp := &GrowthBookApiResp{}
	_ = json.Unmarshal(body, apiResp)
	fmt.Println(apiResp.Features)
	return apiResp.Features
}

func Connect_GB() {
	featureMap := GetFeatureMap()
	features := growthbook.ParseFeatureMap(featureMap)

	context := growthbook.NewContext().
		WithFeatures(features)
	GB = growthbook.New(context)
}
