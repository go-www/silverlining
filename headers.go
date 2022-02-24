package silverlining

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
