package h1

import (
	"bytes"
	"testing"
)

func Test_stricmp(t *testing.T) {
	type args struct {
		a []byte
		b []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"empty",
			args{[]byte(""), []byte("")},
			true,
		},
		{
			"same string",
			args{[]byte("hello"), []byte("hello")},
			true,
		},
		{
			"different string",
			args{[]byte("hello"), []byte("world")},
			false,
		},
		{
			"Case-insensitive",
			args{[]byte("hello"), []byte("hElLO")},
			true,
		},
		{
			"different length",
			args{[]byte("hello"), []byte("hello, world")},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stricmp(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("stricmp() = %v, want %v", got, tt.want)
			}
		})
	}
}

var testdata = [][2][]byte{
	{[]byte("hello"), []byte("hello")},
	{[]byte("hello"), []byte("world")},
	{[]byte("hello"), []byte("hElLO")},
	{[]byte("hello"), []byte("hello, world")},
	{[]byte("hello"), []byte("")},
	{[]byte(""), []byte("hello")},
	{[]byte(""), []byte("")},
	{[]byte("a"), []byte("A")},
	{[]byte("A"), []byte("a")},
	{[]byte("A"), []byte("A")},
	{[]byte("Content-Length"), []byte("content-length")},
	{[]byte("content-length"), []byte("Content-Length")},
	{[]byte("content-length"), []byte("content-length")},
	{[]byte("content-length"), []byte("content-Length")},
	{[]byte("Content-length"), []byte("content-length")},
}

func Benchmark_stricmp(b *testing.B) {
	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				for _, d := range testdata {
					stricmp(d[0], d[1])
				}
			}
		},
	)
}

func Benchmark_Bytes_EqualFold(b *testing.B) {
	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				for _, d := range testdata {
					bytes.EqualFold(d[0], d[1])
				}
			}
		},
	)
}

var testdata2 = [][2][]byte{
	{[]byte("Content-Length"), []byte("Content-Length")},
	{[]byte("Content-Length"), []byte("content-length")},
	{[]byte("Content-Length"), []byte("Accept-Encoding")},
	{[]byte("Content-Length"), []byte("Origin")},
	{[]byte("Content-Length"), []byte("Host")},
	{[]byte("Content-Length"), []byte("User-Agent")},
	{[]byte("Content-Length"), []byte("Referer")},
	{[]byte("Content-Length"), []byte("Accept")},
	{[]byte("Content-Length"), []byte("Accept-Language")},
	{[]byte("Content-Length"), []byte("Accept-Charset")},
	{[]byte("Content-Length"), []byte("Accept-Datetime")},
	{[]byte("Content-Length"), []byte("Authorization")},
	{[]byte("Content-Length"), []byte("Cache-Control")},
	{[]byte("Content-Length"), []byte("Connection")},
	{[]byte("Content-Length"), []byte("Cookie")},
	{[]byte("Content-Length"), []byte("Date")},
	{[]byte("Content-Length"), []byte("Expect")},
}

func Benchmark_ContentLength_stricmp(b *testing.B) {
	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				for _, d := range testdata2 {
					stricmp(d[0], d[1])
				}
			}
		},
	)
}

func Benchmark_ContentLength_Bytes_EqualFold(b *testing.B) {
	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				for _, d := range testdata2 {
					bytes.EqualFold(d[0], d[1])
				}
			}
		},
	)
}
