package silverlining

import "net"

func ListenAndServe(addr string, handler Handler) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	srv := &Server{
		Listener: ln,
		Handler:  handler,
	}

	return srv.Serve(ln)
}
