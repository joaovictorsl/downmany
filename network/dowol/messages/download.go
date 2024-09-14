package messages

import "encoding/binary"

type DownloadRequest struct {
	Start uint64
	End uint64
}

func NewDownloadRequest(start, end uint64) DownloadRequest {
	return DownloadRequest{
		Start: start,
		End: end,
	}
}

func DecodeDownloadRequest(raw []byte) DownloadRequest {
	return DownloadRequest{
		Start: binary.BigEndian.Uint64(raw),
		End: binary.BigEndian.Uint64(raw[8:]),
	}
}

func (msg DownloadRequest) Encode(buf []byte) uint32 {
	binary.BigEndian.AppendUint32(buf, 16)
	binary.BigEndian.PutUint64(buf[8:], msg.Start)
	binary.BigEndian.PutUint64(buf[16:], msg.End)
	return 20
}

type DownloadResponse struct {
	Data []byte
}

func NewDownloadResponse(data []byte) DownloadResponse {
	return DownloadResponse{
		Data: data,
	}
}

func DecodeDownloadResponse(raw []byte) DownloadResponse {
	return DownloadResponse{
		Data: raw,
	}
}

func (msg DownloadResponse) Encode(buf []byte) uint32 {
	payloadLen := uint32(1 + len(msg.Data))
	binary.BigEndian.PutUint32(buf, payloadLen)
	buf[4] = DOWNLOAD_MSG_ID
	return 0
}

