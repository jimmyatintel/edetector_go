package work

import (
	"edetector_go/config"
	"edetector_go/pkg/logger"
	mariadbquery "edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/rabbitmq"
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
	ip, name, err := mariadbquery.GetMachineIPandName(agent)
	if err != nil {
		return err
	}
	taskID := mariadbquery.Load_task_id(agent, "StartMemoryTree", 2)
	UUIDMap := make(map[string]int)
	RelationMap := make(map[int](Relation))
	rootInd := -1
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
			rootInd = child
		} else {
			tmp := RelationMap[parent]
			tmp.Child = append(tmp.Child, RelationMap[child].UUID)
			RelationMap[parent] = tmp
		}
	}
	logger.Info("Record the relation: " + agent)
	headData := MemoryRelation{}
	// send to elastic (details & relation)
	for _, line := range lines {
		values := strings.Split(line, "|")
		if len(values) != 11 {
			if len(values) != 1 {
				logger.Error("Invalid line: " + line)
			}
			continue
		}
		child, err := strconv.Atoi(strings.TrimSpace(values[0]))
		if err != nil {
			return err
		}
		// send details
		err = rabbitmq.ToRabbitMQ_Details(config.Viper.GetString("ELASTIC_PREFIX")+"_memory_tree", &MemoryTree{}, values, RelationMap[child].UUID, agent, ip, name, values[2], values[3], "memory", "", "ed_mid", "StartMemoryTree", taskID)
		if err != nil {
			return err
		}
		// send relation
		data := MemoryRelation{
			Agent:   agent,
			IsRoot:  false,
			Parent:  RelationMap[child].UUID,
			Child:   RelationMap[child].Child,
			Task_id: taskID,
		}
		if child == rootInd {
			headData = data
		} else {
			err := rabbitmq.ToRabbitMQ_Relation("_explorer_relation", data, "ed_mid")
			if err != nil {
				return err
			}
		}
	}
	logger.Info("Send to elastic (details & relation): " + agent)
	// send head relation
	headData.IsRoot = true
	err = rabbitmq.ToRabbitMQ_Relation("_explorer_relation", headData, "ed_mid")
	if err != nil {
		return err
	}
	logger.Info("Send to elastic (head relation): " + agent)
	// send finish signal
	err = rabbitmq.ToRabbitMQ_FinishSignal(agent, "StartMemoryTree", "ed_mid")
	if err != nil {
		return err
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
