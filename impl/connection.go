package impl

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type Connection struct {
	wsConn *websocket.Conn
	inChan chan []byte
	outChan chan []byte
	closeChan chan []byte
	sync.Mutex	// lock
	isClosed bool
}

func InitConnection(wsConn *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConn: wsConn,
		inChan: make(chan []byte, 1000),
		outChan: make(chan []byte, 1000),
		closeChan: make(chan []byte, 1),
	}

	// start read goroutine
	go conn.readLoop()

	// start write goroutine
	go conn.writeLoop()

	return
}

func (conn *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <- conn.inChan:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}

	return 
}

func (conn *Connection) WriteMessage(data []byte) (err error) {
	select {
	case conn.outChan <- data:
	case <- conn.closeChan:
		err = errors.New("connect is closed")
	}

	return
}

func (conn *Connection) Close() {
	// 线程安全的，可重入的API

	// close connection
	conn.wsConn.Close()

	conn.Mutex.Lock()
	if ! conn.isClosed {
		// close chan once
		close(conn.closeChan)
		conn.isClosed = true
	}
	conn.Mutex.Unlock()
}

// 内部实现
func (conn *Connection) readLoop()  {
	var (
		data []byte
		err error
	)

	for {
		if _, data, err = conn.wsConn.ReadMessage(); err != nil {
			goto ERR
		}

		select {
		case conn.inChan <- data:
		case <-conn.closeChan:	// when close chan closeChan
		goto ERR
		}
		
	}

	ERR:
		conn.Close()
}

func (conn *Connection) writeLoop() {
	var (
		data []byte
		err error
	)

	for {
		select {
		case data = <- conn.outChan:
		case <-conn.closeChan:
			goto ERR
		}
		
		if err = conn.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}

	ERR:
		conn.Close()
}