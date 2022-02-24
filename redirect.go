package silverlining

func (r *Context) Redirect(status int, url string) {
	r.ResponseHeaders().Set("Location", url)
	r.SetContentLength(0)
	r.WriteHeader(status)
}
