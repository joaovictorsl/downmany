package messages_test

import (
	"slices"
	"testing"

	"github.com/joaovictorsl/downmany/network/dowol/messages"
)

func TestNewHasFileRequest(t *testing.T) {
	// Setup
	expected := uint64(1010)
	//Action
	req := messages.NewHasFileRequest(expected)
	// Assert
	if req.Hash != expected {
		t.Errorf("expected HasFileRequest's `hash` field to be %d, got %d", expected, req.Hash)
	}
}

func TestDecodeHasFileRequest(t *testing.T) {
	// Setup
	expected := uint64(592)
	buf := []byte{0, 0, 0, 0, 0, 0, 2, 80}
	// Action
	req := messages.DecodeHasFileRequest(buf)
	// Assert
	if req.Hash != expected {
		t.Errorf("expected HasFileRequest's `hash` field to be %d when decoding, got %d", expected, req.Hash)
	}
}

func TestEncodeHasFileRequest(t *testing.T) {
	// Setup
	expected := []byte{0, 0, 0, 9, 3, 0, 0, 0, 0, 0, 0, 2, 80}
	req := messages.NewHasFileRequest(592)
	buf := make([]byte, len(expected))
	// Action
	n := req.Encode(buf)
	// Assert
	if n != uint32(len(expected)) {
		t.Errorf("expected HasFileRequest's encoding to write %d bytes, wrote %d", len(expected), n)
	} else if slices.Compare(expected, buf) != 0 {
		t.Errorf("expected HasFileRequest's encoding to be %v, got %v", expected, buf)
	}
}

func TestNewHasFileResponse(t *testing.T) {
	// Setup
	expected := true
	//Action
	req := messages.NewHasFileResponse(true)
	// Assert
	if req.Has != expected {
		t.Errorf("expected HasFileResponse's `has` field to be %v, got %v", expected, req.Has)
	}
}

func TestDecodeHasFileResponse(t *testing.T) {
	// Setup
	expected := false
	buf := []byte{0}
	// Action
	req := messages.DecodeHasFileResponse(buf)
	// Assert
	if req.Has != expected {
		t.Errorf("expected HasFileResponse's `has` field to be %v when decoding, got %v", expected, req.Has)
	}
}

func TestEncodeHasFileResponse(t *testing.T) {
	// Setup
	expected := []byte{0, 0, 0, 2, 3, 0}
	req := messages.NewHasFileResponse(false)
	buf := make([]byte, 6)
	// Action
	n := req.Encode(buf)
	// Assert
	if n != uint32(len(expected)) {
		t.Errorf("expected HasFileResponse's encoding to write %d bytes, wrote %d", len(expected), n)
	} else if slices.Compare(expected, buf) != 0 {
		t.Errorf("expected HasFileResponse's encoding to be %v, got %v", expected, buf)
	}
}

