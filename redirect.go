package silverlining

func (rctx *RequestContext) Redirect(status int, url string) {
	rctx.ResponseHeaders().Set("Location", url)
	rctx.SetContentLength(0)
	rctx.WriteHeader(status)
}
