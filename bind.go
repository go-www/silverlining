package silverlining

func (rctx *RequestContext) BindJSON(v any) error {
	return rctx.ReadJSON(v)
}
