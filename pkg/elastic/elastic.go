package elastic

// import (
// 	"log"

// 	"github.com/elastic/go-elasticsearch/v8"
// 	"github.com/elastic/go-elasticsearch/v8/esapi"
// )
// var es *elasticsearch.Client

// func setElkClient() {
//     var err error
// 	cfg := elasticsearch.Config{
// 		Addresses: []string{config.Viper.GetString("ELASTIC_HOST") + ":" + config.Viper.GetString("ELASTIC_PORT")"},
// 	}
// 	es, err = elasticsearch.NewClient(cfg)
// 	if err != nil {
// 		panic(err) // 連線失敗
// 	}
// }

// func main() {
// 	setElkClient()
// 	fmt.Println(es.Info())
// }