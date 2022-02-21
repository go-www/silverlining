package main

import (
	"log"
	"net"

	"github.com/go-www/silverlining"
	"github.com/lemon-mint/envaddr"
)

func main() {
	ln, err := net.Listen("tcp", envaddr.Get(":8080"))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening on http://%s", ln.Addr())

	defer ln.Close()

	srv := silverlining.Server{}

	data := []byte("Hello, World!")

	srv.Handler = func(r *silverlining.RequestContext) {
		r.SetContentLength(len(data))
		r.WriteHeader(200)
		r.Write(data)
	}

	err = srv.Serve(ln)
	if err != nil {
		log.Fatal(err)
	}
}
