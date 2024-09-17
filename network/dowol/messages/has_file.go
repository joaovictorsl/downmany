package messages

import "encoding/binary"

type HasFileRequest struct {
	Hash uint64
}

func NewHasFileRequest(hash uint64) HasFileRequest {
	return HasFileRequest{
		Hash: hash,
	}
}

func DecodeHasFileRequest(raw []byte) HasFileRequest {
	return  HasFileRequest{
		Hash: binary.BigEndian.Uint64(raw),
	}
}

func (msg HasFileRequest) Encode(buf []byte) uint32 {
	payloadLen := uint32(1 + 8) // 1 byte for msg id and 8 bytes for hash
	binary.BigEndian.PutUint32(buf, payloadLen)
	buf[4] = HAS_MSG_ID
	binary.BigEndian.PutUint64(buf[5:], msg.Hash)
	return 4 + payloadLen
}

type HasFileResponse struct {
	Has bool
}

func NewHasFileResponse(has bool) HasFileResponse {
	return HasFileResponse{
		Has: has,
	}
}

func DecodeHasFileResponse(raw []byte) HasFileResponse {
	return  HasFileResponse{
		Has: raw[0] == 1,
	}
}

func (msg HasFileResponse) Encode(buf []byte) uint32 {
	payloadLen := uint32(1 + 1) // 1 byte for msg id and 1 byte for has flag
	binary.BigEndian.PutUint32(buf, payloadLen)
	buf[4] = HAS_MSG_ID
	
	if msg.Has {
		buf[5] = 1
	} else {
		buf[5] = 0
	}

	return 4 + payloadLen
}
