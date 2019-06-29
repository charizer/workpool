package workpool

import (
	"fmt"
)

//Worker 工作者单元, 用于执行Job的单元, 数量有限, 由调度中心分配
type Worker struct {
	workerPool chan *Worker
	jobChannel chan Job
	stop       chan struct{}
}

func NewWorker(pool chan *Worker) *Worker {
	return &Worker{
		workerPool: pool,
		jobChannel: make(chan Job),
		stop:       make(chan struct{}),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.workerPool <- w
			select {
			case job := <-w.jobChannel:
				if err := job.Exec(); err != nil {
					fmt.Sprintf("excute job failed with err: %v\n", err)
				}
			case <-w.stop:
				w.stop <- struct{}{}
				return
			}
		}
	}()
}
