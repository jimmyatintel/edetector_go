package clientsearch

import (
	"context"
	config "edetector_go/config"
	channelmap "edetector_go/internal/channelmap"
	packet "edetector_go/internal/packet"
	"edetector_go/internal/taskservice"
	logger "edetector_go/pkg/logger"
	mq "edetector_go/pkg/mariadb/query"
	rq "edetector_go/pkg/redis/query"
	"time"

	// "io/ioutil"
	"net"
	"os"
	"sync"
)

var Client_TCP_Server net.Listener
var Client_detect_TCP_Server net.Listener
var Task_server_TCP_Server net.Listener
var Client_UDP_Server net.PacketConn
var Client_detect_UDP_Server net.PacketConn
var Tcp_enable bool
var Udp_enable bool
var Task_enable bool

func Connect_init() int {
	var err error
	if true {
		logger.Info("TCP is enabled")
		Client_TCP_Server, err = net.Listen(config.Viper.GetString("WORKER_SERVER_TYPE_TCP"), "0.0.0.0"+":"+config.Viper.GetString("WORKER_DEFAULT_WORKER_PORT"))
		if err != nil {
			logger.Panic("Error listening: " + err.Error())
			panic(err)
		}
		Client_detect_TCP_Server, err = net.Listen(config.Viper.GetString("WORKER_SERVER_TYPE_TCP"), "0.0.0.0"+":"+config.Viper.GetString("WORKER_DEFAULT_DETECT_PORT"))
		if err != nil {
			logger.Panic("Error listening: " + err.Error())
			panic(err)
		}
	}
	// if Udp_enable, err = fflag.FFLAG.FeatureEnabled("client_udp"); Udp_enable && err == nil {
	// 	logger.Info("UDP is enabled")
	// 	Client_UDP_Server, err = net.ListenPacket(config.Viper.GetString("WORKER_SERVER_TYPE_UDP"), "0.0.0.0"+":"+config.Viper.GetString("WORKER_DEFAULT_WORKER_PORT"))
	// 	if err != nil {
	// 		logger.Panic("Error listening: " + err.Error())
	// 		panic(err)
	// 	}
	// 	Client_detect_UDP_Server, err = net.ListenPacket(config.Viper.GetString("WORKER_SERVER_TYPE_UDP"), "0.0.0.0"+":"+config.Viper.GetString("WORKER_DEFAULT_DETECT_PORT"))
	// 	if err != nil {
	// 		logger.Panic("Error listening: " + err.Error())
	// 		panic(err)
	// 	}
	// }
	return 0
}
func Conn_TCP_start(c chan string, wg *sync.WaitGroup) {
	channelmap.TaskWorkerChannel = make(map[string](*chan packet.Packet))
	if Client_TCP_Server != nil {
		for {
			conn, err := Client_TCP_Server.Accept()
			if err != nil {
				logger.Error("Error accepting: " + err.Error())
				c <- err.Error()
			}
			new_task_chan := make(chan packet.Packet)
			go handleTCPRequest(conn, new_task_chan, "worker")
		}
	}
	c <- "TCP Server is nil"
}
func Conn_TCP_detect_start(c chan string, ctx context.Context) {
	if Client_detect_TCP_Server != nil {
		for {
			conn, err := Client_detect_TCP_Server.Accept()
			if err != nil {
				logger.Error("Error accepting: " + err.Error())
				c <- err.Error()
			}
			logger.Info("Detect port accepted")
			go handleTCPRequest(conn, nil, "detect")
		}
	}
	c <- "TCP Server is nil"
}
func Conn_UDP_start(c chan string, wg *sync.WaitGroup) {
	if Client_UDP_Server != nil {
		for {
			buf := make([]byte, 65536)
			n, addr, err := Client_UDP_Server.ReadFrom(buf)
			if err != nil {
				c <- err.Error()
			}
			go handleUDPRequest(addr, buf[:n])
		}
	}
	c <- "UDP Server is nil"
}

func Connect_start(ctx context.Context, Connection_close_chan chan<- int) int {
	wg := new(sync.WaitGroup)
	defer wg.Done()
	wg.Add(2)
	TCP_CHANNEL := make(chan string)
	TCP_DETECT_CHANNEL := make(chan string)
	UDP_CHANNEL := make(chan string)
	defer close(TCP_CHANNEL)
	defer close(UDP_CHANNEL)
	Offline_all_clients()
	go checkOnline()
	go Conn_TCP_start(TCP_CHANNEL, wg)
	go Conn_UDP_start(UDP_CHANNEL, wg)
	go Conn_TCP_detect_start(TCP_DETECT_CHANNEL, ctx)
	go taskservice.Start(ctx)
	rt := 0
	if Tcp_enable {
		select {
		case <-ctx.Done():
			logger.Info("Get quit signal")
			Connection_close(Connection_close_chan)
		case ErrTCP := <-TCP_CHANNEL:
			logger.Error("Error TCP listening: " + ErrTCP)
			Connection_close(Connection_close_chan)
			rt = 1
		case ErrTCP := <-TCP_DETECT_CHANNEL:
			logger.Error("Error TCP listening: " + ErrTCP)
			rt = 1
		}
	} else if Udp_enable {
		select {
		case <-ctx.Done():
			logger.Info("Get quit signal")
			Connection_close(Connection_close_chan)
		case ErrUDP := <-UDP_CHANNEL:
			logger.Error("Error UDP listening: " + ErrUDP)
			Connection_close(Connection_close_chan)
			rt = 1
		}
	} else {
		select {
		case <-ctx.Done():
			logger.Info("Get quit signal")
			Connection_close(Connection_close_chan)
		}
	}
	wg.Wait()
	return rt
}
func Connection_close(Connection_close_chan chan<- int) {
	logger.Info("Connection close")
	if Client_TCP_Server != nil {
		Client_TCP_Server.Close()
	}
	if Client_UDP_Server != nil {
		Client_UDP_Server.Close()
	}
	Connection_close_chan <- 1
}

func Offline_all_clients() {
	logger.Info("Offline all clients")
	clients := mq.Load_all_client()
	for _, client := range clients {
		rq.Offline(client, &ClientCount)
	}
}

func checkOnline() {
	for {
		clients := mq.Load_all_client()
		for _, client := range clients {
			if rq.GetStatus(client) == 0 { // offline
				continue
			}
			redisTime, err := time.Parse(time.RFC3339, rq.GetTime(client))
			if err != nil {
				logger.Error("Error parsing time: " + err.Error())
				continue
			}
			currentTime := time.Now()
			difference := currentTime.Sub(redisTime)
			if difference > 65*time.Second {
				logger.Info("Offline: " + client + "- more than 65 seconds without CheckConnect")
				rq.Offline(client, &ClientCount)
			}
		}
		time.Sleep(60 * time.Second)
	}

}

func Main(ctx context.Context, Connection_close_chan chan<- int) {
	if Connect_init() == 1 {
		logger.Panic("Init Connection error")
		os.Exit(1)
	}
	Connect_start(ctx, Connection_close_chan)
}
