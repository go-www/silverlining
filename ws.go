package silverlining

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gobwas/ws"
)

var ErrUpgradeBadRequest = errors.New("upgrade request is not valid")

func (r *Context) writeUpgradeWebSocket() error {
	Upgrade, ok := r.RequestHeaders().Get("Upgrade")
	if !ok {
		r.ResponseHeaders().Set("Content-Type", "text/plain")
		r.WriteFullBodyString(400, "Missing Upgrade header")
		return ErrUpgradeBadRequest
	}
	if Upgrade != "websocket" && strings.EqualFold(Upgrade, "websocket") {
		r.ResponseHeaders().Set("Content-Type", "text/plain")
		r.WriteFullBodyString(400, "Invalid Upgrade header")
		return ErrUpgradeBadRequest
	}

	Connection, ok := r.RequestHeaders().Get("Connection")
	if !ok {
		r.ResponseHeaders().Set("Content-Type", "text/plain")
		r.WriteFullBodyString(400, "Missing Connection header")
		return ErrUpgradeBadRequest
	}
	if !strings.EqualFold(Connection, "Upgrade") {
		r.ResponseHeaders().Set("Content-Type", "text/plain")
		r.WriteFullBodyString(400, "Sorry, only \"Upgrade\" connection is supported")
		return ErrUpgradeBadRequest
	}

	SecVersion, ok := r.RequestHeaders().Get("Sec-WebSocket-Version")
	if !ok {
		r.ResponseHeaders().Set("Content-Type", "text/plain")
		r.WriteFullBodyString(400, "Missing Sec-WebSocket-Version header")
		return ErrUpgradeBadRequest
	}
	if SecVersion != "13" {
		r.ResponseHeaders().Set("Content-Type", "text/plain")
		r.WriteFullBodyString(400, "Sorry, only WebSocket version 13 is supported.")
		return ErrUpgradeBadRequest
	}

	SecKey, ok := r.RequestHeaders().Get("Sec-WebSocket-Key")
	if !ok {
		r.ResponseHeaders().Set("Content-Type", "text/plain")
		r.WriteFullBodyString(400, "Missing Sec-WebSocket-Key header")
		return ErrUpgradeBadRequest
	}

	if len(SecKey) < 24 {
		r.ResponseHeaders().Set("Content-Type", "text/plain")
		r.WriteFullBodyString(400, "Invalid Sec-WebSocket-Key header")
		return ErrUpgradeBadRequest
	}

	r.ResponseHeaders().Set("Connection", "Upgrade")
	r.ResponseHeaders().Set("Upgrade", "Websocket")

	h := sha1.New()
	h.Write(stringToBytes(SecKey))
	h.Write(stringToBytes("258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	r.ResponseHeaders().Set("Sec-WebSocket-Accept", base64.StdEncoding.EncodeToString(h.Sum(nil)))
	r.WriteHeader(http.StatusSwitchingProtocols)
	err := r.Flush()
	if err != nil {
		return err
	}

	return nil
}

// Upgrade the request to WebSocket protocol.
// param:
// 		op: the opcode to use for the WebSocket connection. (ws.OpText|ws.OpBinary)
func (r *Context) UpgradeWebSocket(op ws.OpCode) (conn io.ReadWriteCloser, err error) {
	err = r.writeUpgradeWebSocket()
	if err != nil {
		return nil, err
	}
	conn = r.HijackConn()
	return
}
