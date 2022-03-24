package h1

import (
	"io"
	"sync"
)

type RequestReader struct {
	R io.Reader

	ReadBuffer []byte
	NextBuffer []byte

	Request Request
}

func (r *RequestReader) Reset() {
	r.ReadBuffer = r.ReadBuffer[:cap(r.ReadBuffer)]
	r.NextBuffer = r.ReadBuffer[:0]
	r.Request.Reset()
}

func (r *RequestReader) Fill() (n int, err error) {
	// Copy the remaining bytes to the read buffer
	n0 := copy(r.ReadBuffer[:cap(r.ReadBuffer)], r.NextBuffer)

	// Read more bytes
	n1, err := r.R.Read(r.ReadBuffer[n0:cap(r.ReadBuffer)])
	if err != nil {
		return 0, err
	}

	// Set the next buffer to the read buffer
	r.NextBuffer = r.ReadBuffer[:n0+n1]

	return n1, nil
}

func (r *RequestReader) Next() (remaining int, err error) {
	var retryCount int = 0

	if r.Remaining() == 0 {
		n, err := r.R.Read(r.ReadBuffer[:cap(r.ReadBuffer)])
		if err != nil {
			return 0, err
		}
		r.NextBuffer = r.ReadBuffer[:n]
	}

parse:
	// Reset the request
	r.Request.Reset()

	// Read request line
	r.NextBuffer, err = ParseRequestLine(&r.Request, r.NextBuffer)
	if err != nil {
		if err == ErrBufferTooSmall {
			// Buffer is too small, read more bytes

			_, err = r.Fill()
			if err != nil {
				return 0, err
			}

			if len(r.NextBuffer) == cap(r.ReadBuffer) {
				return 0, ErrRequestHeaderTooLarge
			}

			// Retry parsing
			retryCount++
			if retryCount > 1 {
				return 0, ErrBufferTooSmall
			}
			goto parse
		}
	}

	// Read headers
	r.NextBuffer, err = ParseHeaders(&r.Request, r.NextBuffer)
	if err != nil {
		if err == ErrBufferTooSmall {
			// Buffer is too small, read more bytes

			_, err = r.Fill()
			if err != nil {
				return 0, err
			}

			if len(r.NextBuffer) == cap(r.ReadBuffer) {
				return 0, ErrRequestHeaderTooLarge
			}

			// Retry parsing
			retryCount++
			if retryCount > 1 {
				return 0, ErrBufferTooSmall
			}
			goto parse
		}
	}

	// Parse URI
	r.Request.URI.Parse(r.Request.RawURI)

	return len(r.NextBuffer), nil
}

func (r *RequestReader) Remaining() int {
	return len(r.NextBuffer)
}

func (r *RequestReader) Body() *BodyReader {
	br := GetBodyReader()
	br.Limit = int(r.Request.ContentLength)
	br.Upstream = r
	return br
}

type BodyReader struct {
	Upstream *RequestReader

	Limit int
	Index int
}

func (r *BodyReader) reset() {
	r.Upstream = nil

	r.Limit = 0
	r.Index = 0
}

var BodyReaderPool = &sync.Pool{
	New: func() any {
		return &BodyReader{}
	},
}

func GetBodyReader() *BodyReader {
	return BodyReaderPool.Get().(*BodyReader)
}

func PutBodyReader(r *BodyReader) {
	r.reset()
	BodyReaderPool.Put(r)
}

func (r *BodyReader) Read(p []byte) (n int, err error) {
	if r.Index >= r.Limit {
		return 0, io.EOF
	}

	// If p is bigger than the remaining bytes, set the limit to the remaining bytes
	if len(p) > r.Limit-r.Index {
		p = p[:r.Limit-r.Index]
	}

	// Copy the remaining bytes to the read buffer
	n0 := copy(p, r.Upstream.NextBuffer)
	r.Upstream.NextBuffer = r.Upstream.NextBuffer[n0:]
	r.Index += n0

	if len(p) == n0 {
		// No more bytes to read
		return n0, nil
	}

	// Fill The buffer
	_, err = r.Upstream.Fill()
	if err != nil {
		return n0, err
	}

	// If p is bigger than Uptream.ReadBuffer then read directly from the upstream
	if len(p)-n0 > len(r.Upstream.ReadBuffer) {
		n1, err := r.Upstream.R.Read(p[n0:])
		if err != nil {
			return n0 + n1, err
		}
		return n0 + n1, nil
	}

	// Read more bytes
	n1, err := r.Read(p[n0:])
	if err != nil {
		return n0 + n1, err
	}

	return n0 + n1, nil
}

func (r *BodyReader) Close() error {
	PutBodyReader(r)
	return nil
}

type HijackReader struct {
	Upstream *RequestReader
}

func (h HijackReader) Read(p []byte) (n int, err error) {
	// Copy the remaining bytes to the read buffer
	n0 := copy(p, h.Upstream.NextBuffer)
	h.Upstream.NextBuffer = h.Upstream.NextBuffer[n0:]

	if len(p) == n0 {
		// No more bytes to read
		return n0, nil
	}

	// Fill The buffer
	_, err = h.Upstream.Fill()
	if err != nil {
		return n0, err
	}

	// If p is bigger than Uptream.ReadBuffer then read directly from the upstream
	if len(p)-n0 > len(h.Upstream.ReadBuffer) {
		n1, err := h.Upstream.R.Read(p[n0:])
		if err != nil {
			return n0 + n1, err
		}
		return n0 + n1, nil
	}

	// Read more bytes
	n1, err := h.Read(p[n0:])
	if err != nil {
		return n0 + n1, err
	}

	return n0 + n1, nil
}

func (r *RequestReader) Hijack() HijackReader {
	return HijackReader{
		Upstream: r,
	}
}
