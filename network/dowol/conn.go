package dowol

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"

	"github.com/joaovictorsl/downmany/network/dowol/messages"
)

type DowolPeerConn struct {
	addr       net.Addr
	conn       net.Conn
	lengthBuf  []byte
	payloadBuf []byte
	mutex      *sync.Mutex
}

func NewDowolPeerConn(addr string) (*DowolPeerConn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &DowolPeerConn{
		addr:       conn.RemoteAddr(),
		conn:       conn,
		lengthBuf:  make([]byte, 4),
		payloadBuf: make([]byte, 10000),
		mutex:      &sync.Mutex{},
	}, nil
}

func (dpc *DowolPeerConn) readMessage() (uint32, error) {
	n, err := dpc.conn.Read(dpc.lengthBuf)
	if err != nil {
		return 0, err
	} else if n != 4 {
		return 0, fmt.Errorf("expected to read 4 bytes into length buf, got %d", n)
	}

	bytesToRead := binary.BigEndian.Uint32(dpc.lengthBuf)
	bytesRead := uint32(0)

	for bytesRead != bytesToRead {
		n, err := dpc.conn.Read(dpc.payloadBuf[bytesRead:])
		if err != nil {
			return 0, err
		}

		bytesRead += uint32(n)
	}

	return bytesRead, nil
}

func (dpc *DowolPeerConn) Join(port uint16) error {
	dpc.mutex.Lock()
	defer dpc.mutex.Unlock()

	signal := messages.NewJoinSignal(port)
	n := signal.Encode(dpc.payloadBuf)
	_, err := dpc.conn.Write(dpc.payloadBuf[1:n])
	return err
}

func (dpc *DowolPeerConn) GetIPs() ([]*net.TCPAddr, error) {
	dpc.mutex.Lock()
	defer dpc.mutex.Unlock()

	req := messages.NewGetIPsRequest()
	n := req.Encode(dpc.payloadBuf)

	_, err := dpc.conn.Write(dpc.payloadBuf[:n])
	if err != nil {
		return nil, err
	}

	n, err = dpc.readMessage()
	if err != nil {
		return nil, err
	}

	if dpc.payloadBuf[0] != messages.GET_IPS_MSG_ID {
		return nil, fmt.Errorf("expected to receive message with id %d, got %d", messages.HAS_MSG_ID, dpc.payloadBuf[0])
	}

	res := messages.DecodeGetIPsResponse(dpc.payloadBuf[1:n])

	return res.IPs, nil
}

func (dpc *DowolPeerConn) HasFile(hash uint64) (bool, error) {
	dpc.mutex.Lock()
	defer dpc.mutex.Unlock()

	req := messages.NewHasFileRequest(hash)
	n := req.Encode(dpc.payloadBuf)

	_, err := dpc.conn.Write(dpc.payloadBuf[:n])
	if err != nil {
		return false, err
	}

	n, err = dpc.readMessage()
	if err != nil {
		return false, err
	}

	if dpc.payloadBuf[0] != messages.HAS_MSG_ID {
		return false, fmt.Errorf("expected to receive message with id %d, got %d", messages.HAS_MSG_ID, dpc.payloadBuf[0])
	}

	res := messages.DecodeHasFileResponse(dpc.payloadBuf[1:n])

	return res.Has, nil
}
