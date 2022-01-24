package silverlining

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/go-www/h1"
)

type Handler func(r *RequestContext)

type Server struct {
	Listener   net.Listener // Listener for incoming connections
	ServerName string       // Server Header (default: "SilverLining")

	MaxBodySize int64 // Max body size (default: 2MB)

	Handler Handler // Handler to invoke for each request
}

type Method = h1.Method

const (
	MethodGET     = h1.MethodGET
	MethodHEAD    = h1.MethodHEAD
	MethodPOST    = h1.MethodPOST
	MethodPUT     = h1.MethodPUT
	MethodDELETE  = h1.MethodDELETE
	MethodCONNECT = h1.MethodCONNECT
	MethodOPTIONS = h1.MethodOPTIONS
	MethodTRACE   = h1.MethodTRACE
	MethodPATCH   = h1.MethodPATCH
	MethodBREW    = h1.MethodBREW
)

type RequestContext struct {
	server *Server

	request h1.Request

	r io.Reader
	w io.Writer

	headerBuffer []byte
	next         []byte // sub slice of headerBuffer ReqeustContext.headerBuffer[headerEnd:]
}

var requestPool = sync.Pool{
	New: func() interface{} {
		return &RequestContext{
			headerBuffer: make([]byte, 0, 4096),
		}
	},
}

func GetRequestContext() *RequestContext {
	return requestPool.Get().(*RequestContext)
}

func PutRequestContext(ctx *RequestContext) {
	ctx.reset()
	requestPool.Put(ctx)
}

func (r *RequestContext) reset() {
	r.server = nil
	r.request.Reset()
	r.r = nil
	r.w = nil
	r.headerBuffer = r.headerBuffer[:0]
	r.next = nil // set to nil
}

type BodyReader struct {
	upstream io.Reader

	availableBytes int64
	readBytes      int64

	hcur int
	h    []byte
}

// BodyReader is only valid until the request is done.
// Note: This function must be called only once.
func (r *RequestContext) BodyReader() BodyReader {
	return BodyReader{
		upstream:       r.r,
		availableBytes: r.request.ContentLength,
		readBytes:      0,
		hcur:           0,
		h:              r.next,
	}
}

func (r *BodyReader) Available() int64 {
	return r.availableBytes
}

// Read reads from the BodyReader.
// It is only valid until the request is done.
// Note: This Function is not thread safe.
func (r *BodyReader) Read(p []byte) (n int, err error) {
	if r.availableBytes <= 0 {
		return 0, io.EOF
	}

	if r.availableBytes < int64(len(p)) {
		p = p[:r.availableBytes]
	}

	if r.hcur < len(r.h) {
		n = copy(p, r.h[r.hcur:])
		r.hcur += n
		r.availableBytes -= int64(n)
		return
	}

	n, err = r.upstream.Read(p)
	r.availableBytes -= int64(n)
	return
}

var ErrBodyTooLarge = fmt.Errorf("body too large")

// FullBody returns the full body of the request.
// For Big requests, Use BodyReader instead.
// Body only valid until the request is done.
// Note: This Function is not thread safe.
func (r *RequestContext) Body() ([]byte, error) {
	if r.request.ContentLength <= 0 {
		return nil, nil
	}

	if r.request.ContentLength > r.server.MaxBodySize {
		return nil, ErrBodyTooLarge
	}

	if r.next != nil && len(r.next) >= int(r.request.ContentLength) {
		return r.next[:r.request.ContentLength], nil
	}

	buffer := make([]byte, r.request.ContentLength)
	reader := r.BodyReader()
	_, err := io.ReadAtLeast(&reader, buffer, int(r.request.ContentLength))
	return buffer, err
}
