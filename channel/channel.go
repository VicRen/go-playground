package channel

type WG struct {
	main    chan func()
	allDone chan bool
}

func New(n int) WG {
	res := WG {
		main:make(chan func()),
		allDone:make(chan bool),
	}

	procDone := make(chan bool)
	for i := 0; i < n; i++ {
		go func() {
			for {
				f := <- res.main
				if f == nil {
					procDone <- true
					return
				}
				f()
			}
		}()
	}
	go func() {
		for i := 0; i < n; i++ {
			<- procDone
		}
		res.allDone <- true
	}()
	return res
}

func (wg WG) Add(f func()) {
	wg.main <- f
}

func (wg WG) Wait() {
	close(wg.main)
	<- wg.allDone
}
