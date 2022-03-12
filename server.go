package silverlining

import (
	"errors"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"

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

func PreforkIsChild() bool {
	return os.Getenv("SILVERLINING_PREFORK_CHILD") == "1"
}

var ErrPreforkChildIDNotFound = errors.New("child id not found")

func PreforkChildID() (int, error) {
	if !PreforkIsChild() {
		return 0, nil
	}

	ChildIDEnv, ok := os.LookupEnv("SILVERLINING_PREFORK_CHILD_ID")
	if !ok {
		return 0, ErrPreforkChildIDNotFound
	}

	ChildID, err := strconv.Atoi(ChildIDEnv)
	if err != nil {
		return 0, err
	}

	return ChildID, nil
}

func ListenAndServePrefork(addr string, handler Handler) error {
	runtime.GOMAXPROCS(1)

	if !PreforkIsChild() {
		numCPU := runtime.NumCPU()
		for i := 0; i < numCPU-1; i++ {
			var env []string
			env = append(env, os.Environ()...)
			env = append(env, "GOMAXPROCS=1", "SILVERLINING_PREFORK_CHILD=1", "SILVERLINING_PREFORK_CHILD_ID="+strconv.Itoa(i+1))

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
