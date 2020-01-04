package chat

const (
	MSG_TYPE_MSG string = "msg"
	MSG_TYPE_CONNECTED string = "connected"
	MSG_TYPE_ERROR string = "error"
	MSG_TYPE_CLOSED string = "closed"
)

type Msg struct {
	Type string `json:"type"` // 类型
	Msg  string `json:"msg"`  // 消息文本
	From Client `json:"from"` // 来源
}
