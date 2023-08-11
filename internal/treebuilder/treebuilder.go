package treebuilder

import (
	"edetector_go/config"
	"edetector_go/internal/fflag"
	elasticquery "edetector_go/pkg/elastic/query"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
	"edetector_go/pkg/rabbitmq"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var RelationMap = make(map[string](map[int](Relation)))
var DetailsMap = make(map[string](string))
var Finished = make(chan string)

type Relation struct {
	UUID  string
	Child []string
}

func init() {
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
	if err := mariadb.Connect_init(); err != nil {
		logger.Error("Error connecting to mariadb: " + err.Error())
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("rabbit_enable"); enable && err == nil {
		rabbitmq.Rabbit_init()
		logger.Info("rabbit is enabled.")
	}
}

func Main() {
	logger.Info("Starting tree builder...")
	for {
		var rootInd int
		agent := <-Finished
		fmt_content := DetailsMap[agent]
		RelationMap[agent] = make(map[int](Relation))
		logger.Info("Handling explorer of agent: ", zap.Any("message", agent))

		// send to elastic(main & details) & record the relation
		lines := strings.Split(string(fmt_content), "\n")
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
			if parent == child {
				rootInd = parent
			} else {
				relation := RelationMap[agent][parent]
				relation.Child = append(relation.Child, RelationMap[agent][child].UUID)
				RelationMap[agent][parent] = relation
			}
			//! tmp version: new explorer struct
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
			err = elasticquery.SendToMainElastic(RelationMap[agent][child].UUID, config.Viper.GetString("ELASTIC_PREFIX")+"explorer", agent, values[0], c_time, "file_table", "path(todo)", "ed_low")
			if err != nil {
				logger.Error("Error sending to main elastic: ", zap.Any("error", err.Error()))
				continue
			}
			err = elasticquery.SendToDetailsElastic(RelationMap[agent][child].UUID, config.Viper.GetString("ELASTIC_PREFIX")+"explorer", agent, line, &ExplorerDetails{}, "ed_low")
			if err != nil {
				logger.Error("Error sending to details elastic: ", zap.Any("error", err.Error()))
				continue
			}
		}
		logger.Info("send to elastic(main & details) & record the relation")
		// send to elastic(relation)
		for id, relation := range RelationMap[agent] {
			var isRoot bool
			if rootInd == id {
				isRoot = true
			} else {
				isRoot = false
			}
			data := ExplorerRelation{
				Agent:  agent,
				IsRoot: isRoot,
				Parent: relation.UUID,
				Child:  relation.Child,
			}
			err := elasticquery.SendToRelationElastic(data, "ed_low")
			if err != nil {
				logger.Error("Error sending to relation elastic: ", zap.Any("error", err.Error()))
				continue
			}
		}
		logger.Info("send to elastic(relation)")
		// clear
		RelationMap[agent] = nil
		DetailsMap[agent] = ""
		logger.Info("Finish handling explorer of agent: ", zap.Any("message", agent))
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
		relation := Relation{
			UUID:  uuid.NewString(),
			Child: []string{},
		}
		RelationMap[agent][ind] = relation
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
