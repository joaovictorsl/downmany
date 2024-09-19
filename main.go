package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/joaovictorsl/downmany/core"
)

var fileHash uint64 = 1336902055 // TODO: Get from CLI

func main() {
	runAsServer := flag.Bool("server", false, "Run as server")
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	if *runAsServer {

		core.NewServer(5 * time.Second).Start(":" + *port)
	} else {
		log.Println("Running as client")

		ips := core.Connect()
		connectionsWithFile, failedConnections := core.AskForFile(ips, fileHash)

		if len(failedConnections) > 0 {
			fmt.Println(failedConnections) // TODO: Tratar conexões cuja requisição HasFile() falhou
		}

		for _, connection := range connectionsWithFile {
			fmt.Println(connection.GetAddr().String())
		}

	}
}
