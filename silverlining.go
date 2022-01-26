package silverlining

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/go-www/h1"
)

type Handler func(r *RequestContext)

const (
	ServerStoped uint8 = iota
	ServerStarting
	ServerRunning
	ServerStopping
)

type Server struct {
	Listener   net.Listener // Listener for incoming connections
	ServerName string       // Server Header (default: "SilverLining")

	MaxBodySize int64 // Max body size (default: 2MB)

	Handler Handler // Handler to invoke for each request

	serverTime *[]byte // Server time RFC1123

	serverStatus uint8 // Server status (stoped: 0, starting: 1, running: 2, stopping: 3)
}

var serverHeaderBytes []byte = []byte("Server: ")

func (s *Server) serverTimeWorker() {
	var ts [16][]byte
	for i := 0; i < len(ts); i++ {
		ts[i] = make([]byte, 0, (len("Date: Mon, 02 Jan 2006 15:04:05 MST\r\nServer: ")+len(s.ServerName))*2)
		ts[i] = ts[i][:0]
		ts[i] = time.Now().UTC().AppendFormat(ts[i], "Date: Mon, 02 Jan 2006 15:04:05 MST\r\n")
	}

	var i int

	for {
		i = (i + 1) & 15

		ts[i] = ts[i][:0]
		ts[i] = time.Now().UTC().AppendFormat(ts[i], "Date: Mon, 02 Jan 2006 15:04:05 MST\r\n")
		ts[i] = append(ts[i], serverHeaderBytes...)
		ts[i] = append(ts[i], s.ServerName...)
		ts[i] = append(ts[i], '\r', '\n')
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&s.serverTime)), unsafe.Pointer(&ts[i]))
		time.Sleep(time.Second * 5) // Migrate Cache Timing Attack

		if s.serverStatus == ServerStopping {
			return
		}
	}
}

func (s *Server) Serve(l net.Listener) error {
	s.serverStatus = ServerStarting
	s.Listener = l

	go s.serverTimeWorker()

	for {
		conn, err := l.Accept()
		if err != nil {
			if s.serverStatus == ServerStopping {
				return nil
			}
			return err
		}

		go s.ServeConn(conn)
	}
}

type dummyReadWriter struct{}

func (d dummyReadWriter) Write(p []byte) (n int, err error) { return }
func (d dummyReadWriter) Read(p []byte) (n int, err error)  { return }

var drw = dummyReadWriter{}
var BufWriterPool sync.Pool = sync.Pool{
	New: func() interface{} {
		return bufio.NewWriter(drw)
	},
}

func (s *Server) ServeConn(conn net.Conn) {
	var err error

	defer conn.Close()
	ctx := GetRequestContext()
	defer PutRequestContext(ctx)

	ctx.r.Reset(conn)
	ctx.w.Reset(conn)
	ctx.conn = conn

	fill := func() error {
		var n int
		n = copy(ctx.headerBuffer, ctx.next)
		ctx.next = ctx.headerBuffer[:n]
		n, err = ctx.r.Read(ctx.headerBuffer[n:])
		if err != nil {
			return err
		}
		ctx.next = ctx.headerBuffer[:n]
		return nil
	}

	for {
		err = func() error {
			var err error

			ctx.next, err = h1.ParseRequestLine(&ctx.request, ctx.next)
			if err != nil {
				if err == h1.ErrBufferTooSmall {
					err = fill()
					if err != nil {
						return err
					}
					ctx.next, err = h1.ParseRequestLine(&ctx.request, ctx.next)
					if err != nil {
						return err
					}
				}
				return err
			}

			ctx.next, err = h1.ParseHeaders(&ctx.request, ctx.next)
			if err != nil {
				if err == h1.ErrBufferTooSmall {
					err = fill()
					if err != nil {
						return err
					}
					ctx.next, err = h1.ParseRequestLine(&ctx.request, ctx.next)
					if err != nil {
						return err
					}
					ctx.next, err = h1.ParseHeaders(&ctx.request, ctx.next)
					if err != nil {
						return err
					}
				}
				return err
			}

			ctx.response.reset()
			s.Handler(ctx)

			return nil
		}()
		if err != nil {
			log.Printf("[ERROR] %s", err)
			return
		}

		if ctx.r.Buffered() <= 0 {
			err = ctx.w.Flush()
			if err != nil {
				log.Printf("[ERROR] %s", err)
				return
			}
		}
	}
}

func (r *RequestContext) Write(p []byte) (n int, err error) {
	return r.w.Write(p)
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

	request  h1.Request
	response Response

	r *bufio.Reader
	w *bufio.Writer

	conn net.Conn

	headerBuffer []byte
	next         []byte // sub slice of headerBuffer ReqeustContext.headerBuffer[headerEnd:]
}

type Response struct {
	StatusCode int

	ContentLength      int
	ContentLengthBytes []byte

	Headers []Header
}

type Header struct {
	Disabled bool

	Name  string
	Value string
}

func (resp *Response) reset() {
	resp.StatusCode = 0
	resp.ContentLength = -1
	resp.ContentLengthBytes = resp.ContentLengthBytes[:0]
	resp.Headers = resp.Headers[:0]
}

func (r *RequestContext) SetHeader(name, value string) {
	for i := range r.response.Headers {
		if r.response.Headers[i].Name == name {
			r.response.Headers[i].Disabled = false
			r.response.Headers[i].Value = value
			return
		}
	}
	r.response.Headers = append(r.response.Headers, Header{
		Disabled: false,
		Name:     name,
		Value:    value,
	})
}

func (r *RequestContext) DeleteHeader(name string) {
	for i := range r.response.Headers {
		if r.response.Headers[i].Name == name {
			r.response.Headers[i].Disabled = true
			return
		}
	}
}

var requestPool = sync.Pool{
	New: func() interface{} {
		return &RequestContext{
			headerBuffer: make([]byte, 0, 4096),
			r:            bufio.NewReader(drw),
			w:            bufio.NewWriter(drw),
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
	r.response.reset()
	r.r.Reset(drw)
	r.w.Reset(drw)
	r.conn = nil
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

var contentLengthHeader = []byte("Content-Length: ")
var crlf = []byte("\r\n")

func (r *RequestContext) SetContentLength(length int) {
	r.request.ContentLength = int64(length)
}

func (r *RequestContext) WriteHeader(status int) error {
	_, err := r.w.Write(GetStatusLine(status))
	if err != nil {
		return err
	}
	_, err = r.w.Write(*r.server.serverTime)
	if err != nil {
		return err
	}

	// Write Content-Length
	if r.request.ContentLength >= 0 {
		_, err = r.w.Write(contentLengthHeader)
		if err != nil {
			return err
		}
		r.response.ContentLengthBytes = strconv.AppendInt(r.response.ContentLengthBytes, r.request.ContentLength, 10)
		r.response.ContentLengthBytes = append(r.response.ContentLengthBytes, crlf...)
		_, err = r.w.Write(r.response.ContentLengthBytes)
		if err != nil {
			return err
		}
	}

	for _, header := range r.response.Headers {
		if header.Disabled {
			continue
		}
		_, err = r.w.WriteString(header.Name)
		if err != nil {
			return err
		}
		_, err = r.w.WriteString(": ")
		if err != nil {
			return err
		}
		_, err = r.w.WriteString(header.Value)
		if err != nil {
			return err
		}
		_, err = r.w.WriteString("\r\n")
		if err != nil {
			return err
		}
	}
	return nil
}
