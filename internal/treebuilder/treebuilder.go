package treebuilder

import (
	"edetector_go/config"
	"edetector_go/internal/fflag"
	"edetector_go/internal/work"
	elasticquery "edetector_go/pkg/elastic/query"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
	"edetector_go/pkg/rabbitmq"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ExplorerRelation struct {
	Agent  string   `json:"agent"`
	IsRoot bool     `json:"isRoot"`
	Parent string   `json:"parent"`
	Child  []string `json:"child"`
}

func (n ExplorerRelation) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type ExplorerDetails struct {
	UUID              string `json:"uuid"`
	Agent             string `json:"agent"`
	AgentIP           string `json:"agentIP"`
	AgentName         string `json:"agentName"`
	FileName          string `json:"fileName"`
	IsDeleted         bool   `json:"isDeleted"`
	IsDirectory       bool   `json:"isDirectory"`
	CreateTime        int    `json:"createTime"`
	WriteTime         int    `json:"writeTime"`
	AccessTime        int    `json:"accessTime"`
	EntryModifiedTime int    `json:"entryModifiedTime"`
	Datalen           int    `json:"dataLen"`
}

func (n ExplorerDetails) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

var ParentMap = make(map[string]([]Relation))

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
		agent := <-work.Finished
		fmt_content := work.DetailsMap[agent]
		logger.Info("Handling explorer of agent: ", zap.Any("message", agent))
		// init the uuid of explorer
		var relations []Relation
		for i := 0; i < work.ExplorerTotalMap[agent]; i++ {
			relations = append(relations, Relation{
				UUID:  uuid.NewString(),
				Child: []string{},
			})
		}
		ParentMap[agent] = relations

		// send to elastic(main & details) & record the relation
		lines := strings.Split(fmt_content, "\n")
		for i, line := range lines {
			original := strings.Split(line, "|")
			if len(original) < 2 { //! tmp version
				continue
			}
			parent, err := strconv.Atoi(original[2])
			if err != nil {
				logger.Error("Error converting parent to int")
				continue
			}
			child, err := strconv.Atoi(original[0])
			if err != nil {
				logger.Error("Error converting parent to int")
				continue
			}
			var uuid string
			if child <= work.ExplorerTotalMap[agent] {
				uuid = ParentMap[agent][child].UUID
			} else {
				logger.Error("Short length")
			}
			logger.Debug("uuid", zap.Any("message", fmt.Sprintf("%d: %s", i, uuid)))
			if parent == child {
				rootInd = i
			} else {
				ParentMap[agent][parent].Child = append(ParentMap[agent][parent].Child, uuid)
			}
			//! tmp version: new explorer struct
			layout := "2006/01/02 15:04:05"
			t, err := time.Parse(layout, original[5])
			if err != nil {
				logger.Error("Error parsing time", zap.Any("error", err))
			}
			create_time := t.Unix()
			t, err = time.Parse(layout, original[6])
			if err != nil {
				logger.Error("Error parsing time", zap.Any("error", err))
			}
			write_time := t.Unix()
			t, err = time.Parse(layout, original[7])
			if err != nil {
				logger.Error("Error parsing time", zap.Any("error", err))
			}
			access_time := t.Unix()
			t, err = time.Parse(layout, original[8])
			if err != nil {
				logger.Error("Error parsing time", zap.Any("error", err))
			}
			entry_modified_time := t.Unix()

			line = original[1] + "@|@" + original[3] + "@|@" + original[4] + "@|@" + strconv.FormatInt(create_time, 10) + "@|@" + strconv.FormatInt(write_time, 10) + "@|@" + strconv.FormatInt(access_time, 10) + "@|@" + strconv.FormatInt(entry_modified_time, 10) + "@|@" + original[9]
			values := strings.Split(line, "@|@")

			err = elasticquery.SendToMainElastic(uuid, "ed_explorer", agent, values[0], int(create_time), "file_table", "path(todo)", "ed_low")
			if err != nil {
				logger.Error("Error sending to main elastic: ", zap.Any("error", err.Error()))
				continue
			}
			err = elasticquery.SendToDetailsElastic(uuid, "ed_explorer", agent, line, &ExplorerDetails{}, "ed_low")
			if err != nil {
				logger.Error("Error sending to details elastic: ", zap.Any("error", err.Error()))
				continue
			}
		}
		fmt.Println("send to elastic(main & details) & record the relation")
		// send to elastic(relation)
		for i, relation := range ParentMap[agent] {
			var isRoot bool
			if rootInd == i {
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
			}
		}
		fmt.Println("send to elastic(relation)")
		// clear
		ParentMap[agent] = nil
		work.DetailsMap[agent] = ""
	}
}
