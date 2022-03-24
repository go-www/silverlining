package h1

import (
	"bytes"
	"testing"
)

func Benchmark_Request_Reader(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		buffer := bytes.Buffer{}
		buffer.Write(TestFullReqData)
		buffer.Reset()

		readbuffer := make([]byte, 8192)

		r := &RequestReader{
			R:          &buffer,
			ReadBuffer: readbuffer,
		}

		for p.Next() {
			r.Reset()
			r.ReadBuffer = readbuffer
			r.NextBuffer = nil

			buffer.Write(TestFullReqData)
			_, err := r.Next()
			if err != nil {
				b.Error(err)
			}
			buffer.Reset()
		}
	})
}
