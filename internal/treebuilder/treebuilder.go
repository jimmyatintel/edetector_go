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
	"path/filepath"
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
		explorerFile, agent, disk := file.GetOldestFile(fileUnstagePath, ".txt")
		ip, name, err := query.GetMachineIPandName(agent)
		if err != nil {
			logger.Error("Error getting machine ip and name: " + err.Error())
			query.Failed_task(agent, "StartGetDrive")
			clearBuilder(agent, disk, explorerFile)
			continue
		}
		time.Sleep(3 * time.Second) // wait for fully copy
		explorerContent, err := os.ReadFile(explorerFile)
		if err != nil {
			logger.Error("Read file error: " + err.Error())
			query.Failed_task(agent, "StartGetDrive")
			clearBuilder(agent, disk, explorerFile)
			continue
		}
		logger.Info("Open txt file: " + explorerFile)
		// record the relation
		RelationMap[agent] = make(map[int](Relation))
		rootInd := 0
		lines := strings.Split(string(explorerContent), "\n")
		if terminateDrive(agent, disk, explorerFile) {
			continue
		}
		for _, line := range lines {
			values := strings.Split(line, "|")
			if len(values) != 10 {
				if len(values) != 1 {
					logger.Warn("Invalid line: " + line)
				}
				continue
			}
			parent, child, err := getRelation(values)
			if err != nil {
				logger.Error("Error getting relation: " + err.Error())
				query.Failed_task(agent, "StartGetDrive")
				break
			}
			generateUUID(agent, parent)
			generateUUID(agent, child)
			// record name
			tmp := RelationMap[agent][child]
			tmp.Name = values[0]
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
		if terminateDrive(agent, disk, explorerFile) {
			continue
		}
		// tree traversal & send to elastic(relation)
		treeTraversal(agent, rootInd, true, "", disk)
		logger.Info("Tree traversal & send to elastic (relation)")
		if terminateDrive(agent, disk, explorerFile) {
			continue
		}
		// send to elastic (main & details)
		for _, line := range lines {
			values := strings.Split(line, "|")
			if len(values) != 10 {
				if len(values) != 1 {
					logger.Warn("Invalid line: " + line)
				}
				break
			}
			child, err := strconv.Atoi(values[8])
			if err != nil {
				logger.Error("Error getting child: " + err.Error())
				query.Failed_task(agent, "StartGetDrive")
				break
			}
			values = values[:len(values)-2]
			if values[2] == "2" {
				values[2] = "1"
			}
			err = rabbitmq.ToRabbitMQ_Main(config.Viper.GetString("ELASTIC_PREFIX")+"_explorer", RelationMap[agent][child].UUID, agent, ip, name, values[0], values[3], "file_table", RelationMap[agent][child].Path, "ed_low")
			if err != nil {
				logger.Error("Error sending to rabbitMQ (main): " + err.Error())
				query.Failed_task(agent, "StartGetDrive")
				break
			}
			err = rabbitmq.ToRabbitMQ_Details(config.Viper.GetString("ELASTIC_PREFIX")+"_explorer", &ExplorerDetails{}, values, RelationMap[agent][child].UUID, agent, ip, name, values[0], values[3], "file_table", RelationMap[agent][child].Path, "ed_low")
			if err != nil {
				logger.Error("Error sending to rabbitMQ (details): " + err.Error())
				query.Failed_task(agent, "StartGetDrive")
				break
			}
			time.Sleep(1 * time.Microsecond)
		}
		logger.Info("Send to elastic (main & details)")
		if terminateDrive(agent, disk, explorerFile) {
			continue
		}
		clearBuilder(agent, disk, explorerFile)
		query.Finish_task(agent, "StartGetDrive")
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

func treeTraversal(agent string, ind int, isRoot bool, path string, disk string) {
	relation := RelationMap[agent][ind]
	if disk == "Linux" {
		if path == "" {
			path = relation.Name
		} else {
			path = path + "/" + relation.Name
		}
	} else {
		if path == "" {
			path = disk + ":"
		} else {
			path = path + "\\" + relation.Name
		}
	}

	relation.Path = path
	RelationMap[agent][ind] = relation
	data := ExplorerRelation{
		Agent:  agent,
		IsRoot: isRoot,
		Parent: relation.UUID,
		Child:  relation.Child,
	}
	err := rabbitmq.ToRabbitMQ_Relation("_explorer_relation", data, "ed_low")
	if err != nil {
		logger.Error("Error sending to rabbitMQ (relation): " + err.Error())
		query.Failed_task(agent, "StartGetDrive")
		return
	}
	for _, uuid := range relation.Child {
		treeTraversal(agent, UUIDMap[uuid], false, path, disk)
	}
}

func terminateDrive(agent string, disk string, explorerFile string) bool {
	var flag = false
	if redis.RedisExists(agent+"-terminateDrive") && redis.RedisGetInt(agent+"-terminateDrive") == 1 {
		flag = true
		elastic.DeleteByQueryRequest("agent", agent, "StartGetDrive")
		redis.RedisSet(agent+"-terminateDrive", 0)
		clearBuilder(agent, disk, explorerFile)
	}
	return flag
}

func clearBuilder(agent string, disk string, explorerFile string) {
	UUIDMap = make(map[string]int)
	RelationMap[agent] = nil
	err := file.MoveFile(explorerFile, filepath.Join(fileStagedPath, agent+"."+disk+".txt"))
	if err != nil {
		logger.Error("Error moving file: " + err.Error())
	}
}
