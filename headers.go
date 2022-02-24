package silverlining

type ResponseHeader struct {
	v *RequestContext
}

func (r ResponseHeader) Set(name, value string) {
	r.v.setHeader(name, value)
}

func (r ResponseHeader) Del(name string) {
	r.v.deleteHeader(name)
}

func (rctx *RequestContext) ResponseHeaders() ResponseHeader {
	return ResponseHeader{rctx}
}

func (rctx *RequestContext) setHeader(name, value string) {
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

func (rctx *RequestContext) deleteHeader(name string) {
	for i := range rctx.response.Headers {
		if rctx.response.Headers[i].Name == name {
			rctx.response.Headers[i].Disabled = true
			return
		}
	}
}
