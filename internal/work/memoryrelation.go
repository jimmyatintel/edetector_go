package work

import (
	"edetector_go/pkg/logger"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type Relation struct {
	UUID  string
	Name  string
	Child []string
}

func handleRelation(data []byte, agent string) error {
	// ip, name, err := mariadbquery.GetMachineIPandName(agent)
	// if err != nil {
	// 	return err
	// }
	// taskID := mariadbquery.Load_task_id(agent, "StartMemoryTree", 2)
	UUIDMap := make(map[string]int)
	RelationMap := make(map[int](Relation))
	// rootInd := -1
	strData := strings.ReplaceAll(string(data), "\r", "")
	lines := strings.Split(strData, "\n")
	// record the relation
	for _, line := range lines {
		values := strings.Split(line, "|")
		if len(values) != 11 {
			if len(values) != 1 {
				logger.Error("Invalid line: " + line)
			}
			continue
		}
		parent, child, err := getRelation(values)
		if err != nil {
			return err
		}
		generateUUID(agent, parent, &UUIDMap, &RelationMap)
		generateUUID(agent, child, &UUIDMap, &RelationMap)
		// record name
		tmp := RelationMap[child]
		tmp.Name = values[2]
		RelationMap[child] = tmp
		// record relation
		if parent == -1 {
			// rootInd = child
		} else {
			tmp := RelationMap[parent]
			tmp.Child = append(tmp.Child, RelationMap[child].UUID)
			RelationMap[parent] = tmp
		}
	}
	logger.Info("Record the relation:" + agent)
	// send to elastic (details & relation)
	for _, line := range lines {
		values := strings.Split(line, "|")
		if len(values) != 11 {
			if len(values) != 1 {
				logger.Error("Invalid line: " + line)
			}
			continue
		}
		// child, err := strconv.Atoi(strings.TrimSpace(values[0]))
		// if err != nil {
		// 	return err
		// }
		// // send to elastic
		// err = rabbitmq.ToRabbitMQ_Details(config.Viper.GetString("ELASTIC_PREFIX")+"_memory_tree", &MemoryTree{}, values, RelationMap[child].UUID, agent, ip, name, values[2], values[3], "memory", "", "ed_mid", "StartMemoryTree", taskID)
		// if err != nil {
		// 	return err
		// }
	}
	return nil
}

func getRelation(values []string) (int, int, error) {
	parent, err := strconv.Atoi(strings.TrimSpace(values[1]))
	if err != nil {
		return -1, -1, err
	}
	child, err := strconv.Atoi(strings.TrimSpace(values[0]))
	if err != nil {
		return -1, -1, err
	}
	return parent, child, nil
}

func generateUUID(agent string, ind int, UUIDMap *map[string]int, RelationMap *map[int](Relation)) {
	_, exists := (*RelationMap)[ind]
	if !exists {
		uuid := uuid.NewString()
		relation := Relation{
			UUID:  uuid,
			Name:  "",
			Child: []string{},
		}
		(*RelationMap)[ind] = relation
		(*UUIDMap)[uuid] = ind
	}
}

// func BuildMemoryRelation(agent string, field string, value string, parent string, child string) {
// 	index := config.Viper.GetString("ELASTIC_PREFIX") + "_memory_relation"
// 	searchQuery := fmt.Sprintf(`{
// 			"query": {
// 				"bool": {
// 				  "must": [
// 					{ "term": { "agent": "%s" } },
// 					{ "term": { "parent": "%s" } }
// 				  ]
// 				}
// 			  }
// 			}`, agent, value)
// 	hitsArray := elastic.SearchRequest(index, searchQuery, "uuid")
// 	if len(hitsArray) == 0 { // not exists
// 		data := MemoryRelation{}
// 		if field == "parent" {
// 			data = MemoryRelation{
// 				Agent:  agent,
// 				IsRoot: true,
// 				Parent: value,
// 				Child:  []string{child},
// 			}
// 		} else if field == "child" {
// 			data = MemoryRelation{
// 				Agent:  agent,
// 				IsRoot: false,
// 				Parent: value,
// 				Child:  []string{},
// 			}
// 		}
// 		err := rabbitmq.ToRabbitMQ_Relation("_memory_relation", data, "ed_high")
// 		if err != nil {
// 			logger.Error("Error sending to rabbitMQ (relation): " + err.Error())
// 			return
// 		}
// 	} else if len(hitsArray) == 1 { // exists
// 		hitMap, ok := hitsArray[0].(map[string]interface{})
// 		if !ok {
// 			logger.Error("Hit is not a map")
// 			return
// 		}
// 		docID, ok := hitMap["_id"].(string)
// 		if !ok {
// 			logger.Error("docID not found")
// 			return
// 		}
// 		if field == "parent" {
// 			source, ok := hitMap["_source"].(map[string]interface{})
// 			if !ok {
// 				logger.Error("source not found")
// 				return
// 			}
// 			oldChild, ok := source["child"].([]interface{})
// 			if !ok {
// 				logger.Error("children not found")
// 			}
// 			for _, c := range oldChild {
// 				if c.(string) == child {
// 					return
// 				}
// 			}
// 			script := fmt.Sprintf(`
// 			{
// 				"script": {
// 					"source": "ctx._source.child.add(params.value)",
// 					"lang": "painless",
// 					"params": {
// 						"value": "%s"
// 					}
// 				}
// 			}`, child)
// 			err := elastic.UpdateByDocIDRequest(index, docID, script)
// 			if err != nil {
// 				logger.Error("Error updating parent: " + err.Error())
// 				return
// 			}
// 		} else if field == "child" {
// 			script := `
// 			{
// 				"script": {
// 					"source": "ctx._source.isRoot = params.value",
// 					"lang": "painless",
// 					"params": {
// 						"value": "false"
// 					}
// 				}
// 			}`
// 			err := elastic.UpdateByDocIDRequest(index, docID, script)
// 			if err != nil {
// 				logger.Error("Error updating child: " + err.Error())
// 				return
// 			}
// 		}
// 	} else {
// 		logger.Error("More than one relation found: " + field + " parent " + parent + " child " + child)
// 	}
// }
