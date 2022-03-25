package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	data = bytes.ReplaceAll(data, []byte("%"), []byte("%25"))
	data = bytes.ReplaceAll(data, []byte("\n"), []byte("%0A"))
	data = bytes.ReplaceAll(data, []byte("\r"), []byte("%0D"))

	fmt.Println("::set-output name=body::" + string(data))
}
