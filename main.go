package main

import (
	"flag"
	"log"
	"time"

	"github.com/joaovictorsl/downmany/core"
)

func main() {
	runAsServer := flag.Bool("server", false, "Run as server")
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	if *runAsServer {

		core.NewServer(5 * time.Second).Start(":" + *port)
	} else {
		log.Println("Running as client")
	}
}
