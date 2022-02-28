# silverlining

Silverlining is a low-level HTTP Framework for Go Programming Language.

## Installation

```sh
go get -u github.com/go-www/silverlining
```

## Usage

```go
package main

import (
	"log"
	"net"

	"github.com/go-www/silverlining"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening on http://localhost:8080")

	defer ln.Close()

	srv := silverlining.Server{}

	srv.Handler = func(r *silverlining.Context) {
        r.ResponseHeaders().Set("Content-Type", "text/plain")
		r.WriteFullBodyString(200, "Hello, World!")
	}

	err = srv.Serve(ln)
	if err != nil {
		log.Fatal(err)
	}
}
```
