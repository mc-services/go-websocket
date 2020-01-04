package chat

type Room struct {
	Clients []Client
	Broadcast chan Msg `json:"-"`
}

func (Room) broadcast() {

}
