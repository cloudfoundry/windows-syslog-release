package main

import (
	"log"

	"github.com/jtarchie/syslog/pkg/transports"
	"github.com/jtarchie/syslog/pkg/writers/web"
)

func main() {
	log.Println("starting servers")
	writer := writers.NewServer(8081)
	go func() {
		log.Fatalf("Could not start writer: %s", writer.Start())
	}()

	server, err := transports.NewUDPServer(8088, writer)
	if err != nil {
		log.Fatalf("Could not start server: %s", err)
	}

	server.Start()
}
