package silverlining

type ServerSentEventWriter struct {
	v *Context
}

func (r *Context) ServerSentEventWriter() *ServerSentEventWriter {
	r.ResponseHeaders().Set("Content-Type", "text/event-stream")
	r.SetContentLength(-1)
	return &ServerSentEventWriter{r}
}

func (s ServerSentEventWriter) Flush() error {
	return s.v.Flush()
}

var sse_newline = []byte("\n")
var sse_newline_newline = []byte("\n\n")
var sse_event = []byte("event: ")
var sse_data = []byte("data: ")

func (s ServerSentEventWriter) Send(event string, data string) error {
	_, err := s.v.Write(sse_event)
	if err != nil {
		return err
	}
	_, err = s.v.WriteString(event)
	if err != nil {
		return err
	}
	_, err = s.v.Write(sse_newline)
	if err != nil {
		return err
	}

	_, err = s.v.Write(sse_data)
	if err != nil {
		return err
	}
	_, err = s.v.WriteString(data)
	if err != nil {
		return err
	}
	_, err = s.v.Write(sse_newline_newline)
	if err != nil {
		return err
	}

	return s.v.Flush()
}

var sse_heartbeat = []byte(":heartbeat\n\n")

func (s ServerSentEventWriter) WriteHeartbeat() error {
	_, err := s.v.Write(sse_data)
	if err != nil {
		return err
	}
	_, err = s.v.Write(sse_heartbeat)
	if err != nil {
		return err
	}
	return s.v.Flush()
}
