package messages

import "encoding/binary"

type JoinSignal struct {}

func NewJoinSignal() JoinSignal {
	return JoinSignal{}
}

func DecodeJoinSignal([]byte) JoinSignal {
	return JoinSignal{}
}

func (msg JoinSignal) Encode(buf []byte) uint32 {
	binary.BigEndian.PutUint32(buf, 1)
	buf[4] = JOIN_MSG_ID
	return 5
}

