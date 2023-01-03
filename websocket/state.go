package websocket

const (
	Connecting = 0
	Open       = 1
	Closing    = 2
	Closed     = 3
)

type Instance struct {
	binaryType uint8
	buffered   []byte
	extensions string
	protocol   string
	readyState uint8
	url        string
}

type IWebsocket interface {
	open(event Instance)
	message()
	error()
	close()
}
