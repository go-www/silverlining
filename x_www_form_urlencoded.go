package silverlining

import (
	"sync"

	"github.com/go-www/h1"
)

type KV = h1.Query

type XWWWFormURLEncoded struct {
	value []KV
}

var xwwwFormURLEncodedPool sync.Pool = sync.Pool{
	New: func() interface{} {
		return &XWWWFormURLEncoded{}
	},
}

func (x *XWWWFormURLEncoded) reset() {
	x.value = x.value[:0]
}

func getXWWWFormURLEncoded() *XWWWFormURLEncoded {
	return xwwwFormURLEncodedPool.Get().(*XWWWFormURLEncoded)
}

func PutXWWWFormURLEncoded(x *XWWWFormURLEncoded) {
	x.reset()
	xwwwFormURLEncodedPool.Put(x)
}

func (r *Context) XWWWFormURLEncoded(maxSize int64) (*XWWWFormURLEncoded, error) {
	body, err := r.FastBodyUnsafe(maxSize)
	if err != nil {
		return nil, err
	}

	x := getXWWWFormURLEncoded()
	x.value = h1.ParseRawQuery(body, x.value)
	return x, nil
}

func (x *XWWWFormURLEncoded) Len() int {
	return len(x.value)
}

func (x *XWWWFormURLEncoded) Get(key string) ([]byte, bool) {
	for _, kv := range x.value {
		if string(kv.Key) == key {
			return kv.Value, true
		}
	}
	return nil, false
}

func (x *XWWWFormURLEncoded) GetString(key string) (string, bool) {
	v, ok := x.Get(key)
	return string(v), ok
}

func (x *XWWWFormURLEncoded) GetStringUnsafe(key string) (string, bool) {
	v, ok := x.Get(key)
	if !ok {
		return "", false
	}
	return bytesToString(v), true
}

func (x *XWWWFormURLEncoded) Close() {
	PutXWWWFormURLEncoded(x)
}
