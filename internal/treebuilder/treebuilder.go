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

	"github.com/google/uuid"
)

var fileUnstagePath = "fileUnstage"
var fileStagedPath = "fileStaged"
var limit int
var count int
var cancelMap = map[string][]context.CancelFunc{}

type Relation struct {
	UUID   string
	Name   string
	Path   string
	IsRoot bool
	Child  []string
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
			if len(values) != 10 {
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
			// record name
			tmp := RelationMap[child]
			tmp.Name = values[0]
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
			values := strings.Split(line, "|")
			if len(values) != 10 {
				if len(values) != 1 {
					logger.Error("Invalid line (" + agent + "-" + diskInfo + "): " + line)
				}
				continue
			}
			child, err := strconv.Atoi(values[8])
			if err != nil {
				logger.Error("Error getting child (" + agent + "-" + diskInfo + "): " + err.Error())
				mariadbquery.Failed_task(agent, "StartGetDrive", 6)
				clearBuilder(agent, diskInfo, explorerFile)
				return
			}
			md5_sig := ""
			if fileSystem != "NTFS" {
				md5_sig = values[6]
				values[6] = "0"
			}
			data := Collect_Explorer{
				Explorer: Explorer{
					FileName:          values[0],
					IsDeleted:         values[1] == "1",
					IsDirectory:       values[2] == "2",
					CreateTime:        strToInt(values[3]),
					WriteTime:         strToInt(values[4]),
					AccessTime:        strToInt(values[5]),
					EntryModifiedTime: strToInt(values[6]),
					Datalen:           int64(strToInt(values[7])),
					Path:              RelationMap[child].Path,
					Disk:              diskInfo,
					MD5_Sig:           md5_sig,
					YaraRuleHitCount:  0,
					YaraRuleHit:       "",
					IsRoot:            RelationMap[child].IsRoot,
					Child:             RelationMap[child].Child,
				},
				UUID:      RelationMap[child].UUID,
				Agent:     agent,
				AgentIP:   ip,
				AgentName: name,
				ItemMain:  values[0],
				DateMain:  strToInt(values[3]),
				TypeMain:  "file_table",
				EtcMain:   RelationMap[child].Path,
				Task_id:   taskID,
				Category:  "explorer",
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
	if redis.RedisGetInt(agent+"-DriveUnfinished") == 0 { // last drive -> send finish signal
		err = rabbitmq.ToRabbitMQ_FinishSignal(agent, "StartGetDrive", "ed_low_explorer")
		if err != nil {
			logger.Error("Error sending finish signal to rabbitMQ (" + agent + "): " + err.Error())
			mariadbquery.Failed_task(agent, "StartGetDrive", 6)
			return
		}
	}
	logger.Info("Tree builder task finished: " + agent + "-" + diskInfo)
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

func generateUUID(agent string, ind int, UUIDMap *map[string]int, RelationMap *map[int](Relation)) {
	_, exists := (*RelationMap)[ind]
	if !exists {
		uuid := uuid.NewString()
		relation := Relation{
			UUID:   uuid,
			Name:   "",
			Path:   "",
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

func treeTraversal(agent string, ind int, isRoot bool, path string, diskInfo string, UUIDMap *map[string]int, RelationMap *map[int](Relation), taskID string) {
	disk := strings.Split(diskInfo, "|")[0]
	relation := (*RelationMap)[ind]
	if disk == "Linux" {
		if !isRoot {
			path = path + "/" + relation.Name
		}
	} else {
		if path == "" {
			path = disk + ":"
		} else {
			path = path + "\\" + relation.Name
		}
	}
	if disk == "Linux" && isRoot {
		relation.Path = "/"
	} else {
		relation.Path = path
	}
	(*RelationMap)[ind] = relation
	for _, uuid := range relation.Child {
		treeTraversal(agent, (*UUIDMap)[uuid], false, path, diskInfo, UUIDMap, RelationMap, taskID)
	}
}

func clearBuilder(agent string, disk string, explorerFile string) {
	count--
	cancelMap[agent] = []context.CancelFunc{}
	err := file.MoveFile(explorerFile, filepath.Join(fileStagedPath, agent+"."+disk+".txt"))
	if err != nil {
		logger.Error("Error moving file (" + agent + "-" + disk + "): " + err.Error())
	}
}
