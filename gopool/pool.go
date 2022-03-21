package gopool

func Go(f func()) {
	go f()
}
