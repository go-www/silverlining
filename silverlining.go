package silverlining

import (
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-www/silverlining/gopool"
	"github.com/go-www/silverlining/h1"
)

type Handler func(r *Context)

const (
	ServerStoped uint8 = iota
	ServerStarting
	ServerRunning
	ServerStopping
)

type Server struct {
	connid uint64 // Connection id

	Listener net.Listener // Listener for incoming connections

	MaxBodySize int64 // Max body size (default: 2MB)

	Handler Handler // Handler to invoke for each request

	serverStatus uint8 // Server status (stoped: 0, starting: 1, running: 2, stopping: 3)

	ReadTimeout time.Duration
}

func (s *Server) Serve(l net.Listener) error {
	s.serverStatus = ServerStarting
	s.Listener = l

	if s.MaxBodySize == 0 {
		s.MaxBodySize = 2 * 1024 * 1024
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			if s.serverStatus == ServerStopping {
				return nil
			}
			return err
		}

		gopool.Go(func() {
			s.ServeConn(conn)
		})
	}
}

var buffer8kPool sync.Pool = sync.Pool{
	New: func() interface{} {
		v := make([]byte, 8*1024)
		return &v
	},
}

func getBuffer8k() *[]byte {
	return buffer8kPool.Get().(*[]byte)
}

func putBuffer8k(b *[]byte) {
	*b = (*b)[:cap(*b)]
	buffer8kPool.Put(b)
}

func (s *Server) ServeConn(conn net.Conn) {
	var hijack bool
	readBuffer := getBuffer8k()
	reqCtx := GetRequestContext(conn)

	defer func() {
		if hijack {
			return
		}
		conn.Close()
		putBuffer8k(readBuffer)
		PutRequestContext(reqCtx)
	}()
	reqCtx.connID = atomic.AddUint64(&s.connid, 1)
	reqCtx.server = s
	reqCtx.conn = conn
	reqCtx.rawconn = conn
	reqCtx.reqR = h1.RequestReader{
		R:          conn,
		ReadBuffer: *readBuffer,
		NextBuffer: nil,
		Request:    h1.Request{},
	}

	for {
		if s.ReadTimeout > 0 {
			conn.SetReadDeadline(time.Now().Add(s.ReadTimeout))
		}
		reqCtx.resetSoft()
		_, err := reqCtx.reqR.Next()
		if err != nil {
			return
		}

		s.Handler(reqCtx)

		if reqCtx.hijack {
			hijack = true
			return
		}

		if reqCtx.respW.ContentLength == -1 {
			_ = reqCtx.respW.Flush()
			return
		}

		if reqCtx.reqR.Remaining() == 0 {
			err = reqCtx.respW.Flush()
			if err != nil {
				return
			}
		}
	}
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

type Context struct {
	server *Server

	response Response
	hwt      bool
	br       *h1.BodyReader

	respW *h1.Response
	reqR  h1.RequestReader

	conn io.ReadWriteCloser

	rawconn net.Conn
	connID  uint64

	hijack bool
}

func (r *Context) Write(p []byte) (n int, err error) {
	r.WriteHeader(r.response.StatusCode)
	return r.respW.Write(p)
}

func (r *Context) WriteString(s string) (n int, err error) {
	r.WriteHeader(r.response.StatusCode)
	return r.respW.WriteString(s)
}

func (r *Context) WriteHeader(status int) {
	if !r.hwt {
		r.response.StatusCode = status
		r.hwt = true
		r.respW.WriteHeader(r.response.StatusCode)
		err := r.writeUserHeader()
		if err != nil {
			// log.Println(err)
			return
		}
		_, err = r.respW.Write(crlf)
		if err != nil {
			// log.Println(err)
			return
		}
	}
}

var headersep = []byte(": ")
var crlf = []byte("\r\n")

func (r *Context) writeUserHeader() error {
	for i := range r.response.Headers {
		if !r.response.Headers[i].Disabled {
			_, err := r.respW.WriteString(r.response.Headers[i].Name)
			if err != nil {
				return err
			}
			_, err = r.respW.Write(headersep)
			if err != nil {
				return err
			}
			_, err = r.respW.WriteString(r.response.Headers[i].Value)
			if err != nil {
				return err
			}
			_, err = r.respW.Write(crlf)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Context) BodyReader() *h1.BodyReader {
	if r.br == nil {
		r.br = r.reqR.Body()
	}
	return r.br
}
func (r *Context) CloseBodyReader() {
	if r.br != nil {
		h1.PutBodyReader(r.br)
		r.br = nil
	}
}

type Response struct {
	Headers []Header

	StatusCode int
}

type Header struct {
	Disabled bool

	Name  string
	Value string
}

func (r *Response) reset() {
	r.StatusCode = 200
	r.Headers = r.Headers[:0]
}

func (r *Context) RawURI() []byte {
	return r.reqR.Request.RawURI
}

func (r *Context) Path() []byte {
	return r.reqR.Request.URI.Path()
}

func (r *Context) QueryParams() []h1.Query {
	return r.reqR.Request.URI.Query()
}

func (r *Context) GetQueryParam(name []byte) ([]byte, error) {
	return r.reqR.Request.URI.QueryValue(name)
}

func (r *Context) GetQueryParamString(name string) (string, error) {
	v, err := r.reqR.Request.URI.QueryValue(stringToBytes(name))
	if err != nil {
		return "", err
	}
	return bytesToString(v), nil
}

func (r *Context) SetContentLength(length int) {
	r.respW.ContentLength = length
}

var requestPool = sync.Pool{
	New: func() interface{} {
		v := new(Context)
		return v
	},
}

func GetRequestContext(upstream io.Writer) *Context {
	ctx := requestPool.Get().(*Context)
	ctx.respW = h1.GetResponse(upstream)
	return ctx
}

func PutRequestContext(ctx *Context) {
	ctx.resetHard()
	requestPool.Put(ctx)
}

func (r *Context) resetSoft() {
	r.hwt = false
	r.CloseBodyReader()
	r.response.reset()
}

func (r *Context) resetHard() {
	r.resetSoft()
	r.conn = nil
	r.rawconn = nil
	r.server = nil
	r.hijack = false
	r.reqR.Reset()
	h1.PutResponse(r.respW)
}

func (r *Context) HijackConn() (bufR h1.HijackReader, bufW *h1.Response, conn net.Conn) {
	bufR = r.reqR.Hijack()
	bufW = r.respW
	conn = r.rawconn
	r.hijack = true

	return
}

func (r *Context) Method() h1.Method {
	return r.reqR.Request.Method
}

func (r *Context) ConnectionClose() {
	r.SetContentLength(-1)
}

func (r *Context) KillConn() error {
	r.ConnectionClose()
	return r.conn.Close()
}

// (*Context).ConnID returns the connection ID.
func (r *Context) ConnID() uint64 {
	return r.connID
}

// (*Context).RemoteAddr returns the remote address.
func (r *Context) RemoteAddr() net.Addr {
	return r.rawconn.RemoteAddr()
}
