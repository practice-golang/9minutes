package wsock

import (
	"log"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

var WSockWorkers map[uint64]WebSocketChatWorker = map[uint64]WebSocketChatWorker{}
var Message chan MsgShape
var Idx uint64 = 0

func Publisher() {
	for {
		select {
		case msg := <-Message:
			for _, worker := range WSockWorkers {
				if worker.conn != nil {
					// log.Println("Send message to worker #:", worker.Idx)
					worker.msgCH <- msg
				}
			}
		}
	}
}

func WebSocketChat(w http.ResponseWriter, r *http.Request) {
	var err error
	worker := WebSocketChatWorker{}
	worker.msgCH = make(chan MsgShape)
	worker.conn, _, _, err = ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Println("ws UpgradeHTTP:", err)
	}

	idx := Idx
	worker.Idx = idx
	Idx++
	WSockWorkers[idx] = worker

	// Publish message to all workers
	go func() {
		defer func() {
			WSockWorkers[idx] = WebSocketChatWorker{}
			(worker.conn).Close()
		}()
		for {
			recv, _, err := wsutil.ReadClientData(worker.conn)
			if err != nil {
				log.Println("ws ReadClientData:", err)
				break
			}
			Message <- MsgShape{msg: string(recv)}
		}
	}()

	// Receives messages from the worker.
	go func() {
		for {
			select {
			case msg := <-worker.msgCH:
				// log.Println("#", worker.Idx, "Received:", msg)
				err = wsutil.WriteServerMessage(worker.conn, ws.OpText, []byte(msg.msg))
				if err != nil {
					log.Println("ws WriteServerMessage:", err)
				}
			}
		}
	}()
}

func InitWebSocketChat() {
	Message = make(chan MsgShape)
	// Publisher or Broadcaster
	go Publisher()
}
