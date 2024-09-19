package core

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/joaovictorsl/downmany/network/dowol"
)

var hashMap map[uint64]string

func treat(err error) {
	if err != nil {
		panic(err)
	}
}

var port uint16 = 3000

func Connect() []*net.TCPAddr {
	sumsMap, err := Sum("./dataset")
	hashMap = sumsMap
	treat(err)
	go openTCPPort()

	// var serverIp string = "192.168.1.1:8000"
	// serverConnection, err := dowol.NewDowolPeerConn(serverIp) // TODO: Get the ip from CLI param
	// treat(err)

	// serverConnection.Join(port) // Joins the server, only then starts renewing
	// go renew(serverConnection)

	// ips, err := serverConnection.GetIPs()
	// treat(err)

	// Mocked server
	ips := []*net.TCPAddr{
		{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 4000,
		},
	}

	return ips
}

func AskForFile(ips []*net.TCPAddr, fileHash uint64) ([]*dowol.DowolPeerConn, []*dowol.DowolPeerConn) {

	clientConnections, failedConnectionsIps := makeConnections(ips) // TODO: add channels and goroutines
	if len(failedConnectionsIps) > 0 {
		fmt.Println(failedConnectionsIps) // TODO: Tratar conexões que falharam
	}

	return getClientsWithFile(clientConnections, fileHash) // TODO: add channels, goroutines and select statement

}

// func renew(serverConnection *dowol.DowolPeerConn) {
// 	for {
// 		// O parâmetro de porta só é necessário na primeira chamada do join
// 		err := serverConnection.Join(port) // TODO: Get port from CLI
// 		treat(err)
// 		time.Sleep(5 * time.Second)
// 	}
// }

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
	failedConnections := make([]*dowol.DowolPeerConn, 0)

	for _, clientConnection := range clientConnections {
		hasFile, err := clientConnection.HasFile(fileHash)
		if err != nil {
			failedConnections = append(failedConnections, clientConnection)
		} else if hasFile {
			connectionsWithFile = append(connectionsWithFile, clientConnection)
		}
	}
	return connectionsWithFile, failedConnections
}

func openTCPPort() {
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
