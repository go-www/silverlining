package h1

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type FastDateServer struct {
	serverName string

	dates []*[]byte
	index int

	current *[]byte

	once sync.Once
	stop chan struct{}
}

func NewFastDateServer(serverName string) *FastDateServer {
	fds := &FastDateServer{
		dates:      make([]*[]byte, 32),
		serverName: serverName,
		index:      0,
		stop:       make(chan struct{}),
	}

	for i := range fds.dates {
		v := make([]byte, 0, len("Server: ")+len(serverName)+len("\r\nDate: ")+len(time.RFC1123)*2+len("\r\n\r\n"))
		fds.dates[i] = &v
	}

	fds.updateDate(time.Now().UTC())

	return fds
}

func (fds *FastDateServer) updateDate(date time.Time) {
	new := fds.dates[fds.index%len(fds.dates)]
	fds.index++

	*new = (*new)[:0]
	*new = append(*new, "Date: "...)
	*new = append(*new, date.Format(time.RFC1123)...)
	*new = append(*new, "\r\nServer: "...)
	*new = append(*new, fds.serverName...)
	*new = append(*new, "\r\n"...)

	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&fds.current)), unsafe.Pointer(new))
}

func (fds *FastDateServer) GetDate() []byte {
	return *fds.current
}

func (fds *FastDateServer) Start() {
	fds.once.Do(func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fds.updateDate(time.Now().UTC())
			case <-fds.stop:
				return
			}
		}
	})
}

func (fds *FastDateServer) Stop() {
	fds.stop <- struct{}{}
	close(fds.stop)

	fds.once = sync.Once{}
}
