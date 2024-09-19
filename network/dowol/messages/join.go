package messages

import "encoding/binary"

type JoinRequest struct {
	port uint16
}

func NewJoinRequest(port uint16) JoinRequest {
	return JoinRequest{
		port: port,
	}
}

func DecodeJoinRequest(buf []byte) JoinRequest {
	port := binary.BigEndian.Uint16(buf)
	return JoinRequest{
		port: port,
	}
}

func (msg JoinRequest) Encode(buf []byte) uint32 {
	binary.BigEndian.PutUint32(buf, 3)
	buf[4] = JOIN_MSG_ID
	binary.BigEndian.PutUint16(buf[5:], msg.port)
	return 7
}

func (msg JoinRequest) Port() uint16 {
	return msg.port
}

type JoinResponse struct {}

func NewJoinResponse() JoinResponse {
	return JoinResponse{}
}

func DecodeJoinResponse(buf []byte) JoinResponse {
	return JoinResponse{}
}

func (msg JoinResponse) Encode(buf []byte) uint32 {
	binary.BigEndian.PutUint32(buf, 1)
	buf[4] = JOIN_MSG_ID
	return 5
}

