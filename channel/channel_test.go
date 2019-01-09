package channel

import (
	"sync"
	"testing"
	"time"
)

func Test_Waiting(t *testing.T) {
	wg := New(50)
	p := 0
	var l sync.Mutex
	for i := 0; i < 1000; i++ {
		wg.Add(func() {
			time.Sleep(time.Second / 10)
			l.Lock()
			p ++
			l.Unlock()
		})
	}

	wg.Wait()
	if p != 1000 {
		t.Fatalf("p != 1000, instead == %d", p)
	}
}