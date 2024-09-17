package core

import (
	"fmt"
	"net"
	"time"

	"github.com/joaovictorsl/downmany/network/dowol"
)

var hashMap map[uint64]string

func Connect() {
	sumsMap, err := Sum()
	hashMap = sumsMap
	if err != nil {
		panic(err)
	}

	serverConnection, err := dowol.NewDowolPeerConn("192.168.1.1:8000") // TODO: Get the ip from CLI param
	if err != nil {
		panic(err)
	}

	go renew(serverConnection)

	ips, err := serverConnection.GetIPs()
	if err != nil {
		panic(err)
	}

	clientConnections, failedConnectionsIps := makeConnections(ips)
	if len(failedConnectionsIps) > 0 {
		fmt.Println(failedConnectionsIps) // TODO: Tratar conexões que falharam
	}

	connectionsWithFile, failedConnections := getClientsWithFile(clientConnections, 123454312) // TODO: receber caminho do arquivo na chamada do script e calcular o sum
	if len(failedConnections) > 0 {
		fmt.Println(failedConnections) // TODO: Tratar conexões cuja requisição HasFile() falhou
	}

	for _, connection := range connectionsWithFile {
		fmt.Println(connection.GetAddr().String())
	}
}

func renew(serverConnection *dowol.DowolPeerConn) {
	for {
		// O parâmetro de porta só é necessário na primeira chamada do join
		serverConnection.Join(3000) // TODO: Get port from CLI
		time.Sleep(5 * time.Second)
	}
}

func makeConnections(ips []net.Addr) ([]*dowol.DowolPeerConn, []net.Addr) {
	connections := make([]*dowol.DowolPeerConn, len(ips))
	failedConnectionsIps := make([]net.Addr, len(ips))

	for _, ip := range ips {
		connection, err := dowol.NewDowolPeerConn(ip.String())
		if err != nil {
			failedConnectionsIps = append(failedConnectionsIps, ip)
		}
		connections = append(connections, connection)
	}

	return connections, failedConnectionsIps
}

func getClientsWithFile(clientConnections []*dowol.DowolPeerConn, fileSum uint64) ([]*dowol.DowolPeerConn, []*dowol.DowolPeerConn) {
	connectionsWithFile := make([]*dowol.DowolPeerConn, len(clientConnections))
	failedConnections := make([]*dowol.DowolPeerConn, len(clientConnections))

	for _, clientConnection := range clientConnections {
		hasFile, err := clientConnection.HasFile(fileSum)
		if err != nil {
			failedConnections = append(failedConnections, clientConnection)
		} else if hasFile {
			connectionsWithFile = append(connectionsWithFile, clientConnection)
		}
	}
	return connectionsWithFile, failedConnections
}

// Open connection in the port used in the Join
