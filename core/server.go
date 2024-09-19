package core

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/joaovictorsl/downmany/network/dowol/messages"
)

type Server struct {
	mapIPs map[*net.TCPAddr]time.Time 
	timeout time.Duration
}

func NewServer(timeout time.Duration) *Server {
	return &Server{
		mapIPs: make(map[*net.TCPAddr]time.Time),
		timeout: timeout,
	}
}

func (s *Server) Start(port uint16) error {
	log.Println("Running as server on port", port)

	go s.cleanIPs()

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))	
	if err != nil { 
		return err
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	b := make([]byte, 2048)

    for {
        // Read the message
        n, err := conn.Read(b)
        if err != nil {
            return
        }

        fmt.Println("server received")
        fmt.Println(b[:n])

        // Decode the message
        id := b[4]
        switch id {
        case messages.JOIN_MSG_ID:
            s.handleJoin(b, conn)
        case messages.GET_IPS_MSG_ID:
            s.handleGetIPs(b, conn)
        default:
            fmt.Printf("Error: unknown message id %d\n", id)
        }
    }
}

func (s *Server) handleJoin(b []byte, conn net.Conn) {
	msg := messages.DecodeJoinRequest(b[5:7])
	addr := conn.RemoteAddr().(*net.TCPAddr)

	addr.Port = int(msg.Port())
	s.mapIPs[addr] = time.Now()

    n := messages.NewJoinResponse().Encode(b)
    conn.Write(b[:n])
}

func (s *Server) handleGetIPs(b []byte, conn net.Conn) {
	ips := make([]*net.TCPAddr, 0, len(s.mapIPs))

	for addr := range s.mapIPs {
		if addr == conn.RemoteAddr() {
			continue
		}

        ips = append(ips, addr)
	}

	resp := messages.NewGetIPsResponse(ips)
	n := resp.Encode(b)
	conn.Write(b[:n])
}

func (s *Server) cleanIPs() {
	for {
		time.Sleep(s.timeout)

		for addr, t := range s.mapIPs {
			if s.timedOut(t) {
				delete(s.mapIPs, addr)
			}
		}
	}

}

func (s *Server) timedOut(t time.Time) bool {
	return time.Since(t) > s.timeout 
}

