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
	jsonData := map[string]string{"hello": "world"}

	srv.Handler = func(r *silverlining.RequestContext) {
		switch string(r.URI()) {
		case "/":
			r.WriteFullBody(200, data)
		case "/json":
			r.WriteJSON(200, jsonData)
		default:
			r.WriteFullBody(404, nil)
		}
	}

	err = srv.Serve(ln)
	if err != nil {
		log.Fatal(err)
	}
}
