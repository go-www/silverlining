package gopool

import "sync"

var pool = sync.Pool{}

func init() {
	pool.New = func() interface{} {
		ch := make(chan func(), 1)
		go func() {
			for f := range ch {
				f()
				pool.Put(ch)
			}
		}()
		return ch
	}
}

func Go(f func()) {
	ch := pool.Get().(chan func())
	ch <- f
}
