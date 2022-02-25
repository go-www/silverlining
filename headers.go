package silverlining

import "github.com/go-www/h1"

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

func (r RequestHeaders) Get(name string) (string, bool) {
	return r.v.getHeader(name)
}

func (r RequestHeaders) GetBytes(name []byte) ([]byte, bool) {
	return r.v.getHeaderBytes(name)
}

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
