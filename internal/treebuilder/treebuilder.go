package treebuilder

import (
	"edetector_go/config"
	"edetector_go/internal/fflag"
	"edetector_go/internal/file"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/rabbitmq"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
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

	fflag.Get_fflag()
	if fflag.FFLAG == nil {
		logger.Error("Error loading feature flag")
		return
	}
	vp, err := config.LoadConfig()
	if vp == nil {
		logger.Error("Error loading config file", zap.Any("error", err.Error()))
		return
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("logger_enable"); enable && err == nil {
		logger.InitLogger(config.Viper.GetString("BUILDER_LOG_FILE"), "treebuilder", "TREEBUILDER")
		logger.Info("logger is enabled please check all out info in log file: ", zap.Any("message", config.Viper.GetString("BUILDER_LOG_FILE")))
	}
	if err := mariadb.Connect_init(); err != nil {
		logger.Error("Error connecting to mariadb: " + err.Error())
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("rabbit_enable"); enable && err == nil {
		rabbitmq.Rabbit_init()
		logger.Info("rabbit is enabled.")
	}
}

func Main(version string) {
	builder_init()
	logger.Info("Welcome to edetector tree builder: ", zap.Any("version", version))
	for {
		explorerFile, agent := file.GetOldestFile(fileUnstagePath, ".txt")
		ip, name := query.GetMachineIPandName(agent)
		time.Sleep(3 * time.Second) // wait for fully copy
		explorerContent, err := os.ReadFile(explorerFile)
		if err != nil {
			logger.Error("Read file error", zap.Any("message", err.Error()))
			continue
		}
		RelationMap[agent] = make(map[int](Relation))
		logger.Info("Open txt file: ", zap.Any("message", explorerFile))
		// record the relation
		var rootInd int
		lines := strings.Split(string(explorerContent), "\n")
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}
			values := strings.Split(line, "|")
			parent, child, err := getRelation(values)
			if err != nil {
				logger.Error("error getting relation: ", zap.Any("message", err))
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
		logger.Info("record the relation")
		// tree traversal & send to elastic(relation)
		treeTraversal(agent, rootInd, true, "")
		logger.Info("tree traversal & send to elastic (relation)")
		// send to elastic (main & details)
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}
			values := strings.Split(line, "|")
			child, err := strconv.Atoi(values[8])
			if err != nil {
				logger.Error("error getting child: ", zap.Any("message", err))
				continue
			}
			values = values[:len(values)-2]
			// err = rabbitmq.ToRabbitMQ_Main(config.Viper.GetString("ELASTIC_PREFIX")+"_explorer", RelationMap[agent][child].UUID, agent, ip, name, values[0], values[3], "file_table", RelationMap[agent][child].Path, "ed_low")
			if err != nil {
				logger.Error("Error sending to main elastic: ", zap.Any("error", err.Error()))
				continue
			}
			err = rabbitmq.ToRabbitMQ_Details(config.Viper.GetString("ELASTIC_PREFIX")+"_explorer", &ExplorerDetails{}, values, RelationMap[agent][child].UUID, agent, ip, name, values[0], values[3], "file_table", RelationMap[agent][child].Path, "ed_low")
			if err != nil {
				logger.Error("Error sending to details elastic: ", zap.Any("error", err.Error()))
				continue
			}
			time.Sleep(1 * time.Microsecond)
		}
		logger.Info("send to elastic (main & details)")
		// clear
		UUIDMap = make(map[string]int)
		RelationMap[agent] = nil
		dstPath := strings.ReplaceAll(explorerFile, fileUnstagePath, fileStagedPath)
		err = file.MoveFile(explorerFile, dstPath)
		if err != nil {
			logger.Error("Error moving file: ", zap.Any("error", err.Error()))
		}
		logger.Info("Task finished: ", zap.Any("message", agent))
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
		// logger.Debug("uuid", zap.Any("message", strconv.Itoa(ind)+": "+uuid))
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
		logger.Error("Error sending to relation elastic: ", zap.Any("error", err.Error()))
	}
	for _, uuid := range relation.Child {
		treeTraversal(agent, UUIDMap[uuid], false, path)
	}
}
