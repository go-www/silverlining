package h1

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"sync"
)

type Request struct {
	// Request line
	Method  Method
	RawURI  []byte
	Version []byte

	// Headers
	Headers []Header

	// Parsed URI
	URI URI

	ContentLength int64
}

var requestPool = sync.Pool{
	New: func() any {
		return &Request{}
	},
}

func (r *Request) Reset() {
	r.Method = MethodInvalid
	r.RawURI = nil
	r.Version = nil
	r.Headers = r.Headers[:0]
	r.ContentLength = 0
}

func (r *Request) GetHeader(name []byte) (*Header, bool) {
	for i := range r.Headers {
		if stricmp(r.Headers[i].Name, name) {
			return &r.Headers[i], true
		}
	}
	return nil, false
}

func GetRequest() *Request {
	return requestPool.Get().(*Request)
}

func PutRequest(r *Request) {
	r.Reset()
	requestPool.Put(r)
}

var ErrInvalidMethod = errors.New("invalid method")
var ErrInvalidURI = errors.New("invalid uri")
var ErrInvalidVersion = errors.New("invalid version")

var ErrBufferTooSmall = errors.New("buffer too small")
var ErrRequestHeaderTooLarge = errors.New("request header too large")

func splitLine(src []byte) (line, rest []byte, err error) {
	idx := bytes.IndexByte(src, '\n')
	if idx < 1 { // 0: cr 1: lf
		return nil, src, ErrBufferTooSmall
	}

	if src[idx-1] == '\r' {
		line = src[:idx-1]
		rest = src[idx+1:]
		return
	}
	return src[:idx], src[idx+1:], nil
}

func parseRequestLineforTest(src []byte) (method Method, uri []byte, version []byte, next []byte, err error) {
	req := Request{}
	next, err = ParseRequestLine(&req, src)
	if err != nil {
		return MethodInvalid, nil, nil, nil, err
	}
	return req.Method, req.RawURI, req.Version, next, nil
}

var methodTable = [256]Method{}

var _ = func() int {
	//GET
	const GETIndex = 'G' ^ 'E' + 'T'
	methodTable[GETIndex] = MethodGET
	//PUT
	const PUTIndex = 'P' ^ 'U' + 'T'
	methodTable[PUTIndex] = MethodPUT
	//HEAD
	const HEADIndex = 'H' ^ 'E' + 'A'
	methodTable[HEADIndex] = MethodHEAD
	//POST
	const POSTIndex = 'P' ^ 'O' + 'S'
	methodTable[POSTIndex] = MethodPOST
	//BREW
	const BREWIndex = 'B' ^ 'R' + 'E'
	methodTable[BREWIndex] = MethodBREW
	//TRACE
	const TRACEIndex = 'T' ^ 'R' + 'A'
	methodTable[TRACEIndex] = MethodTRACE
	//PATCH
	const PATCHIndex = 'P' ^ 'A' + 'T'
	methodTable[PATCHIndex] = MethodPATCH
	//DELETE
	const DELETEIndex = 'D' ^ 'E' + 'L'
	methodTable[DELETEIndex] = MethodDELETE
	//CONNECT
	const CONNECTIndex = 'C' ^ 'O' + 'N'
	methodTable[CONNECTIndex] = MethodCONNECT
	//OPTIONS
	const OPTIONSIndex = 'O' ^ 'P' + 'T'
	methodTable[OPTIONSIndex] = MethodOPTIONS

	// all methods should have distinct index number
	var _ = map[int]Method{
		GETIndex:     MethodGET,
		PUTIndex:     MethodPUT,
		HEADIndex:    MethodHEAD,
		POSTIndex:    MethodPOST,
		BREWIndex:    MethodBREW,
		TRACEIndex:   MethodTRACE,
		PATCHIndex:   MethodPATCH,
		DELETEIndex:  MethodDELETE,
		CONNECTIndex: MethodCONNECT,
		OPTIONSIndex: MethodOPTIONS,
	}

	return 0
}()

