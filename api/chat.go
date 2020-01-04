package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gp-websoket/chat"
	"gp-websoket/impl"
	"gp-websoket/model"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		//  allow cors
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// 房间
	rooms = make(map[string]chat.Room)
)

func Rooms(c *gin.Context) {
	//claims := c.MustGet("auth_user")
	var res = []struct {
		Name  string `json:"name"`
		Count uint `json:"count"`
		HasMe bool `json:"has_me"`
	}{}

	for k := range rooms {
		var clients = rooms[k].Clients
		var user = c.MustGet("auth_user").(model.User)
		var HasMe bool = false

		for _, client := range clients {
			if client.User.Id == user.Id {
				HasMe = true
			}
		}

		res = append(res, struct {
			Name  string `json:"name"`
			Count uint `json:"count"`
			HasMe bool `json:"has_me"`
		}{Name: k, Count: uint(len(clients)), HasMe: HasMe})
	}

	c.JSON(http.StatusOK, res)
	c.JSON(http.StatusOK, rooms)
}

func CreateRoom(c *gin.Context) {
	var form struct {
		Name string `form:"name"`
	}

	if c.ShouldBind(&form) == nil {
		fmt.Println("name")
		fmt.Println(form.Name)
		fmt.Println(rooms)
		if _, ok := rooms[form.Name]; ok == true {
			c.JSON(http.StatusUnprocessableEntity, "房间已存在")
			return
		}

		user := c.MustGet("auth_user").(model.User)

		fmt.Println(23333)
		fmt.Println(user)

		rooms[form.Name] = chat.Room{
			Clients: []chat.Client{
				{
					User: user,
					Cone: impl.Connection{},
				},
			},
		}

		c.JSON(http.StatusOK, "添加成功")
	}
}

func JoinRoom(c *gin.Context) {
	user := c.MustGet("auth_user").(model.User)
	name := c.Param("name")

	room, ok := rooms[name]
	if ok == false {
		c.JSON(http.StatusUnprocessableEntity, "房间不存在")
		return
	}

	for _, client := range room.Clients {
		if client.User.Id == user.Id {
			c.JSON(http.StatusUnprocessableEntity, "已存在与此房间")
			return
		}
	}

	room.Clients = append(room.Clients, chat.Client{
		User: user,
	})

	c.JSON(http.StatusOK, room)
}

func WsHandler(c *gin.Context) {
	var (
		wsConn *websocket.Conn
		conn   *impl.Connection
		data   []byte
		err    error
	)

	// 🤝
	if wsConn, err = upgrader.Upgrade(c.Writer, c.Request, nil); err != nil {
		return
	}

	// 初始化连接，开启读协程与写协程
	if conn, err = impl.InitConnection(wsConn); err != nil {
		conn.Close()
		return
	}

	// heartbeat
	go func() {
		var err error
		for {
			msg, _ := json.Marshal(chat.Msg{Type: chat.MSG_TYPE_MSG, Msg: "heartbeat"})
			if err = conn.WriteMessage(msg); err != nil {
				return
			}

			// heartbeat
			time.Sleep(1 * time.Second)
		}
	}()

	// 不断读取消息并写入消息
	go func() {
		for {
			if data, err = conn.ReadMessage(); err != nil {
				conn.Close()
				return
			}

			if err = conn.WriteMessage(data); err != nil {
				conn.Close()
				return
			}
		}
	}()
}
