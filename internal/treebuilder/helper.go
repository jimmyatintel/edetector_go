package treebuilder

import (
	"context"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func getEntryModifiedTime(fileSystem string, entryModifiedTime string) int {
	if fileSystem == "NTFS" {
		return strToInt(entryModifiedTime)
	}

	return 0
}

func getMD5Sig(fileSystem string, signature string) string {
	if fileSystem != "NTFS" {
		return signature
	}

	return ""
}

func getStartCluster(fileSystem string, explorerDataRow []string) int {
	if fileSystem == "FAT32" {
		return strToInt(explorerDataRow[10])
	}

	return 0
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
			UUID:    uuid,
			Name:    "",
			Path:    "",
			DataLen: 0,
			IsRoot:  false,
			Child:   []string{},
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
	if disk == "Ubuntu" {
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
	if disk == "Ubuntu" && isRoot {
		relation.Path = "/"
	} else {
		relation.Path = path
	}
	(*RelationMap)[ind] = relation
	for _, uuid := range relation.Child {
		treeTraversal(agent, (*UUIDMap)[uuid], false, path, diskInfo, UUIDMap, RelationMap, taskID)
	}
}

func countFileSize(uuid int, UUIDMap *map[string]int, RelationMap *map[int](Relation)) int64 {
	relation := (*RelationMap)[uuid]

	// return dataLen if the dataLen is already calculated or it is not a directory
	if relation.DataLen > 0 || len(relation.Child) == 0 {
		return relation.DataLen
	}

	var childDataLen int64 = 0

	for _, childUUID := range relation.Child {
		childDataLen += countFileSize((*UUIDMap)[childUUID], UUIDMap, RelationMap)
	}

	relation.DataLen = childDataLen
	(*RelationMap)[uuid] = relation

	return relation.DataLen
}

func clearBuilder(agent string, disk string, explorerFile string) {
	count--
	cancelMap[agent] = []context.CancelFunc{}
	err := file.MoveFile(explorerFile, filepath.Join(fileStagedPath, agent+"."+disk+".txt"))
	if err != nil {
		logger.Error("Error moving file (" + agent + "-" + disk + "): " + err.Error())
	}
}
