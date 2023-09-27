package treebuilder

import (
	"edetector_go/config"
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/rabbitmq"
	"edetector_go/pkg/redis"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var RelationMap = make(map[string](map[int](Relation)))
var UUIDMap = make(map[string]int)
var fileUnstagePath = "fileUnstage"
var fileStagedPath = "fileStaged"

type Relation struct {
	UUID  string
	Name  string
	Path  string
	Child []string
}

func builder_init() {
	file.CheckDir(fileUnstagePath)
	file.CheckDir(fileStagedPath)

	// fflag.Get_fflag()
	// if fflag.FFLAG == nil {
	// 	logger.Panic("Error loading feature flag")
	// 	panic("Error loading feature flag")
	// }
	vp, err := config.LoadConfig()
	if vp == nil {
		logger.Panic("Error loading config file: " + err.Error())
		panic(err)
	}
	if true {
		logger.InitLogger(config.Viper.GetString("BUILDER_LOG_FILE"), "treebuilder", "TREEBUILDER")
		logger.Info("logger is enabled please check all out info in log file: " + config.Viper.GetString("BUILDER_LOG_FILE"))
	}
	connString, err := mariadb.Connect_init()
	if err != nil {
		logger.Panic("Error connecting to mariadb: " + err.Error())
		panic(err)
	} else {
		logger.Info("Mariadb connectionString: " + connString)
	}
	if true {
		if db := redis.Redis_init(); db == nil {
			logger.Panic("Error connecting to redis")
			panic(err)
		}
	}
	if true {
		rabbitmq.Rabbit_init()
		logger.Info("Rabbit is enabled.")
	}
	if true {
		elastic.Elastic_init()
		logger.Info("Elastic is enabled.")
	}
}

func Main(version string) {
	builder_init()
	logger.Info("Welcome to edetector tree builder: " + version)
	for {
		explorerFile, agent := file.GetOldestFile(fileUnstagePath, ".txt")
		ip, name := query.GetMachineIPandName(agent)
		time.Sleep(3 * time.Second) // wait for fully copy
		explorerContent, err := os.ReadFile(explorerFile)
		if err != nil {
			logger.Error("Read file error: " + err.Error())
			continue
		}
		RelationMap[agent] = make(map[int](Relation))
		logger.Info("Open txt file: " + explorerFile)
		// record the relation
		var rootInd int
		lines := strings.Split(string(explorerContent), "\n")
		if terminateDrive(agent, explorerFile) {
			continue
		}
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}
			values := strings.Split(line, "|")
			parent, child, err := getRelation(values)
			if err != nil {
				logger.Error("Error getting relation: " + err.Error())
				continue
			}
			generateUUID(agent, parent)
			generateUUID(agent, child)
			// record name
			tmp := RelationMap[agent][child]
			tmp.Name = values[1]
			RelationMap[agent][child] = tmp
			// record relation
			if parent == child {
				rootInd = parent
			} else {
				tmp := RelationMap[agent][parent]
				tmp.Child = append(tmp.Child, RelationMap[agent][child].UUID)
				RelationMap[agent][parent] = tmp
			}
		}
		logger.Info("Record the relation")
		if terminateDrive(agent, explorerFile) {
			continue
		}
		// tree traversal & send to elastic(relation)
		treeTraversal(agent, rootInd, true, "")
		logger.Info("Tree traversal & send to elastic (relation)")
		if terminateDrive(agent, explorerFile) {
			continue
		}
		// send to elastic (main & details)
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}
			values := strings.Split(line, "|")
			child, err := strconv.Atoi(values[8])
			if err != nil {
				logger.Error("Error getting child: " + err.Error())
				continue
			}
			values = values[:len(values)-2]
			// err = rabbitmq.ToRabbitMQ_Main(config.Viper.GetString("ELASTIC_PREFIX")+"_explorer", RelationMap[agent][child].UUID, agent, ip, name, values[0], values[3], "file_table", RelationMap[agent][child].Path, "ed_low")
			if err != nil {
				logger.Error("Error sending to rabbitMQ (main): " + err.Error())
				continue
			}
			err = rabbitmq.ToRabbitMQ_Details(config.Viper.GetString("ELASTIC_PREFIX")+"_explorer", &ExplorerDetails{}, values, RelationMap[agent][child].UUID, agent, ip, name, values[0], values[3], "file_table", RelationMap[agent][child].Path, "ed_low")
			if err != nil {
				logger.Error("Error sending to rabbitMQ (details): " + err.Error())
				continue
			}
			time.Sleep(1 * time.Microsecond)
		}
		logger.Info("Send to elastic (main & details)")
		if terminateDrive(agent, explorerFile) {
			continue
		}
		clearBuilder(agent, explorerFile)
		logger.Info("Tree builder task finished: " + agent)
	}
}

func getRelation(values []string) (int, int, error) {
	values[9] = strings.TrimSpace(values[9])
	parent, err := strconv.Atoi(values[9])
	if err != nil {
		return -1, -1, err
	}
	child, err := strconv.Atoi(values[8])
	if err != nil {
		return -1, -1, err
	}
	return parent, child, nil
}

func generateUUID(agent string, ind int) {
	_, exists := RelationMap[agent][ind]
	if !exists {
		uuid := uuid.NewString()
		relation := Relation{
			UUID:  uuid,
			Name:  "",
			Path:  "",
			Child: []string{},
		}
		RelationMap[agent][ind] = relation
		UUIDMap[uuid] = ind
	}
}

func treeTraversal(agent string, ind int, isRoot bool, path string) {
	relation := RelationMap[agent][ind]
	path = path + "/" + relation.Name
	relation.Path = path
	RelationMap[agent][ind] = relation
	data := ExplorerRelation{
		Agent:  agent,
		IsRoot: isRoot,
		Parent: relation.UUID,
		Child:  relation.Child,
	}
	err := rabbitmq.ToRabbitMQ_Relation(data, "ed_low")
	if err != nil {
		logger.Error("Error sending to rabbitMQ (relation): " + err.Error())
	}
	for _, uuid := range relation.Child {
		treeTraversal(agent, UUIDMap[uuid], false, path)
	}
}

func terminateDrive(agent string, explorerFile string) bool {
	var flag = false
	if redis.RedisExists(agent+"-terminateFinishIteration") && redis.RedisGetInt(agent+"-terminateFinishIteration") == 0 {
		return flag
	}
	if redis.RedisExists(agent+"-terminateDrive") && redis.RedisGetInt(agent+"-terminateDrive") == 1 {
		flag = true
		elastic.DeleteByQueryRequest("agent", agent, "StartGetDrive")
		redis.RedisSet(agent+"-terminateDrive", 0)
		clearBuilder(agent, explorerFile)
	}
	if redis.RedisExists(agent+"-terminateDrive") && redis.RedisExists(agent+"-terminateCollect") && redis.RedisGetInt(agent+"-terminateDrive") == 0 && redis.RedisGetInt(agent+"-terminateCollect") == 0 {
		query.Finish_task(agent, "Terminate")
	}
	return flag
}

func clearBuilder(agent string, explorerFile string) {
	UUIDMap = make(map[string]int)
	RelationMap[agent] = nil
	dstPath := strings.ReplaceAll(explorerFile, fileUnstagePath, fileStagedPath)
	err := file.MoveFile(explorerFile, dstPath)
	if err != nil {
		logger.Error("Error moving file: " + err.Error())
	}
}
