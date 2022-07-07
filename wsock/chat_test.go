package wsock

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/websocket"
)

func Test_WebSocketChat(t *testing.T) {
	t.Run("WebsocketChat", func(t *testing.T) {
		InitWebSocketChat()

		serverHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			WebSocketChat(w, r)
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

		require.Equal(t, len(msg), i, "WebsocketChat not equal")
	})
}
