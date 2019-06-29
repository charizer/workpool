package workpool

type Job interface {
	Exec() error
}
