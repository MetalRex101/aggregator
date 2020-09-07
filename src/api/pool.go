package api

import (
	"sync"
)

func NewPool(size int) *Pool {
	return &Pool{
		workers: make(chan struct{}, size),
	}
}

type Pool struct {
	workers chan struct{}
}

func (p *Pool) Run(callback func()) {
	p.workers <- struct{}{}
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer func() {
			<-p.workers
			wg.Done()
		}()

		callback()
	}(wg)

	wg.Wait()
}