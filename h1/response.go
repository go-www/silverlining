package h1

import (
	"io"
	"strconv"
	"sync"
)

var ResponsePool = sync.Pool{
	New: func() any {
		return &Response{
			upstream:      nil,
			buf:           make([]byte, 8192),
			itoaBuf:       make([]byte, 0, 32),
			n:             0,
			ContentLength: -1,
			//Connection:    ConnectionKeepAlive,
		}
	},
}

func GetResponse(upstream io.Writer) *Response {
	r := ResponsePool.Get().(*Response)
	r.upstream = upstream
	return r
}

func PutResponse(r *Response) {
	r.Reset()
	r.upstream = nil
	ResponsePool.Put(r)
}

type Response struct {
	upstream io.Writer

	buf []byte // Note: Do not use append() to add bytes to this buffer. Use Write() instead. This is to avoid unnecessary memory allocations.
	n   int

	// Itoa Buffer
	itoaBuf []byte // buffer for itoa

	// Standard Hop-by-Hop response headers.
	ContentLength int
	//Connection    Connection
}

func (r *Response) Reset() {
	r.n = 0
	r.ContentLength = -1
}

var DefaultFastDateServer = NewFastDateServer("h1")

var _ = func() int {
	go DefaultFastDateServer.Start()
	return 0
}()

var DateServerHeaderFunc = func() []byte {
	return DefaultFastDateServer.GetDate()
}

func (r *Response) Flush() error {
	if r.upstream == nil || r.n == 0 {
		return nil
	}

	_, err := r.upstream.Write(r.buf[:r.n])
	if err != nil {
		return err
	}

	r.Reset()
	return nil
}

func (r *Response) Write(b []byte) (int, error) {
	// Check if buffer is full
	if len(r.buf) == r.n {
		if err := r.Flush(); err != nil {
			return 0, err
		}
	}

	var n int
	for len(b) > 0 {
		n = copy(r.buf[r.n:], b)
		r.n += n
		b = b[n:]

		if r.n == len(r.buf) {
			if err := r.Flush(); err != nil {
				return 0, err
			}
		}

		if len(b) > len(r.buf) {
			// direct write
			_, err := r.upstream.Write(b)
			return len(b), err
		}
	}

	return n, nil
}

func (r *Response) WriteString(b string) (int, error) {
	return r.Write(stringToBytes(b))
}

func (r *Response) WriteInt(i int) (int, error) {
	r.itoaBuf = r.itoaBuf[:0]
	r.itoaBuf = strconv.AppendInt(r.itoaBuf, int64(i), 10)
	return r.Write(r.itoaBuf)
}

func (r *Response) WriteUint(u uint) (int, error) {
	r.itoaBuf = r.itoaBuf[:0]
	r.itoaBuf = strconv.AppendUint(r.itoaBuf, uint64(u), 10)
	return r.Write(r.itoaBuf)
}

func (r *Response) WriteInt64(i int64) (int, error) {
	r.itoaBuf = r.itoaBuf[:0]
	r.itoaBuf = strconv.AppendInt(r.itoaBuf, i, 10)
	return r.Write(r.itoaBuf)
}

func (r *Response) WriteUint64(u uint64) (int, error) {
	r.itoaBuf = r.itoaBuf[:0]
	r.itoaBuf = strconv.AppendUint(r.itoaBuf, u, 10)
	return r.Write(r.itoaBuf)
}

func (r *Response) WriteUint64Hex(u uint64) (int, error) {
	r.itoaBuf = r.itoaBuf[:0]
	r.itoaBuf = strconv.AppendUint(r.itoaBuf, u, 16)
	return r.Write(r.itoaBuf)
}

func (r *Response) WriteStatusLine(status int) error {
	_, err := r.Write(GetStatusLine(status))
	return err
}

var contentLengthHeader = []byte("Content-Length: ")
var crlf = []byte("\r\n")

func (r *Response) WriteHeader(status int) error {
	err := r.WriteStatusLine(status)
	if err != nil {
		return err
	}
	// Write standard hop-by-hop response headers

	_, err = r.Write(DateServerHeaderFunc())
	if err != nil {
		return err
	}

	// Content-Length
	if r.ContentLength >= 0 {
		_, err = r.Write(contentLengthHeader)
		if err != nil {
			return err
		}
		_, err = r.WriteInt(r.ContentLength)
		if err != nil {
			return err
		}
		_, err = r.Write(crlf)
		if err != nil {
			return err
		}
	}

	return nil
}
