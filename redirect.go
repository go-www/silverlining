package silverlining

func (rctx *RequestContext) Redirect(status int, url string) {
	rctx.SetHeader("Location", url)
	rctx.SetContentLength(0)
	rctx.WriteHeader(status)
}
