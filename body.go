package silverlining

import (
	"errors"
	"io"
)

var ErrBodyTooLarge = errors.New("body too large")

func (r *Context) Body() ([]byte, error) {
	if r.reqR.Request.ContentLength > r.server.MaxBodySize {
		return nil, ErrBodyTooLarge
	}

	var buf []byte = make([]byte, r.reqR.Request.ContentLength)
	n, err := io.ReadAtLeast(r.BodyReader(), buf, len(buf))
	r.CloseBodyReader()
	return buf[:n], err
}
