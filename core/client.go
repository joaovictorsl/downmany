package core

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/joaovictorsl/downmany/network/dowol"
)

var hashMap map[uint64]string

func treat(err error) {
	if err != nil {
		panic(err)
	}
}

func Connect(serverAddr string, port uint16) []*net.TCPAddr {
	sumsMap, err := Sum("./dataset")
	hashMap = sumsMap
	treat(err)
	go openTCPPort(port)

	serverConnection, err := dowol.NewDowolPeerConn(serverAddr)
	treat(err)

	serverConnection.Join(port) // Joins the server, only then starts renewing
	go renew(serverConnection, port)

	ips, err := serverConnection.GetIPs()
	treat(err)

	return ips
}

func AskForFile(ips []*net.TCPAddr, fileHash uint64) ([]*dowol.DowolPeerConn, []*dowol.DowolPeerConn, []*net.TCPAddr) {

	clientConnections, failedConnectionsIps := makeConnections(ips) // TODO: add channels and goroutines
	if len(failedConnectionsIps) > 0 {
		// TODO: Tratar conexões que falharam
	}

	clientsWithFile, clientsWithoutFile := getClientsWithFile(clientConnections, fileHash)

	return clientsWithFile, clientsWithoutFile, failedConnectionsIps // TODO: add channels, goroutines, and select statement

}

func renew(serverConnection *dowol.DowolPeerConn, port uint16) {
	for {
		// O parâmetro de porta só é necessário na primeira chamada do join
		err := serverConnection.Join(port) // TODO: Get port from CLI
		treat(err)
		time.Sleep(5 * time.Second)
	}
}

func makeConnections(ips []*net.TCPAddr) ([]*dowol.DowolPeerConn, []*net.TCPAddr) {
	connections := make([]*dowol.DowolPeerConn, 0)
	failedConnectionsIps := make([]*net.TCPAddr, 0)
	for _, ip := range ips {
		connection, err := dowol.NewDowolPeerConn(ip.String())
		if err != nil {
			failedConnectionsIps = append(failedConnectionsIps, ip)
		} else {
			connections = append(connections, connection)
		}
	}

	return connections, failedConnectionsIps
}

func getClientsWithFile(clientConnections []*dowol.DowolPeerConn, fileHash uint64) ([]*dowol.DowolPeerConn, []*dowol.DowolPeerConn) {
	connectionsWithFile := make([]*dowol.DowolPeerConn, 0)
	failedHasFileConnections := make([]*dowol.DowolPeerConn, 0)

	for _, clientConnection := range clientConnections {
		hasFile, err := clientConnection.HasFile(fileHash)
		if err != nil {
			failedHasFileConnections = append(failedHasFileConnections, clientConnection)
		} else if hasFile {
			connectionsWithFile = append(connectionsWithFile, clientConnection)
		}
	}
	return connectionsWithFile, failedHasFileConnections
}

func openTCPPort(port uint16) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	treat(err)
	defer listener.Close()

	fmt.Printf("Escutando na porta %d...\n", port)

	for {
		conn, err := listener.Accept()
		treat(err)

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var bufferCapacity uint16 = 1024
	buffer := make([]byte, bufferCapacity)
	for {
		n, err := conn.Read(buffer)
		treat(err)

		if n > 0 {
			hashFile := binary.BigEndian.Uint64(buffer[5:n])

			response := []byte{0, 0, 0, 2, 3, 0}

			if hashMap[hashFile] != "" {
				response[5] = 1
			}

			conn.Write(response)
		}
	}
}
