package silverlining

import (
	"net"

	reuse "github.com/libp2p/go-reuseport"
)

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

func ListenAndServeReusePort(addr string, handler Handler) error {
	ln, err := reuse.Listen("tcp", addr)
	if err != nil {
		return err
	}

	srv := &Server{
		Listener: ln,
		Handler:  handler,
	}

	return srv.Serve(ln)
}
