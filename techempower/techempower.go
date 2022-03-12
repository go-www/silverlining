package main

import (
	"flag"
	"log"

	"github.com/go-www/silverlining"
	"github.com/lemon-mint/envaddr"
)

var (
	prefork = flag.Bool("prefork", false, "use prefork")
)

func main() {
	flag.Parse()

	Handler := func(r *silverlining.Context) {
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

	var err error
	if *prefork {
		var id int
		id, err = silverlining.PreforkChildID()
		if err != nil {
			log.Fatalln(err)
		}

		if id == 0 {
			log.Println("Starting prefork leader process")
		} else {
			log.Printf("Starting prefork replica process %d", id)
		}
		err = silverlining.ListenAndServePrefork(envaddr.Get(":8080"), Handler)
	} else {
		err = silverlining.ListenAndServe(envaddr.Get(":8080"), Handler)
	}
	if err != nil {
		log.Fatalln(err)
	}
}
