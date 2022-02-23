package main

import (
	"log"
	"net"
	"net/http"

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
	healthz := []byte("OK")
	jsonData := map[string]string{"hello": "world"}

	srv.Handler = func(r *silverlining.RequestContext) {
		switch string(r.Path()) {
		case "/":
			r.SetHeader("Content-Type", "text/plain")
			r.WriteFullBody(200, data)
		case "/json":
			r.WriteJSON(200, jsonData)
		case "/healthz":
			r.SetHeader("Content-Type", "text/plain")
			r.WriteFullBody(200, healthz)
		case "/redirect":
			r.Redirect(http.StatusSeeOther, "/")
		default:
			r.WriteFullBody(404, nil)
		}
	}

	err = srv.Serve(ln)
	if err != nil {
		log.Fatal(err)
	}
}
