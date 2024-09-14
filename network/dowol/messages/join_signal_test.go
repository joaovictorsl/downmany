package messages_test

import (
	"slices"
	"testing"

	"github.com/joaovictorsl/downmany/network/dowol/messages"
)

func TestEncodeJoinSignal(t *testing.T) {
	// Setup
	expected := []byte{0, 0, 0, 1, 0}
	buf := make([]byte, len(expected))
	// Action
	n := messages.NewJoinSignal().Encode(buf)
	// Assert
	if n != uint32(len(expected)) {
		t.Errorf("expected JoinSignal's encoding to write %d bytes, wrote %d", len(expected), n)
	} else if slices.Compare(expected, buf) != 0 {
		t.Errorf("expected JoinSignal's encoding to be %v, got %v", expected, buf)
	}
}
