package silverlining

import (
	"net"
	"os"
	"os/exec"
	"runtime"

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

func ListenAndServePrefork(addr string, handler Handler) error {
	runtime.GOMAXPROCS(1)

	IsChildEnv := os.Getenv("SILVERLINING_PREFORK_CHILD")
	IsChild := IsChildEnv == "1"

	if !IsChild {
		numCPU := runtime.NumCPU()
		var env []string
		env = append(env, os.Environ()...)
		env = append(env, "GOMAXPROCS=1", "SILVERLINING_PREFORK_CHILD=1")

		for i := 0; i < numCPU-1; i++ {
			cmd := exec.Command(os.Args[0], os.Args[1:]...)
			cmd.Env = env
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			if err := cmd.Start(); err != nil {
				return err
			}
		}
	}

	return ListenAndServeReusePort(addr, handler)
}
