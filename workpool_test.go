package workpool

import (
	"io/ioutil"
	"log"
	"runtime"
	"testing"
)

func init() {
	log.Println("using MAXPROC")
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)
}

type TestBenchJob struct {
	Num int
}

func (t *TestBenchJob) Exec() error {
	log.Printf("I am worker! Number %d\n", t.Num)
	return nil
}

func BenchmarkPool(b *testing.B) {
	pool := NewPool(1, 10)
	defer pool.Release()

	log.SetOutput(ioutil.Discard)

	for n := 0; n < b.N; n++ {
		pool.Put(&TestBenchJob{Num: n})
	}
}
