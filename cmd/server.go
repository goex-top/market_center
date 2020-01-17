package main

import (
	"context"
	"io"
	"log"
	"market_center/api"
	"market_center/config"
	"market_center/data"
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
	received := make([]byte, 0)
	for {
		buf := make([]byte, 512)
		count, err := c.Read(buf)
		received = append(received, buf[:count]...)
		if err != nil {
			ProcessMessage(c, received)
			if err != io.EOF {
				log.Printf("Error on read: %s", err)
			}
			break
		}
	}
}

func main() {
	log.Println("Starting USD server")
	ln, err := net.Listen("unix", "/tmp/goex.market.center")
	if err != nil {
		log.Fatal("Listen error: ", err)
	}

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
