package silverlining

import (
	"io"
	"log"
	"net"
	"sync"

	"github.com/go-www/h1"
	"github.com/go-www/silverlining/gopool"
)

type Handler func(r *RequestContext)

const (
	ServerStoped uint8 = iota
	ServerStarting
	ServerRunning
	ServerStopping
)

type Server struct {
	Listener net.Listener // Listener for incoming connections

	MaxBodySize int64 // Max body size (default: 2MB)

	Handler Handler // Handler to invoke for each request

	serverStatus uint8 // Server status (stoped: 0, starting: 1, running: 2, stopping: 3)
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

type logConn struct {
	c net.Conn
}

func (lc *logConn) Read(b []byte) (n int, err error) {
	n, err = lc.c.Read(b)
	log.Printf("Read: %d bytes\n", n)
	log.Println(string(b[:n]))
	return
}

func (lc *logConn) Write(b []byte) (n int, err error) {
	log.Printf("Write: %d bytes\n", len(b))
	log.Println(string(b))
	return lc.c.Write(b)
}

func (s *Server) ServeConn(conn net.Conn) {
	defer conn.Close()

	readBuffer := getBuffer8k()
	defer putBuffer8k(readBuffer)

	//reqCtx := GetRequestContext(&logConn{conn})
	reqCtx := GetRequestContext(conn)
	defer PutRequestContext(reqCtx)
	reqCtx.server = s
	//reqCtx.conn = &logConn{conn}
	reqCtx.conn = conn
	reqCtx.reqR = h1.RequestReader{
		//R: &logConn{conn},
		R:          conn,
		ReadBuffer: *readBuffer,
		NextBuffer: nil,
		Request:    h1.Request{},
	}

	for {
		_, err := reqCtx.reqR.Next()
		if err != nil {
			if err == io.EOF {
				log.Println("EOF")
				return
			}
			log.Println(err)
			return
		}
		reqCtx.resetSoft()

		//println("Request:", reqCtx.reqR.Request.Method, string(reqCtx.reqR.Request.URI))

		s.Handler(reqCtx)

		//println("Response:", reqCtx.response.StatusCode)

		if reqCtx.reqR.Remaining() == 0 {
			err = reqCtx.respW.Flush()
			if err != nil {
				log.Println(err)
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

type RequestContext struct {
	server *Server

	response Response
	hwt      bool
	br       *h1.BodyReader

	respW *h1.Response
	reqR  h1.RequestReader

	conn io.ReadWriter
}

func (rctx *RequestContext) Write(p []byte) (n int, err error) {
	rctx.WriteHeader(rctx.response.StatusCode)
	return rctx.respW.Write(p)
}

func (rctx *RequestContext) WriteHeader(status int) {
	if !rctx.hwt {
		rctx.response.StatusCode = status
		rctx.hwt = true
		rctx.respW.WriteHeader(rctx.response.StatusCode)
		err := rctx.writeUserHeader()
		if err != nil {
			log.Println(err)
		}
		_, err = rctx.respW.Write(crlf)
		if err != nil {
			log.Println(err)
		}
	}
}

var headersep = []byte(": ")
var crlf = []byte("\r\n")

func (rctx *RequestContext) writeUserHeader() error {
	for i := range rctx.response.Headers {
		if !rctx.response.Headers[i].Disabled {
			_, err := rctx.respW.WriteString(rctx.response.Headers[i].Name)
			if err != nil {
				return err
			}
			_, err = rctx.respW.Write(headersep)
			if err != nil {
				return err
			}
			_, err = rctx.respW.WriteString(rctx.response.Headers[i].Value)
			if err != nil {
				return err
			}
			_, err = rctx.respW.Write(crlf)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (rctx *RequestContext) BodyReader() *h1.BodyReader {
	if rctx.br == nil {
		rctx.br = rctx.reqR.Body()
	}
	return rctx.br
}
func (rctx *RequestContext) CloseBodyReader() {
	if rctx.br != nil {
		h1.PutBodyReader(rctx.br)
		rctx.br = nil
	}
}

type ResponseType uint8

const (
	ResponseTypeNone ResponseType = iota
	ResponseTypeFullBody
	ResponseTypeStream
	ResponseTypeHijack
	ResponseTypeUser
)

type Response struct {
	Headers []Header

	// FullBody []byte    // for ResponseTypeFullBody
	// Stream   io.Reader // for ResponseTypeStream
	// Hijack   func() (io.ReadCloser, io.Writer)

	StatusCode int

	// BodyType ResponseType
}

type Header struct {
	Disabled bool

	Name  string
	Value string
}

func (resp *Response) reset() {
	resp.StatusCode = 200
	resp.Headers = resp.Headers[:0]
}

func (resp *RequestContext) RawURI() []byte {
	return resp.reqR.Request.RawURI
}

func (resp *RequestContext) Path() []byte {
	return resp.reqR.Request.URI.Path()
}

func (resp *RequestContext) Params() []h1.Query {
	return resp.reqR.Request.URI.Query()
}

func (rctx *RequestContext) WriteFullBody(status int, body []byte) error {
	rctx.response.StatusCode = status
	rctx.SetContentLength(len(body))
	_, err := rctx.Write(body)
	return err
}

func (rctx *RequestContext) WriteStream(status int, stream io.Reader) error {
	rctx.response.StatusCode = status
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
		_, err = rctx.Write((*buf)[:n])
		if err != nil {
			return err
		}
	}
}

func (rctx *RequestContext) SetContentLength(length int) {
	rctx.respW.ContentLength = length
}

func (rctx *RequestContext) SetHeader(name, value string) {
	for i := range rctx.response.Headers {
		if rctx.response.Headers[i].Name == name {
			rctx.response.Headers[i].Disabled = false
			rctx.response.Headers[i].Value = value
			return
		}
	}
	rctx.response.Headers = append(rctx.response.Headers, Header{
		Disabled: false,
		Name:     name,
		Value:    value,
	})
}

func (rctx *RequestContext) DeleteHeader(name string) {
	for i := range rctx.response.Headers {
		if rctx.response.Headers[i].Name == name {
			rctx.response.Headers[i].Disabled = true
			return
		}
	}
}

var requestPool = sync.Pool{
	New: func() interface{} {
		v := new(RequestContext)
		return v
	},
}

func GetRequestContext(upstream io.Writer) *RequestContext {
	ctx := requestPool.Get().(*RequestContext)
	ctx.respW = h1.GetResponse(upstream)
	return ctx
}

func PutRequestContext(ctx *RequestContext) {
	ctx.resetHard()
	requestPool.Put(ctx)
}

func (rctx *RequestContext) resetSoft() {
	rctx.hwt = false
	rctx.CloseBodyReader()
	rctx.response.reset()
}

func (rctx *RequestContext) resetHard() {
	rctx.resetSoft()
	rctx.conn = nil
	rctx.server = nil
	rctx.reqR.Reset()
	h1.PutResponse(rctx.respW)
}
