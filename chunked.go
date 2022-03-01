package silverlining

type ChunkedBodyWriter struct {
	v *Context
}

func (c ChunkedBodyWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	c.v.WriteHeader(c.v.response.StatusCode)
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

func (c ChunkedBodyWriter) WriteString(p string) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	c.v.WriteHeader(c.v.response.StatusCode)
	_, err := c.v.respW.WriteUint64Hex(uint64(len(p)))
	if err != nil {
		return 0, err
	}
	_, err = c.v.respW.Write(crlf)
	if err != nil {
		return 0, err
	}
	_, err = c.v.respW.WriteString(p)
	if err != nil {
		return 0, err
	}
	_, err = c.v.respW.Write(crlf)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

var chunkClose = []byte("0\r\n\r\n")

func (c ChunkedBodyWriter) Close() error {
	_, err := c.v.respW.Write(chunkClose)
	if err != nil {
		return err
	}
	return nil
}

func (r *Context) ChunkedBodyWriter() ChunkedBodyWriter {
	r.SetContentLength(-2)
	r.ResponseHeaders().Set("Transfer-Encoding", "chunked")
	return ChunkedBodyWriter{r}
}
