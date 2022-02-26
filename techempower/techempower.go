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

	srv := silverlining.Server{
		//ReadTimeout: time.Minute,
	}

	srv.Handler = func(r *silverlining.Context) {
		switch string(r.Path()) {
		case "/plaintext":
			r.ResponseHeaders().Set("Content-Type", "text/plain")
			r.WriteFullBodyString(200, "Hello, World!")
		case "/json":
			type Message struct {
				Message string `json:"message"`
			}
			msg := Message{Message: "Hello, World!"}
			r.WriteJSON(200, msg)
		default:
			r.WriteFullBody(404, nil)
		}
	}

	err = srv.Serve(ln)
	if err != nil {
		log.Fatal(err)
	}
}
