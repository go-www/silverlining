package h1

import (
	"fmt"
	"net/url"
	"testing"
)

func Benchmark_Net_URL_Parse(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, err := url.Parse("/oauth2/authorize?response_type=code&client_id=foo&redirect_uri=http%3A%2F%2Fexample.com%2Fcb&scope=email%20profile&state=xyz&nonce=abc")
			if err != nil {
				panic(err)
			}
		}
	})
}

func Benchmark_H1_URI_Parse(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		uri := URI{}
		data := []byte("/oauth2/authorize?response_type=code&client_id=foo&redirect_uri=http%3A%2F%2Fexample.com%2Fcb&scope=email%20profile&state=xyz&nonce=abc")
		for p.Next() {
			uri.Parse(data)
		}
	})
}

func Benchmark_H1_URI_Query(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		uri := URI{}
		data := []byte("/oauth2/authorize?response_type=code&client_id=foo&redirect_uri=http%3A%2F%2Fexample.com%2Fcb&scope=email%20profile&state=xyz&nonce=abc")
		for p.Next() {
			uri.Parse(data)
			uri.Query()
		}
	})
}

func Test_H1_URI_Parse(t *testing.T) {
	uri := URI{}
	urls := []string{
		"/some/path?foo=bar&baz=qux",
		"/oauth2/authorize?response_type=code&client_id=foo&redirect_uri=http%3A%2F%2Fexample.com%2Fcb&scope=email%20profile&state=xyz&nonce=abc",
		"/?foo=bar&baz=qux",
	}

	for _, u := range urls {
		t.Run(fmt.Sprintf("H1_URI_Parse (%s)", u), func(t *testing.T) {
			uri.Parse([]byte(u))
			q := uri.Query()

			netu, err := url.Parse(u)
			if err != nil {
				t.Fatal(err)
			}

			if len(q) != len(netu.Query()) {
				t.Fatalf("expected %d query args got %d", len(netu.Query()), len(q))
			}
		})
	}
}
