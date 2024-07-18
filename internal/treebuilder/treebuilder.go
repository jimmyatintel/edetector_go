package treebuilder

//TBD
import (
	"bufio"
	"context"
	"edetector_go/config"
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
	mariadbquery "edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/rabbitmq"
	"edetector_go/pkg/redis"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var fileUnstagePath = filepath.Join("static", "fileUnstage")
var fileStagedPath = filepath.Join("static", "fileStaged")
var limit int
var count int
var cancelMap = map[string][]context.CancelFunc{}

type Relation struct {
	UUID    string
	Name    string
	Path    string
	DataLen int64
	IsRoot  bool
	Child   []string
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
	limit = config.Viper.GetInt("PARSER_BUILDER_LIMIT")
}

func Main(version string) {
	builder_init()
	logger.Info("Welcome to edetector tree builder: " + version)
	count = 0
	go terminateDrive()
	for {
		if count < limit {
			explorerFile, agent, diskInfo := file.GetOldestFile(fileUnstagePath, ".txt")
			count++
			ctx, cancel := context.WithCancel(context.Background())
			cancelMap[agent] = append(cancelMap[agent], cancel)
			go treeBuilder(ctx, explorerFile, agent, diskInfo)
		}
		time.Sleep(10 * time.Second)
	}
}

func terminateDrive() {
	terminating := 5
	for {
		time.Sleep(10 * time.Second)
		handlingTasks, err := mariadbquery.Load_stored_task("nil", "nil", terminating, "StartGetDrive")
		if err != nil {
			logger.Error("Error loading stored task: " + err.Error())
			continue
		}
		for _, t := range handlingTasks {
			logger.Info("Received terminate drive: " + t[1])
			for i, c := range cancelMap[t[1]] {
				if c != nil {
					cancelMap[t[1]][i]()
				}
			}
			mariadbquery.Terminated_task(t[1], "StartGetDrive", terminating)
		}
	}
}

func treeBuilder(ctx context.Context, explorerFile string, agent string, diskInfo string) {
	parts := strings.Split(diskInfo, "|")
	if len(parts) != 2 {
		logger.Error("Invalid diskInfo: " + diskInfo)
		mariadbquery.Failed_task(agent, "StartGetDrive", 6)
		return
	}
	fileSystem := parts[1]
	time.Sleep(3 * time.Second) // wait for fully copy
	UUIDMap := make(map[string]int)
	RelationMap := make(map[int](Relation))
	ip, name, err := mariadbquery.GetMachineIPandName(agent)
	if err != nil {
		logger.Error("Error getting machine ip and name (" + agent + "-" + diskInfo + "): " + err.Error())
		mariadbquery.Failed_task(agent, "StartGetDrive", 6)
		clearBuilder(agent, diskInfo, explorerFile)
		return
	}
	file1, err := os.Open(explorerFile)
	if err != nil {
		logger.Error("Read file error (" + agent + "-" + diskInfo + "): " + err.Error())
		mariadbquery.Failed_task(agent, "StartGetDrive", 6)
		clearBuilder(agent, diskInfo, explorerFile)
		return
	}
	defer file1.Close()
	logger.Info("Open txt file: " + explorerFile)
	// record the relation
	rootInd := 0
	scanner1 := bufio.NewScanner(file1)
	for scanner1.Scan() {
		line := scanner1.Text()
		select {
		case <-ctx.Done():
			logger.Info("Terminate drive (" + diskInfo + "): " + agent)
			clearBuilder(agent, diskInfo, explorerFile)
			return
		default:
			values := strings.Split(line, "|")
			if len(values) != 10 && len(values) != 11 {
				if len(values) != 1 {
					logger.Error("Invalid line (" + agent + "-" + diskInfo + "): " + line)
				}
				continue
			}
			parent, child, err := getRelation(values)
			if err != nil {
				logger.Error("Error getting relation (" + agent + "-" + diskInfo + "): " + err.Error())
				mariadbquery.Failed_task(agent, "StartGetDrive", 6)
				clearBuilder(agent, diskInfo, explorerFile)
				return
			}
			generateUUID(agent, parent, &UUIDMap, &RelationMap)
			generateUUID(agent, child, &UUIDMap, &RelationMap)
			// record name and dataLen
			tmp := RelationMap[child]
			tmp.Name = values[0]
			tmp.DataLen, err = strconv.ParseInt(values[7], 10, 64)
			if err != nil {
				logger.Error("Error getting dataLen (" + agent + "-" + diskInfo + "): " + err.Error())
				mariadbquery.Failed_task(agent, "StartGetDrive", 6)
				clearBuilder(agent, diskInfo, explorerFile)
				return
			}
			RelationMap[child] = tmp
			// record relation
			if parent == child {
				tmp := RelationMap[parent]
				tmp.IsRoot = true
				RelationMap[parent] = tmp
				rootInd = parent
			} else {
				tmp := RelationMap[parent]
				tmp.Child = append(tmp.Child, RelationMap[child].UUID)
				RelationMap[parent] = tmp
			}
		}
	}
	if err := scanner1.Err(); err != nil {
		logger.Error("Read file error (" + agent + "-" + diskInfo + "): " + err.Error())
		mariadbquery.Failed_task(agent, "StartGetDrive", 6)
		clearBuilder(agent, diskInfo, explorerFile)
		return
	}
	file1.Close()
	logger.Info("Record the relation (" + agent + "-" + diskInfo + ")")

	// tree traversal
	taskID := mariadbquery.Load_task_id(agent, "StartGetDrive", 2)
	treeTraversal(agent, rootInd, true, "", diskInfo, &UUIDMap, &RelationMap, taskID)
	logger.Info("Tree traversal & send relation to elastic (" + agent + "-" + diskInfo + ")")

	// count file size
	countFileSize(rootInd, &UUIDMap, &RelationMap)
	logger.Info("Finsh counting file size (" + agent + "-" + diskInfo + ")")

	// send to elastic
	headData := Collect_Explorer{}
	file2, err := os.Open(explorerFile)
	if err != nil {
		logger.Error("Read file error (" + agent + "-" + diskInfo + "): " + err.Error())
		mariadbquery.Failed_task(agent, "StartGetDrive", 6)
		clearBuilder(agent, diskInfo, explorerFile)
		return
	}
	defer file2.Close()

	scanner2 := bufio.NewScanner(file2)
	for scanner2.Scan() {
		line := scanner2.Text()
		select {
		case <-ctx.Done():
			logger.Info("Terminate drive (" + diskInfo + "): " + agent)
			clearBuilder(agent, diskInfo, explorerFile)
			return
		default:
			explorerDataRaw := strings.Split(line, "|")
			if len(explorerDataRaw) != 10 && len(explorerDataRaw) != 11 {
				if len(explorerDataRaw) != 1 {
					logger.Error("Invalid line (" + agent + "-" + diskInfo + "): " + line)
				}
				continue
			}

			child, err := strconv.Atoi(explorerDataRaw[8])
			if err != nil {
				logger.Error("Error getting child (" + agent + "-" + diskInfo + "): " + err.Error())
				mariadbquery.Failed_task(agent, "StartGetDrive", 6)
				clearBuilder(agent, diskInfo, explorerFile)
				return
			}

			timestamp, err := mariadbquery.GetTaskTimestamp(taskID)
			if err != nil {
				logger.Error("Error getting task timestamp (" + agent + "-" + diskInfo + "): " + err.Error())
				mariadbquery.Failed_task(agent, "StartGetDrive", 6)
				clearBuilder(agent, diskInfo, explorerFile)
				return
			}

			data := Collect_Explorer{
				Explorer: Explorer{
					FileName:          explorerDataRaw[0],
					FileId:            strToInt(explorerDataRaw[8]),
					IsDeleted:         explorerDataRaw[1] == "1",
					IsDirectory:       explorerDataRaw[2] != "0",
					CreateTime:        strToInt(explorerDataRaw[3]),
					WriteTime:         strToInt(explorerDataRaw[4]),
					AccessTime:        strToInt(explorerDataRaw[5]),
					EntryModifiedTime: getEntryModifiedTime(fileSystem, explorerDataRaw[6]),
					Datalen:           RelationMap[child].DataLen,
					Path:              RelationMap[child].Path,
					Disk:              diskInfo,
					MD5_Sig:           getMD5Sig(fileSystem, explorerDataRaw[6]),
					StartCluster:      getStartCluster(fileSystem, explorerDataRaw),
					YaraRuleHitCount:  0,
					YaraRuleHit:       "",
					IsRoot:            RelationMap[child].IsRoot,
					Child:             RelationMap[child].Child,
				},
				UUID:          RelationMap[child].UUID,
				Agent:         agent,
				AgentIP:       ip,
				AgentName:     name,
				ItemMain:      explorerDataRaw[0],
				DateMain:      strToInt(explorerDataRaw[3]),
				TypeMain:      "file_table",
				EtcMain:       RelationMap[child].Path,
				Task_id:       taskID,
				Category:      "explorer",
				TaskTimestamp: strToInt(timestamp),
			}
			if RelationMap[child].IsRoot {
				headData = data
				continue
			}
			err = rabbitmq.ToRabbitMQ_Tree(config.Viper.GetString("ELASTIC_PREFIX")+"_explorer", data, "ed_low_explorer")
			if err != nil {
				logger.Error("Error sending to details rabbitMQ (" + agent + "-" + diskInfo + "): " + err.Error())
				mariadbquery.Failed_task(agent, "StartGetDrive", 6)
				clearBuilder(agent, diskInfo, explorerFile)
				return
			}
		}
		time.Sleep(1 * time.Microsecond)
	}
	file2.Close()
	logger.Info("Send details to elastic (" + agent + "-" + diskInfo + ")")
	// send ExplorerTreeHead in the end
	err = rabbitmq.ToRabbitMQ_Tree(config.Viper.GetString("ELASTIC_PREFIX")+"_explorer", headData, "ed_low_explorer")
	if err != nil {
		logger.Error("Error sending to details rabbitMQ (" + agent + "-" + diskInfo + "): " + err.Error())
		mariadbquery.Failed_task(agent, "StartGetDrive", 6)
		clearBuilder(agent, diskInfo, explorerFile)
		return
	}
	clearBuilder(agent, diskInfo, explorerFile)

	redis.RedisSet_AddInteger(agent+"-DriveUnfinished", -1)
	driveUnfinished, err := redis.RedisGetInt(agent + "-DriveUnfinished")
	if err != nil {
		logger.Error("Error getting drive unfinished (" + agent + "): " + err.Error())
		mariadbquery.Failed_task(agent, "StartGetDrive", 6)
		return
	}

	if driveUnfinished == 0 { // last drive -> send finish signal
		err = rabbitmq.ToRabbitMQ_FinishSignal(agent, "StartGetDrive", "ed_low_explorer")
		if err != nil {
			logger.Error("Error sending finish signal to rabbitMQ (" + agent + "): " + err.Error())
			mariadbquery.Failed_task(agent, "StartGetDrive", 6)
			return
		}
	}
	logger.Info("Tree builder task finished: " + agent + "-" + diskInfo)
}
