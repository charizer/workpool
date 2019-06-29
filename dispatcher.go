package workpool

type Dispatcher struct {
	workerPool chan *Worker
	jobQueue   chan Job
	stop       chan struct{}
}

func NewDispatcher(workerPool chan *Worker, jobQueue chan Job) *Dispatcher {
	d := &Dispatcher{
		workerPool: workerPool,
		jobQueue:   jobQueue,
		stop:       make(chan struct{}),
	}

	for i := 0; i < cap(d.workerPool); i++ {
		worker := NewWorker(d.workerPool)
		worker.Start()
	}
	go d.Dispatch()
	return d
}

func (d *Dispatcher) Dispatch() {
	for {
		select {
		case job := <-d.jobQueue:
			worker := <-d.workerPool
			worker.jobChannel <- job
		case <-d.stop:
			for i := 0; i < cap(d.workerPool); i++ {
				worker := <-d.workerPool
				worker.stop <- struct{}{}
				<-worker.stop
			}
			d.stop <- struct{}{}
			return
		}
	}
}
