package main

type Pool struct {
	Chan chan *Worker
}

type Worker struct {
	Id      int
	fun     func()
	Working bool
}

func NewPool(size int) *Pool {
	p := &Pool{
		Chan: make(chan *Worker, size),
	}
	go p.Consumer()
	return p
}
func (p *Pool) Go(f func()) {
	work := &Worker{
		Id:      len(p.Chan),
		fun:     f,
		Working: true,
	}
	p.Chan <- work
}
func (p *Pool) Consumer() {
	for {
		select {
		case w := <-p.Chan:
			w.fun()
		}
	}
}
