package closable

import (
	"fmt"
	"io"
	"os"
	"sync"
)

var (
	closer []func() error
	mu sync.Mutex
)

func Push(c io.Closer) {
	mu.Lock()
	defer mu.Unlock()
	closer = append(closer, c.Close)
}

func PushCloseFun(closeFunc func() error)  {
	mu.Lock()
	defer mu.Unlock()
	closer = append(closer, closeFunc)
}

func Done() {
	for l := len(closer); l > 0; l = len(closer) {
		if err := closer[l-1](); err != nil {
			_, _ = fmt.Fprint(os.Stderr, err)
		}
		closer = closer[:l-1]
	}
}