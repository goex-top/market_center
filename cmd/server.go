package main

import (
	"context"
	. "github.com/goex-top/market_center"
	"github.com/goex-top/market_center/api"
	"github.com/goex-top/market_center/config"
	"github.com/goex-top/market_center/data"
	log "github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	Api    *api.Api
	Ctx    context.Context
	Cancel func()
	Cfg    *config.Config
	Data   *data.Data
	logger = log.New()
)

//生成随机字符串
func getRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := make([]byte, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func HandleConn(c net.Conn, id string) {
	logger.Infof("%s - connected, RemoteAddr:%s, LocalAddr:%s", id, c.RemoteAddr().String(), c.LocalAddr().String())
	for {
		buf := make([]byte, 1024)
		count, err := c.Read(buf[:])
		if err != nil {
			if err != io.EOF {
				logger.Errorf("Error on read: %s", err.Error())
			}
			logger.Infof("%s - disconnected", id)
			return
		} else if count > 0 {
			ProcessMessage(c, buf[:count])
		}
	}
}

func main() {
	os.Remove(UDS_PATH)
	logger.SetLevel(log.DebugLevel)
	ln, err := net.Listen("unix", UDS_PATH)
	if err != nil {
		logger.Fatal("Listen error: ", err)
	}
	defer os.Remove(UDS_PATH)
	logger.Printf("Starting USD server:%s", UDS_PATH)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(ln net.Listener, c chan os.Signal) {
		sig := <-c
		Cancel()
		logger.Printf("Caught signal %s: shutting down.", sig)
		ln.Close()
		os.Remove(UDS_PATH)
		os.Exit(0)
	}(ln, sigc)

	for {
		fd, err := ln.Accept()
		if err != nil {
			logger.Fatal("Accept error: ", err)
		}
		id := getRandomString(10)
		go HandleConn(fd, id)
	}
}
