package messages

import "encoding/binary"

type JoinSignal struct {
	port uint16
}

func NewJoinSignal(port uint16) JoinSignal {
	return JoinSignal{
		port: port,
	}
}

func DecodeJoinSignal(buf []byte) JoinSignal {
	port := binary.BigEndian.Uint16(buf)
	return JoinSignal{
		port: port,
	}
}

func (msg JoinSignal) Encode(buf []byte) uint32 {
	binary.BigEndian.PutUint32(buf, 3)
	buf[4] = JOIN_MSG_ID
	binary.BigEndian.PutUint16(buf[5:], msg.port)
	return 7
}
