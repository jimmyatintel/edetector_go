package work

import (
	"edetector_go/config"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/memory"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	elasticquery "edetector_go/pkg/elastic/query"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"net"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func GiveDetectNetwork(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveDetectNetwork: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	go detectNetworkElastic(p)
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func detectNetworkElastic(p packet.Packet) {
	ip, name := query.GetMachineIPandName(p.GetRkey())
	networkSet := make(map[string]struct{})
	lines := strings.Split(p.GetMessage(), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		uuid := uuid.NewString()
		values := strings.Split(line, "|")
		key := values[0] + "," + values[3]
		networkSet[key] = struct{}{}
		err := elasticquery.SendToDetailsElastic(config.Viper.GetString("ELASTIC_PREFIX")+"_memory_network_detect", &(memory.MemoryNetworkDetect{}), values, uuid, p.GetRkey(), ip, name, "0", "0", "0", "0", "ed_mid")
		if err != nil {
			logger.Error("Error sending to details elastic: ", zap.Any("error", err.Error()))
		}
	}
	elasticquery.UpdateNetworkInfo(p.GetRkey(), networkSet)
}
