package main

import (
	"github.com/gorilla/websocket"
	"gp-websoket/impl"
	"net/http"
	"time"
)

const (
	host = "0.0.0.0:8888"
)

var (
	upgrader = websocket.Upgrader{
		//  allow cors
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {
	http.HandleFunc("/ws", wsHandler)

	http.ListenAndServe(host, nil)
}

func wsHandler(writer http.ResponseWriter, request *http.Request) {
	var (
		wsConn *websocket.Conn
		conn *impl.Connection
		data []byte
		err error
	)

	// 🤝
	if wsConn, err = upgrader.Upgrade(writer, request, nil); err != nil {
		return
	}

	if conn, err = impl.InitConnection(wsConn); err != nil {
		goto ERR
	}

	// heartbeat test connection
	go func() {
		var err error
		for {
			if err = conn.WriteMessage([]byte("heartbeat")); err != nil {
				return
			}

			// heartbeat
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		if data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}

		if err = conn.WriteMessage(data); err !=nil {
			goto ERR
		}
	}

	ERR:
		conn.Close()
}

