package workpool

import (
	"log"
	"runtime"
	"sync/atomic"
	"testing"
)

func init() {
	log.Println("using MAXPROC")
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)
}

type TestBenchJob struct {
	Num  uint64
	Pool *Pool
}

func (t *TestBenchJob) Exec() error {
	defer t.Pool.JobDone()
	log.Printf("I am worker! Number %d\n", t.Num)
	return nil
}

func TestNewPool(t *testing.T) {
	pool := NewPool(1000, 10000)
	defer pool.Release()

	iterations := 100
	pool.WaitCount(iterations)
	arg := uint64(1)
	for i := 0; i < iterations; i++ {
		job := &TestBenchJob{
			Num:  arg,
			Pool: pool,
		}
		atomic.AddUint64(&arg, 1)
		pool.JobQueue <- job
	}
	pool.WaitAll()
	t.Log("work run over")
}
