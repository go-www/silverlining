package h1

/*
type Connection uint8

const (
	ConnectionUnset Connection = iota
	ConnectionClose
	ConnectionKeepAlive
	ConnectionUpgrade
)

var connectionHeaderTable [16][]byte

var _ = func() int {
	nilString := []byte("")

	for i := range connectionHeaderTable {
		connectionHeaderTable[i] = nilString
	}

	connectionHeaderTable[ConnectionClose] = []byte("Connection: close\r\n")
	connectionHeaderTable[ConnectionKeepAlive] = []byte("Connection: keep-alive\r\n")
	connectionHeaderTable[ConnectionUpgrade] = []byte("Connection: upgrade\r\n")

	return 0
}

func getConnectionHeader(c Connection) []byte {
	return connectionHeaderTable[c%16]
}
*/
