package messages_test

import (
	"net"
	"slices"
	"testing"

	"github.com/joaovictorsl/downmany/network/dowol/messages"
)

func compareAddrSlice(a1, a2 *net.TCPAddr) int {
	if slices.Compare(a1.IP, a2.IP) != 0 {
		return 1
	} else if a1.Port != a2.Port {
		return 1
	}

	return 0
}

func TestGetIPsRequestEncode(t *testing.T) {
	// Setup
	expected := []byte{0, 0, 0, 1, 1}
	buf := make([]byte, len(expected))
	req := messages.NewGetIPsRequest()
	// Action
	n := req.Encode(buf)
	// Assert
	if n != uint32(len(expected)) {
		t.Errorf("expected GetIPsRequest's encoding to write %d bytes, wrote %d", len(expected), n)
	} else if slices.Compare(expected, buf) != 0 {
		t.Errorf("expected GetIPsRequest's encoding to be %v, got %v", expected, buf)
	}
}

func TestNewGetIPsResponse(t *testing.T) {
	// Setup
	expected := []*net.TCPAddr{
		&net.TCPAddr{
			IP: []byte{192, 168, 0, 1},
			Port: 80,
		},
		&net.TCPAddr{
			IP: []byte{192, 123, 1, 2},
			Port: 3000,
		},
	}
	// Action
	res := messages.NewGetIPsResponse(expected)
	// Assert
	if slices.CompareFunc(expected, res.IPs, compareAddrSlice) != 0 {
		t.Errorf("expected GetIPsResponse's `IPs` field to be %v, got %v", expected, res.IPs)
	}
}

func TestDecodeGetIPsResponse(t *testing.T) {
	// Setup
	expected := []*net.TCPAddr{
		&net.TCPAddr{
			IP: []byte{192, 168, 0, 1},
			Port: 80,
		},
		&net.TCPAddr{
			IP: []byte{192, 168, 0, 2},
			Port: 82,
		},
	}
	raw := []byte{192, 168, 0, 1, 0, 80, 192, 168, 0, 2, 0, 82}
	// Action
	res := messages.DecodeGetIPsResponse(raw)
	// Assert
	if slices.CompareFunc(expected, res.IPs, compareAddrSlice) != 0 {
		t.Errorf("expected GetIPsResponse's `IPs` field to be %v, got %v", expected, res.IPs)
	}
}

func TestEncodeGetIPsResponse(t *testing.T) {
	// Setup
	expected := []byte{0, 0, 0, 13, 1, 192, 168, 0, 1, 0, 80, 192, 168, 0, 2, 0, 82}
	addrs := []*net.TCPAddr{
		&net.TCPAddr{
			IP: []byte{192, 168, 0, 1},
			Port: 80,
		},
		&net.TCPAddr{
			IP: []byte{192, 168, 0, 2},
			Port: 82,
		},
	}
	res := messages.NewGetIPsResponse(addrs)
	buf := make([]byte, 4 + 1 + len(addrs) * 6)
	// Action
	n := res.Encode(buf)
	// Assert
	if n != uint32(len(buf)) {
		t.Errorf("expected GetIPsResponse's encoding to write %d bytes, wrote %d", len(buf), n)
	} else if slices.Compare(expected, buf) != 0 {
		t.Errorf("expected GetIPsResponse's encoding to be %v, got %v", expected, buf)
	}
}

