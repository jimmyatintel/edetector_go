package clientsearch

import (
	"context"
	config "edetector_go/config"
	fflag "edetector_go/internal/fflag"
	logger "edetector_go/pkg/logger"
	"fmt"

	"go.uber.org/zap"

	// "io/ioutil"
	"net"
	"os"
	"sync"
)

var Client_TCP_Server net.Listener
var Client_detect_TCP_Server net.Listener
var Client_UDP_Server net.PacketConn
var Client_detect_UDP_Server net.PacketConn
var Tcp_enable bool
var Udp_enable bool

func Connect_init() int {
	var err error
	if Tcp_enable, err = fflag.FFLAG.FeatureEnabled("client_tcp"); Tcp_enable && err == nil {
		logger.Info("tcp is enabled")
		Client_TCP_Server, err = net.Listen(config.Viper.GetString("WORKER_SERVER_TYPE_TCP"), "0.0.0.0"+":"+config.Viper.GetString("WORKER_DEFAULT_WORKER_PORT"))
		if err != nil {
			logger.Error("Error listening:", zap.Any("error", err.Error()))
			return 1
		}
		Client_detect_TCP_Server, err = net.Listen(config.Viper.GetString("WORKER_SERVER_TYPE_TCP"), "0.0.0.0"+":"+config.Viper.GetString("WORKER_DEFAULT_DETECT_PORT"))
		if err != nil {
			logger.Error("Error listening:", zap.Any("error", err.Error()))
			return 1
		}
	}
	if Udp_enable, err = fflag.FFLAG.FeatureEnabled("client_udp"); Udp_enable && err == nil {
		logger.Info("udp is enabled")
		Client_UDP_Server, err = net.ListenPacket(config.Viper.GetString("WORKER_SERVER_TYPE_UDP"), "0.0.0.0"+":"+config.Viper.GetString("WORKER_DEFAULT_WORKER_PORT"))
		if err != nil {
			logger.Error("Error listening:", zap.Any("error", err.Error()))
			return 1
		}
		Client_detect_UDP_Server, err = net.ListenPacket(config.Viper.GetString("WORKER_SERVER_TYPE_UDP"), "0.0.0.0"+":"+config.Viper.GetString("WORKER_DEFAULT_DETECT_PORT"))
		if err != nil {
			logger.Error("Error listening:", zap.Any("error", err.Error()))
			return 1
		}
	}
	return 0
}
func Conn_TCP_start(c chan string, wg *sync.WaitGroup) {
	if Client_TCP_Server != nil {
		for {
			conn, err := Client_TCP_Server.Accept()
			if err != nil {
				// fmt.Println("Error accepting: ", err.Error())
				c <- err.Error()
			}
			go handleTCPRequest(conn)
		}
	}
	c <- "TCP Server is nil"
	return
}
func Conn_TCP_detect_start(c chan string, ctx context.Context) {
	if Client_detect_TCP_Server != nil {
		for {
			conn, err := Client_detect_TCP_Server.Accept()
			if err != nil {
				// fmt.Println("Error accepting: ", err.Error())
				c <- err.Error()
			}
			go handleTCPRequest(conn)
		}
	}
	c <- "TCP Server is nil"
	return
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
	return
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
	go Conn_TCP_start(TCP_CHANNEL, wg)
	go Conn_UDP_start(UDP_CHANNEL, wg)
	go Conn_TCP_detect_start(TCP_DETECT_CHANNEL, ctx)
	rt := 0
	if Tcp_enable {
		select {
		case <-ctx.Done():
			fmt.Println("Get quit signal")
			Connection_close(Connection_close_chan)
		case ErrTCP := <-TCP_CHANNEL:
			logger.Error("Error TCP listening:", zap.Any("error", ErrTCP))
			Connection_close(Connection_close_chan)
			rt = 1
		case ErrTCP := <-TCP_DETECT_CHANNEL:
			logger.Error("Error TCP listening:", zap.Any("error", ErrTCP))
			rt = 1
		}
	} else if Udp_enable {
		select {
		case <-ctx.Done():
			fmt.Println("Get quit signal")
			Connection_close(Connection_close_chan)
		case ErrUDP := <-UDP_CHANNEL:
			logger.Error("Error UDP listening:", zap.Any("error", ErrUDP))
			Connection_close(Connection_close_chan)
			rt = 1
		}
	} else {
		select {
		case <-ctx.Done():
			fmt.Println("Get quit signal")
			Connection_close(Connection_close_chan)
		}
	}
	wg.Wait()
	return rt
}
func Connection_close(Connection_close_chan chan<- int) {
	// logger.Info("Connection close")
	if Client_TCP_Server != nil {
		Client_TCP_Server.Close()
	}
	if Client_UDP_Server != nil {
		Client_UDP_Server.Close()
	}
	Connection_close_chan <- 1
}

func Main(ctx context.Context, Connection_close_chan chan<- int) {
	if Connect_init() == 1 {
		logger.Error("Init Connection error")
		os.Exit(1)
	}
	Connect_start(ctx, Connection_close_chan)

	// if Connect_start() == 1 {
	// 	logger.Error("Start Connection error")
	// 	os.Exit(1)
	// }
}
