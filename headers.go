package silverlining

import "github.com/go-www/silverlining/h1"

type ResponseHeaders struct {
	v *Context
}

func (r ResponseHeaders) Set(name, value string) {
	r.v.setHeader(name, value)
}

func (r ResponseHeaders) Del(name string) {
	r.v.deleteHeader(name)
}

func (r *Context) ResponseHeaders() ResponseHeaders {
	return ResponseHeaders{r}
}

func (r *Context) setHeader(name, value string) {
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

func (r *Context) deleteHeader(name string) {
	for i := range r.response.Headers {
		if r.response.Headers[i].Name == name {
			r.response.Headers[i].Disabled = true
			return
		}
	}
}

type RequestHeaders struct {
	v *Context
}

// (RequestHeaders).Get returns the value of the header with the given name.
//
// Returned value is valid until the next request.
func (r RequestHeaders) Get(name string) (string, bool) {
	return r.v.getHeader(name)
}

// (RequestHeaders).GetBytes returns the value of the header with the given name.
//
// Returned value is valid until the next request.
func (r RequestHeaders) GetBytes(name []byte) ([]byte, bool) {
	return r.v.getHeaderBytes(name)
}

// (RequestHeaders).List returns a slice of all the headers.
//
// Returned value is valid until the next request.
func (r RequestHeaders) List() []h1.Header {
	return r.v.reqR.Request.Headers
}

func (r *Context) RequestHeaders() RequestHeaders {
	return RequestHeaders{r}
}

func (r *Context) getHeader(name string) (string, bool) {
	h, ok := r.reqR.Request.GetHeader(stringToBytes(name))
	if !ok {
		return "", false
	}
	return bytesToString(h.RawValue), true
}

func (r *Context) getHeaderBytes(name []byte) ([]byte, bool) {
	h, ok := r.reqR.Request.GetHeader(name)
	if !ok {
		return nil, false
	}
	return h.RawValue, true
}

var host_header = stringToBytes("Host")

// (*Context).Host returns the value of the Host header.
// Returned value is valid until the next request.
func (r *Context) Host() string {
	h, ok := r.getHeaderBytes(host_header)
	if !ok {
		return ""
	}
	return bytesToString(h)
}

// (*Context).HostBytes returns the value of the Host header.
// Returned value is valid until the next request.
func (r *Context) HostBytes() []byte {
	h, ok := r.getHeaderBytes(host_header)
	if !ok {
		return nil
	}
	return h
}
