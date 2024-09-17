package messages

import (
	"encoding/binary"
	"net"
)

type GetIPsRequest struct{}

func NewGetIPsRequest() GetIPsRequest {
	return GetIPsRequest{}
}

func DecodeGetIPsRequest(raw []byte) GetIPsRequest {
	return GetIPsRequest{}
}

func (msg GetIPsRequest) Encode(buf []byte) uint32 {
	binary.BigEndian.PutUint32(buf, 1)
	buf[4] = GET_IPS_MSG_ID
	return 5
}

type GetIPsResponse struct {
	IPs []net.Addr
}

func NewGetIPsResponse(ips []net.Addr) GetIPsResponse {
	return GetIPsResponse{
		IPs: ips,
	}
}

func DecodeGetIPsResponse(raw []byte) GetIPsResponse {
	totalIPs := len(raw) / 6 // an IP takes 6 bytes
	ips := make([]net.Addr, totalIPs)

	for i := 0; i < totalIPs; i++ {
		startAddr := i * 6
		startPort := startAddr + 4 // 4 is length of addr
		ips[i] = &net.TCPAddr{
			IP:   raw[startAddr:startPort],
			Port: int(binary.BigEndian.Uint16(raw[startPort : startPort+2])),
		}
	}

	return GetIPsResponse{
		ips,
	}
}

func (msg GetIPsResponse) Encode(buf []byte) uint32 {
	payloadLen := uint32(len(msg.IPs)) * 6 // 6 bytes per IP (addr + port)
	payloadLen = payloadLen + 1            // byte for msg id
	// Put payload length
	binary.BigEndian.PutUint32(buf, payloadLen)
	buf[4] = GET_IPS_MSG_ID

	start := 5
	for _, ip := range msg.IPs {
		tcpIp := ip.(*net.TCPAddr)
		copy(buf[start:], tcpIp.IP)
		binary.BigEndian.PutUint16(buf[start+4:], uint16(tcpIp.Port))
		start += 6
	}

	return 4 + payloadLen
}
