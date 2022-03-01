package silverlining

type ChunckedBodyWriter struct {
	v *Context
}

func (c ChunckedBodyWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	_, err := c.v.respW.WriteUint64Hex(uint64(len(p)))
	if err != nil {
		return 0, err
	}
	_, err = c.v.respW.Write(crlf)
	if err != nil {
		return 0, err
	}
	_, err = c.v.respW.Write(p)
	if err != nil {
		return 0, err
	}
	_, err = c.v.respW.Write(crlf)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (c ChunckedBodyWriter) WriteString(s string) (int, error) {
	return c.Write(stringToBytes(s))
}

var chunkClose = []byte("0\r\n\r\n")

func (c ChunckedBodyWriter) Close() error {
	_, err := c.v.respW.Write(chunkClose)
	if err != nil {
		return err
	}
	return nil
}

func (r *Context) ChunckedBodyWriter() ChunckedBodyWriter {
	r.SetContentLength(-2)
	r.ResponseHeaders().Set("Transfer-Encoding", "chunked")
	return ChunckedBodyWriter{r}
}
