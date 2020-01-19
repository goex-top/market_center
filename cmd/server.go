package main

import (
	"context"
	"fmt"
	. "github.com/goex-top/market_center"
	"github.com/goex-top/market_center/api"
	"github.com/goex-top/market_center/config"
	"github.com/goex-top/market_center/data"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	Api  *api.Api
	Ctx  context.Context
	Cfg  *config.Config
	Data *data.Data
)

func HandleConn(c net.Conn) {
	for {
		buf := make([]byte, 1024)
		count, err := c.Read(buf[:])
		if err != nil {
			fmt.Println("err:", err)
			if err != io.EOF {
				fmt.Errorf("Error on read: %s", err.Error())
			}
			return
		} else if count > 0 {
			ProcessMessage(c, buf[:count])
		}
	}
}

func main() {
	os.Remove(UDS_PATH)
	log.Println("Starting USD server")
	ln, err := net.Listen("unix", UDS_PATH)
	if err != nil {
		log.Fatal("Listen error: ", err)
	}
	defer os.Remove(UDS_PATH)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(ln net.Listener, c chan os.Signal) {
		sig := <-c
		log.Printf("Caught signal %s: shutting down.", sig)
		ln.Close()
		os.Exit(0)
	}(ln, sigc)

	for {
		fd, err := ln.Accept()
		if err != nil {
			log.Fatal("Accept error: ", err)
		}

		go HandleConn(fd)
	}
}
