# workpool
功能：一个轻量的goroutine pool实现，它由Job、dispatcher、work、pool组成。通过pool.put投放job到pool中，dispatcher负责分发job到可用的worker中，然后worker执行job，自定义job实现Job的Exec方法即可。执行完job后返回pool。

### 基本数据结构：
```
//工作池，将Job投放到JobQueue中
type Pool struct {
    //投放job的chan
	JobQueue   chan Job
	//调度中心分派器
	dispatcher *Dispatcher
	//waitgroup
	wg         sync.WaitGroup
}

//调度中心分派器，从jobQueue中获取Job，从workerPool中获取一个可用worker，然后将Job投放到worker的jobChannel中
type Dispatcher struct {
	workerPool chan *Worker
	jobQueue   chan Job
	stop       chan struct{}
}

//Worker 工作者单元, 用于执行Job的单元, 从jobChannel中获取Job
type Worker struct {
	workerPool chan *Worker
	jobChannel chan Job
	stop       chan struct{}
}

```

### 使用示例：
```
type MyJob struct {
	Num  uint64
	Pool *Pool
}

func (m *MyJob) Exec() error {
	defer t.Pool.JobDone()
	log.Printf("I am worker! Number %d\n", m.Num)
	return nil
}
pool := NewPool(1000, 10000)
defer pool.Release()

iterations := 100
pool.WaitCount(iterations)
arg := uint64(1)
for i := 0; i < iterations; i++ {
	job := &MyJob{
		Num:  arg,
		Pool: pool,
	}
	atomic.AddUint64(&arg, 1)
	pool.JobQueue <- job
}
pool.WaitAll()
log.Println("work run over")

```



