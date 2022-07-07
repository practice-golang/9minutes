package wsock

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/websocket"
)

func TestWebSocketEcho(t *testing.T) {
	serverHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WebSocketEcho(w, r)
	})
	s := httptest.NewServer(serverHandler)
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	w, err := websocket.Dial(u, "", s.URL)
	if err != nil {
		t.Errorf("Dial error %v", err)
	}

	msg := []byte("Hello")

	i, err := w.Write(msg)
	if err != nil {
		t.Errorf("Write error %v", err)
	}

	require.Equal(t, len(msg), i, "WebsocketEcho not equal")
	log.Println(i)
}
