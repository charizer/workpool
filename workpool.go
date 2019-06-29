package workpool

import "sync"

type Pool struct {
	JobQueue   chan Job
	dispatcher *Dispatcher
	wg         sync.WaitGroup
}

func NewPool(numWorkers int, jobQueueLen int) *Pool {
	jobQueue := make(chan Job, jobQueueLen)
	workerPool := make(chan *Worker, numWorkers)

	pool := &Pool{
		JobQueue:   jobQueue,
		dispatcher: NewDispatcher(workerPool, jobQueue),
	}
	return pool
}

func (p *Pool) Put(job Job){
	p.JobQueue <- job
}

func (p *Pool) JobDone() {
	p.wg.Done()
}

func (p *Pool) WaitCount(count int) {
	p.wg.Add(count)
}

func (p *Pool) WaitAll() {
	p.wg.Wait()
}

func (p *Pool) Release() {
	p.dispatcher.stop <- struct{}{}
	<-p.dispatcher.stop
}
