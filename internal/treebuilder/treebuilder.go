package treebuilder

//TBD
import (
	"context"
	"edetector_go/config"
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
	mariadbquery "edetector_go/pkg/mariadb/query"
	elasticquery "edetector_go/pkg/elastic/query"
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
	if redis.RedisGetString(agent+"-DriveUnfinished") == redis.RedisGetString(agent+"-DriveTotal") {
		elasticquery.DeleteRepeat(agent, "StartGetDrive")
	}
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
	explorerContent, err := os.ReadFile(explorerFile)
	if err != nil {
		logger.Error("Read file error (" + agent + "-" + diskInfo + "): " + err.Error())
		mariadbquery.Failed_task(agent, "StartGetDrive", 6)
		clearBuilder(agent, diskInfo, explorerFile)
		return
	}
	logger.Info("Open txt file: " + explorerFile)
	// record the relation
	rootInd := 0
	lines := strings.Split(string(explorerContent), "\n")
	for _, line := range lines {
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
				rootInd = parent
			} else {
				tmp := RelationMap[parent]
				tmp.Child = append(tmp.Child, RelationMap[child].UUID)
				RelationMap[parent] = tmp
			}
		}
	}
	logger.Info("Record the relation (" + agent + "-" + diskInfo + ")")
	// tree traversal & send to elastic(relation)
	treeTraversal(agent, rootInd, true, "", diskInfo, &UUIDMap, &RelationMap)
	logger.Info("Tree traversal & send relation to elastic (" + agent + "-" + diskInfo + ")")
	// send to elastic (main & details)
	for _, line := range lines {
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
			values = values[:len(values)-2]
			if values[2] == "2" {
				values[2] = "1"
			}
			values = append(values, RelationMap[child].Path)
			values = append(values, diskInfo)
			if fileSystem == "NTFS" {
				values = append(values, "")
			} else {
				values = append(values, values[6])
				values[6] = "0"
			}
			values = append(values, "0", "") // yara rule hit count & yara rule hit
			err = rabbitmq.ToRabbitMQ_Main(config.Viper.GetString("ELASTIC_PREFIX")+"_explorer", RelationMap[child].UUID, agent, ip, name, values[0], values[3], "file_table", RelationMap[child].Path, "ed_low")
			if err != nil {
				logger.Error("Error sending to main rabbitMQ (" + agent + "-" + diskInfo + "): " + err.Error())
				mariadbquery.Failed_task(agent, "StartGetDrive", 6)
				clearBuilder(agent, diskInfo, explorerFile)
				return
			}
			err = rabbitmq.ToRabbitMQ_Details(config.Viper.GetString("ELASTIC_PREFIX")+"_explorer", &ExplorerDetails{}, values, RelationMap[child].UUID, agent, ip, name, values[0], values[3], "file_table", RelationMap[child].Path, "ed_low", "StartGetDrive")
			if err != nil {
				logger.Error("Error sending to details rabbitMQ (" + agent + "-" + diskInfo + "): " + err.Error())
				mariadbquery.Failed_task(agent, "StartGetDrive", 6)
				clearBuilder(agent, diskInfo, explorerFile)
				return
			}
			time.Sleep(1 * time.Microsecond)
		}
	}
	logger.Info("Send main & details to elastic (" + agent + "-" + diskInfo + ")")
	clearBuilder(agent, diskInfo, explorerFile)
	redis.RedisSet_AddInteger(agent+"-DriveUnfinished", -1)
	if redis.RedisGetInt(agent+"-DriveUnfinished") == 0 { // last drive -> send finish signal
		err = rabbitmq.ToRabbitMQ_FinishSignal(agent, "StartGetDrive", "ed_low")
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
			UUID:  uuid,
			Name:  "",
			Path:  "",
			Child: []string{},
		}
		(*RelationMap)[ind] = relation
		(*UUIDMap)[uuid] = ind
	}
}

func treeTraversal(agent string, ind int, isRoot bool, path string, diskInfo string, UUIDMap *map[string]int, RelationMap *map[int](Relation)) {
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
	data := ExplorerRelation{
		Agent:  agent,
		IsRoot: isRoot,
		Parent: relation.UUID,
		Child:  relation.Child,
	}
	err := rabbitmq.ToRabbitMQ_Relation("_explorer_relation", data, "ed_low")
	if err != nil {
		logger.Error("Error sending to relation rabbitMQ (" + agent + "-" + diskInfo + "): " + err.Error())
		mariadbquery.Failed_task(agent, "StartGetDrive", 6)
		clearBuilder(agent, diskInfo, "")
		return
	}
	for _, uuid := range relation.Child {
		treeTraversal(agent, (*UUIDMap)[uuid], false, path, diskInfo, UUIDMap, RelationMap)
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
