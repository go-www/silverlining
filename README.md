[![GitHub last commit](https://img.shields.io/github/last-commit/go-www/silverlining?style=for-the-badge)](https://github.com/go-www/silverlining/commits/main)
[![Go Reference](https://img.shields.io/badge/Go-Reference-007d9c?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/go-www/silverlining)

# silverlining

Silverlining is a low-level HTTP Framework for Go Programming Language.

## Installation

```sh
go get -u github.com/go-www/silverlining
```

## Usage

```go
package main

import "github.com/go-www/silverlining"

func main() {
	silverlining.ListenAndServe(":8080", func(r *silverlining.Context) {
		r.WriteFullBodyString(200, "Hello, World!")
	})
}

```
