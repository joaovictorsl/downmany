package network

import "net"

type PeerConn interface {
	Join() error
	GetIPs() ([]net.Addr, error)
	HasFile(hash uint64) (bool, error)
	Download(start, end uint64) error
}

