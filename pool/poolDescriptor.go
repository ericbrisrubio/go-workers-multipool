package pool

type Descriptor interface {
	SetWorkerFunc(fn func(interface{}) bool)
	AddTask(data interface{}) error
	AddWorkers(amount int) error
	KillWorkers(amount int) error
	EditWorkersAmount(amount int) error
	PauseAllWorkers()
	ResumeAllWorkers()
	Wait() error
}

