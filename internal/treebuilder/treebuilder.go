package treebuilder

import (
	"edetector_go/config"
	"edetector_go/internal/fflag"
	"edetector_go/internal/file"
	elasticquery "edetector_go/pkg/elastic/query"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
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

func init() {
	file.CheckDir(fileUnstagePath)
	file.CheckDir(fileStagedPath)

	fflag.Get_fflag()
	if fflag.FFLAG == nil {
		logger.Error("Error loading feature flag")
		return
	}
	vp := config.LoadConfig()
	if vp == nil {
		logger.Error("Error loading config file")
		return
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("logger_enable"); enable && err == nil {
		logger.InitLogger(config.Viper.GetString("BUILDER_LOG_FILE"))
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

func Main() {
	for {
		explorerFile := file.GetOldestFile(fileUnstagePath, ".txt")
		path := strings.Split(strings.Split(explorerFile, ".txt")[0], "/")
		agent := strings.Split(path[len(path)-1], "-")[0]
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
			original := strings.Split(line, "|")
			parent, child, err := getRelation(original)
			if err != nil {
				logger.Error("error getting relation: ", zap.Any("message", err))
				continue
			}
			generateUUID(agent, parent)
			generateUUID(agent, child)
			// record name
			tmp := RelationMap[agent][child]
			tmp.Name = original[1]
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
			original := strings.Split(line, "|")
			// ! tmp version: new explorer struct
			create_time, write_time, access_time, entry_modified_time, err := tmpGetTime(original)
			if err != nil {
				logger.Error("error parsing time: ", zap.Any("message", err))
				continue
			}
			line = original[1] + "@|@" + original[3] + "@|@" + original[4] + "@|@" + create_time + "@|@" + write_time + "@|@" + access_time + "@|@" + entry_modified_time + "@|@" + original[9]
			values := strings.Split(line, "@|@")
			c_time, err := strconv.Atoi(create_time)
			if err != nil {
				logger.Error("error converting time")
			}
			_, child, err := getRelation(original)
			if err != nil {
				logger.Error("error getting relation: ", zap.Any("message", err))
				continue
			}
			// err = elasticquery.SendToMainElastic(RelationMap[agent][child].UUID, config.Viper.GetString("ELASTIC_PREFIX")+"_explorer", agent, values[0], c_time, "file_table", RelationMap[agent][child].Path, "ed_low")
			// if err != nil {
			// 	logger.Error("Error sending to main elastic: ", zap.Any("error", err.Error()))
			// 	continue
			// }
			err = elasticquery.SendToDetailsElastic(RelationMap[agent][child].UUID, config.Viper.GetString("ELASTIC_PREFIX")+"_explorer", agent, line, &ExplorerDetails{}, "ed_low", values[0], c_time, "file_table", RelationMap[agent][child].Path)
			if err != nil {
				logger.Error("Error sending to details elastic: ", zap.Any("error", err.Error()))
				continue
			}
		}
		logger.Info("send to elastic (main & details)")
		// clear
		UUIDMap = nil
		RelationMap[agent] = nil
		dstPath := strings.ReplaceAll(explorerFile, fileUnstagePath, fileStagedPath)
		err = file.MoveFile(explorerFile, dstPath)
		if err != nil {
			logger.Error("Error moving file: ", zap.Any("error", err.Error()))
		}
		logger.Info("Task finished: ", zap.Any("message", agent))
	}
}

func getRelation(original []string) (int, int, error) {
	parent, err := strconv.Atoi(original[2])
	if err != nil {
		return -1, -1, err
	}
	child, err := strconv.Atoi(original[0])
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
	err := elasticquery.SendToRelationElastic(data, "ed_low")
	if err != nil {
		logger.Error("Error sending to relation elastic: ", zap.Any("error", err.Error()))
	}
	for _, uuid := range relation.Child {
		treeTraversal(agent, UUIDMap[uuid], false, path)
	}
}

func tmpGetTime(original []string) (string, string, string, string, error) { //!tmp version
	layout := "2006/01/02 15:04:05"
	t, err := time.Parse(layout, original[5])
	if err != nil {
		return "", "", "", "", err
	}
	create_time := strconv.FormatInt(t.Unix(), 10)
	t, err = time.Parse(layout, original[6])
	if err != nil {
		return "", "", "", "", err
	}
	write_time := strconv.FormatInt(t.Unix(), 10)
	t, err = time.Parse(layout, original[7])
	if err != nil {
		return "", "", "", "", err
	}
	access_time := strconv.FormatInt(t.Unix(), 10)
	t, err = time.Parse(layout, original[8])
	if err != nil {
		return "", "", "", "", err
	}
	entry_modified_time := strconv.FormatInt(t.Unix(), 10)

	return create_time, write_time, access_time, entry_modified_time, err
}
