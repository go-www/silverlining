package silverlining

type ResponseHeaders struct {
	v *RequestContext
}

func (r ResponseHeaders) Set(name, value string) {
	r.v.setHeader(name, value)
}

func (r ResponseHeaders) Del(name string) {
	r.v.deleteHeader(name)
}

func (rctx *RequestContext) ResponseHeaders() ResponseHeaders {
	return ResponseHeaders{rctx}
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
