package core

import (
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

var port uint16 = 3000

func Connect() {
	sumsMap, err := Sum()
	hashMap = sumsMap
	treat(err)
	go receive()

	var serverIp string = "192.168.1.1:8000"
	serverConnection, err := dowol.NewDowolPeerConn(serverIp) // TODO: Get the ip from CLI param
	treat(err)

	serverConnection.Join(port) // Joins the server, only then starts renewing
	go renew(serverConnection)

	ips, err := serverConnection.GetIPs()
	treat(err)

	clientConnections, failedConnectionsIps := makeConnections(ips)
	if len(failedConnectionsIps) > 0 {
		fmt.Println(failedConnectionsIps) // TODO: Tratar conexões que falharam
	}

	var filehash uint64 = 1337086202
	connectionsWithFile, failedConnections := getClientsWithFile(clientConnections, filehash) // TODO: receber caminho do arquivo na chamada do script e calcular o sum
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
		err := serverConnection.Join(port) // TODO: Get port from CLI
		treat(err)
		time.Sleep(5 * time.Second)
	}
}

func makeConnections(ips []*net.TCPAddr) ([]*dowol.DowolPeerConn, []*net.TCPAddr) {
	connections := make([]*dowol.DowolPeerConn, len(ips))
	failedConnectionsIps := make([]*net.TCPAddr, len(ips))

	for _, ip := range ips {
		connection, err := dowol.NewDowolPeerConn(ip.String())
		if err != nil {
			failedConnectionsIps = append(failedConnectionsIps, ip)
		}
		connections = append(connections, connection)
	}

	return connections, failedConnectionsIps
}

func getClientsWithFile(clientConnections []*dowol.DowolPeerConn, fileHash uint64) ([]*dowol.DowolPeerConn, []*dowol.DowolPeerConn) {
	connectionsWithFile := make([]*dowol.DowolPeerConn, len(clientConnections))
	failedConnections := make([]*dowol.DowolPeerConn, len(clientConnections))

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

// Open connection in the port used in the Join
func receive() {
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
			fmt.Printf("Recebido: %s\n", string(buffer[:n]))
		}
	}
}