func ParseRequestLine(dst *Request, src []byte) (next []byte, err error) {
	next = src
	var line []byte
	line, next, err = splitLine(next)
	if err != nil {
		return next, err
	}
	MethodIndex := bytes.IndexByte(line, ' ')
	if MethodIndex < 0 || MethodIndex < 3 {
		return next, ErrInvalidMethod
	}
	URIIndex := bytes.IndexByte(line[MethodIndex+1:], ' ')
	if URIIndex < 0 {
		return next, ErrInvalidURI
	}
	dst.RawURI = line[MethodIndex+1 : MethodIndex+1+URIIndex]
	dst.Version = line[MethodIndex+1+URIIndex+1:]

	m := line[:MethodIndex]

	dst.Method = methodTable[m[0]^m[1]+m[2]]
	return next, nil
}

var ContentLengthHeader = []byte("Content-Length")

func ParseHeaders(dst *Request, src []byte) (next []byte, err error) {
	next = src
	var line []byte
	for {
		line, next, err = splitLine(next)
		if err != nil {
			return next, err
		}
		if len(line) == 0 {
			break
		}
		h := Header{}
		h.Name, h.RawValue = ParseHeaderLine(line)
		dst.Headers = append(dst.Headers, h)

		if stricmp(h.Name, ContentLengthHeader) {
			dst.ContentLength, err = ParseContentLength(h.RawValue)
			if err != nil {
				return next, err
			}
		}
	}
	return next, nil
}

func ParseContentLength(src []byte) (int64, error) {
	srcS := bytesToString(src)
	return strconv.ParseInt(srcS, 10, 64)
}

func ParseHeaderLine(src []byte) (name []byte, value []byte) {
	idx := bytes.IndexByte(src, ':')
	if idx < 0 {
		return src[:0], nil
	}
	// RFC2616 Section 4.2
	// Remove all leading and trailing LWS on field contents

	// skip leading LWS
	var i int = idx + 1
	for ; i < len(src); i++ {
		if src[i] != ' ' && src[i] != '\t' {
			break
		}
	}
	// skip trailing LWS
	var j int = len(src) - 1
	for ; j > i; j-- {
		if src[j] != ' ' && src[j] != '\t' {
			break
		}
	}
	return src[:idx], src[i : j+1]
}

type Header struct {
	raw []byte

	Name     []byte
	RawValue []byte
}

func (h *Header) Reset() {
	h.raw = nil
	h.Name = nil
	h.RawValue = nil
}

const BufferPoolSize = 4096

var bufferPool = sync.Pool{
	New: func() any {
		buffer := make([]byte, BufferPoolSize)
		return &buffer
	},
}

func GetBuffer() *[]byte {
	return bufferPool.Get().(*[]byte)
}

func PutBuffer(b *[]byte) {
	if cap(*b) >= BufferPoolSize {
		bufferPool.Put(b)
	}
}

// Do not use this function in production code.
// This function is only for testing purpose.
func ParseRequest(dst *Request, r io.Reader) (next []byte, err error) {
	dst.Reset()
	var buffer *[]byte = GetBuffer()
	//defer PutBuffer(buffer) // Allow GC to collect the buffer
	n, err := r.Read(*buffer)
	if err != nil {
		return nil, err
	}
	next = (*buffer)[:n]
retryRead:

	// This function can't parse request line correctly if the request line is too long (>=4096)
	next, err = ParseRequestLine(dst, next)
	if err == ErrBufferTooSmall {
		buffer = GetBuffer()
		//defer PutBuffer(buffer)
		remainBytes := copy(*buffer, next)
		n, err = r.Read((*buffer)[remainBytes:])
		if err != nil {
			return next, err
		}
		next = (*buffer)[:remainBytes+n]
		goto retryRead
	} else if err != nil {
		return next, err
	}

	for {
		next, err = ParseHeaders(dst, next)
		if err == ErrBufferTooSmall {
			buffer = GetBuffer()
			//defer PutBuffer(buffer)
			remainBytes := copy(*buffer, next)
			n, err = r.Read((*buffer)[remainBytes:])
			if err != nil {
				return next, err
			}
			next = (*buffer)[:remainBytes+n]
			continue
		}
		if err != nil {
			return next, err
		}
		break
	}
	return next, nil
}

func parseRequestForTest(data []byte) (*Request, error) {
	r := &Request{}
	_, err := ParseRequest(r, bytes.NewReader(data))
	return r, err
}

func parseRequestForTestIsValid(data []byte) bool {
	_, err := parseRequestForTest(data)
	return err == nil
}
