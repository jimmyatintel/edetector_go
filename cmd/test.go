package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
)

func main() {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://ela-master.ed.qa:9200"},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		fmt.Println("Error creating Elasticsearch client: ", err)
		return
	}

	// Define your update query and body
	updateQuery := `
	{
		"script": {
			"source": "ctx._source.processConnectIP = params.processConnectIP",
			"lang": "painless",
			"params": {
				"processConnectIP": "test"
			}
		},
		"query": {
			"term": {
				"processId": 42284
			}
		}
	}
	`

	// Create an update request
	req := esapi.UpdateByQueryRequest{
		Index: []string{"peggy_memory"},
		Body:  strings.NewReader(updateQuery),
	}

	// Execute the update request
	res, err := req.Do(context.Background(), es)
	if err != nil {
		fmt.Println("Error executing update request: ", err)
		return
	}
	defer res.Body.Close()

	fmt.Println("Update operation completed")
}
