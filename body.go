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

func (r *Context) WriteFullBody(status int, body []byte) error {
	r.response.StatusCode = status
	r.SetContentLength(len(body))
	_, err := r.Write(body)
	return err
}

func (r *Context) WriteFullBodyString(status int, body string) error {
	r.response.StatusCode = status
	r.SetContentLength(len(body))
	_, err := r.Write(stringToBytes(body))
	return err
}

func (r *Context) WriteStream(status int, stream io.Reader) error {
	r.response.StatusCode = status
	buf := getBuffer8k()
	defer putBuffer8k(buf)
	for {
		n, err := stream.Read((*buf)[:cap(*buf)])
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		_, err = r.Write((*buf)[:n])
		if err != nil {
			return err
		}
	}
}

func (r *Context) Flush() error {
	return r.respW.Flush()
}
