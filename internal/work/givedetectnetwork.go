package work

import (
	"edetector_go/config"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/ip2location"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/rabbitmq"
	"net"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func GiveDetectNetwork(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveDetectNetwork: " + p.GetRkey())
	go detectNetworkElastic(p)
	err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func detectNetworkElastic(p packet.Packet) {
	ip, name, err := query.GetMachineIPandName(p.GetRkey())
	if err != nil {
		logger.Error("Error getting machine ip and name: " + err.Error())
		return
	}
	networkSet := make(map[string]struct{})
	lines := strings.Split(p.GetMessage(), "\n")
	for _, line := range lines {
		uuid := uuid.NewString()
		values := parseNetowrk(line, &networkSet, ip)
		if values == nil {
			continue
		}
		err = rabbitmq.ToRabbitMQ_Details(config.Viper.GetString("ELASTIC_PREFIX")+"_memory_network", &(MemoryNetwork{}), values, uuid, p.GetRkey(), ip, name, "0", "0", "0", "0", "ed_mid")
		if err != nil {
			logger.Error("Error sending to rabbitMQ (details): " + err.Error())
		}
	}
	elastic.UpdateNetworkInfo(p.GetRkey(), networkSet)
}

func parseNetowrk(line string, networkSet *map[string]struct{}, ip string) []string{
	values := strings.Split(line, "|")
	if len(values) != 6 {
		if len(values) != 1 {
			logger.Warn("Invalid line: " + line)
		}
		return nil
	}
	key := values[0] + "," + values[3]
	(*networkSet)[key] = struct{}{}
	addr, port := strings.Split(values[1], ":")[0], strings.Split(values[1], ":")[1]
	agentCountry, err := ip2location.ToCountry(ip)
	if err != nil {
		logger.Error("Error getting self country: " + err.Error())
	}
	agentLa, agentLo, err := ip2location.ToLatitudeLongtitude(ip)
	if err != nil {
		logger.Error("Error getting latitude and longtitude: " + err.Error())
	}
	otherCountry, err := ip2location.ToCountry(addr)
	if err != nil {
		logger.Error("Error getting other country: " + err.Error())
	}
	otherLo, otherLa, err := ip2location.ToLatitudeLongtitude(addr)
	if err != nil {
		logger.Error("Error getting latitude and longtitude: " + err.Error())
	}
	var modifiedStr string
	if values[4] == "0" { // in -> i am dst
		modifiedStr = values[0] + "|" + values[3] + "|" + values[2] + "|" + addr + "|" + port + "|" + ip + "|" + values[5] + "|" +
			"unknown" + "|in|detect|" +
			values[5] + "|" + agentCountry + "|" + strconv.Itoa(agentLo) + "|" + strconv.Itoa(agentLa) + "|" +
			addr + "|" + port + "|" + otherCountry + "|" + strconv.Itoa(otherLo) + "|" + strconv.Itoa(otherLa)
	} else { // out -> i am src
		modifiedStr = values[0] + "|" + values[3] + "|" + values[2] + "|" + ip + "|" + values[5] + "|" + addr + "|" + port + "|" +
			"unknown" + "|out|detect|" +
			values[5] + "|" + agentCountry + "|" + strconv.Itoa(agentLo) + "|" + strconv.Itoa(agentLa) + "|" +
			addr + "|" + port + "|" + otherCountry + "|" + strconv.Itoa(otherLo) + "|" + strconv.Itoa(otherLa)
	}
	values = strings.Split(modifiedStr, "|")
	return values
}