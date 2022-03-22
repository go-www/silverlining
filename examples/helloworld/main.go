package main

import "github.com/go-www/silverlining"

func main() {
	silverlining.ListenAndServe(":8080", func(r *silverlining.Context) {
		r.WriteFullBodyString(200, "Hello, World!")
	})
}
