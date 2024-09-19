package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/joaovictorsl/downmany/core"
)

//var fileHash uint64 = 1336902055

type cliArgs struct {
	runAsServer    bool
	port           uint16
	clientPort     uint16
	timeoutSeconds uint
	fileHash       uint64
	serverAddr     string
}

func getArgs() cliArgs {
	runAsServer := flag.Bool("server", false, "Run as server")
	port := flag.String("port", "8000", "Port to listen on")
	clientPort := flag.String("clientPort", "8080", "Client port to listen on")
	timeoutSeconds := flag.Uint("timeout", 5, "Timeout in seconds to be used when cleaning clients")
	fileHash := flag.Uint64("file_hash", 0, "Hash of the file to be downloaded")
	serverAddr := flag.String("server_addr", "", "Address of the tracker server")
	flag.Parse()

	if !*runAsServer {
		if *fileHash == 0 {
			panic("Client should specify file hash for download")
		} else if *serverAddr == "" {
			panic("Client should specify server address")
		}
	}

	portNumber, err := strconv.ParseUint(*port, 10, 16)
	if err != nil {
		panic(err)
	}
	clientPortNumber, err := strconv.ParseUint(*clientPort, 10, 16)
	if err != nil {
		panic(err)
	}

	return cliArgs{
		*runAsServer,
		uint16(portNumber),
		uint16(clientPortNumber),
		*timeoutSeconds,
		*fileHash,
		*serverAddr,
	}
}

func main() {
	args := getArgs()

	if args.runAsServer {
		timeout := time.Duration(args.timeoutSeconds) * time.Second
		server := core.NewServer(timeout)
		if err := server.Start(args.port); err != nil {
			log.Fatal(err)
		}
	} else {
		done := make(chan bool)
		log.Println("Running as client")

		ips := core.Connect(args.serverAddr, args.clientPort)
		fmt.Println("Connected with the server")

		if args.fileHash != 0 {
			fmt.Printf("Requesting file with hash %v\n", args.fileHash)
			connectionsWithFile, failedHasFileConnections, failedConnectionsIps := core.AskForFile(ips, args.fileHash)

			if len(failedConnectionsIps) > 0 {
				fmt.Println("Failed connections:")
				fmt.Println(failedConnectionsIps) // TODO: Tratar conexões cuja requisição HasFile() falhou
			}

			if len(failedHasFileConnections) > 0 {
				fmt.Println("Failed HasFile connections:")
				fmt.Println(failedHasFileConnections) // TODO: Tratar conexões cuja requisição HasFile() falhou
			}

			if len(connectionsWithFile) > 0 {
				fmt.Println("Connections with file:")
				for _, connection := range connectionsWithFile {
					fmt.Println(connection.GetAddr().String())
				}
			} else {
				fmt.Println("There are no connections with the file you requested")
			}
		}

		<-done

	}
}
