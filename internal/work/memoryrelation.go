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
	UUID   string
	Name   string
	IsRoot bool
	Child  []string
}

func handleRelation(data []byte, agent string) error {
	ip, name, err := mariadbquery.GetMachineIPandName(agent)
	if err != nil {
		return err
	}
	taskID := mariadbquery.Load_task_id(agent, "StartMemoryTree", 2)
	UUIDMap := make(map[string]int)
	RelationMap := make(map[int](Relation))
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
			tmp := RelationMap[child]
			tmp.IsRoot = true
			RelationMap[child] = tmp
		} else {
			tmp := RelationMap[parent]
			tmp.Child = append(tmp.Child, RelationMap[child].UUID)
			RelationMap[parent] = tmp
		}
	}
	logger.Info("Record the relation: " + agent)
	headData := Collect_MemoryTree{}
	// send to elastic
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
		data := Collect_MemoryTree{
			MemoryTree: MemoryTree{
				ProcessId:               strToInt(values[0]),
				ParentProcessId:         strToInt(values[1]),
				ProcessName:             values[2],
				ProcessCreateTime:       strToInt(values[3]),
				ParentProcessName:       values[4],
				ParentProcessCreateTime: strToInt(values[5]),
				ProcessPath:             values[6],
				UserName:                values[7],
				IsPacked:                values[8] == "1",
				DynamicCommand:          values[9],
				IsHide:                  values[10] == "1",
				IsRoot:                  RelationMap[child].IsRoot,
				Child:                   RelationMap[child].Child,
			},
			UUID:      RelationMap[child].UUID,
			Agent:     agent,
			AgentIP:   ip,
			AgentName: name,
			ItemMain:  values[2],
			DateMain:  strToInt(values[3]),
			TypeMain:  "memory",
			EtcMain:   "",
			Task_id:   taskID,
			Category:  "memory_tree",
		}

		if RelationMap[child].IsRoot {
			headData = data
			continue
		}
		// send details
		err = rabbitmq.ToRabbitMQ_Tree(config.Viper.GetString("ELASTIC_PREFIX")+"_memory", data, "ed_mid")
		if err != nil {
			return err
		}
	}
	logger.Info("Send to elastic: " + agent)
	// send head
	err = rabbitmq.ToRabbitMQ_Tree(config.Viper.GetString("ELASTIC_PREFIX")+"_memory", headData, "ed_mid")
	if err != nil {
		return err
	}
	logger.Info("Send to elastic (head): " + agent)
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
			UUID:   uuid,
			Name:   "",
			IsRoot: false,
			Child:  []string{},
		}
		(*RelationMap)[ind] = relation
		(*UUIDMap)[uuid] = ind
	}
}

func strToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return num
}
