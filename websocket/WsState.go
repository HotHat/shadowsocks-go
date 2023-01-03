package websocket

const (
	Connection = 0
	Open       = 1
	Closing    = 2
	Closed     = 3
)

type State uint8

type WsEvent interface {
	OnOpen()
	OnClose()
	OnError()
	OnMessage(data []byte)
	Send(data []byte)
}

type WsState struct {
	url        string
	readyState State
	protocol   string
	extension  string
}

func NewWebSocket(url, protocol string) *WsState {
	return &WsState{
		url:        url,
		readyState: Connection,
		protocol:   protocol,
		extension:  "",
	}
}
