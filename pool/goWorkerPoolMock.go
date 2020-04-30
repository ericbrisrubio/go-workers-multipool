package pool

import "github.com/enriquebris/goworkerpool"

type GoWorkerPoolDefinerMock struct {
	totalWorkers int
	SetWorkerFuncHasBeenCalled bool
	AddTaskFuncHasBeenCalled   bool
	AddWorkersHasBeenCalled    bool
	KillWorkersHasBeenCalled    bool
	*goworkerpool.Pool
}

func (definer *GoWorkerPoolDefinerMock) SetWorkerFunc(fn func(interface{}) bool) {
	definer.SetWorkerFuncHasBeenCalled = true
}

func (definer *GoWorkerPoolDefinerMock) AddTask(data interface{}) error {
	definer.AddTaskFuncHasBeenCalled = true
	return nil
}

func (definer *GoWorkerPoolDefinerMock) AddWorkers(amount int) error {
	definer.AddWorkersHasBeenCalled = true
	return nil
}

func (definer *GoWorkerPoolDefinerMock) KillWorkers(amount int) error {
	definer.KillWorkersHasBeenCalled = true
	return nil
}

func (definer *GoWorkerPoolDefinerMock) EditWorkersAmount(workersAmount int) error {
	return nil
}

func (definer *GoWorkerPoolDefinerMock) PauseAllWorkers() {
}

func (definer *GoWorkerPoolDefinerMock) ResumeAllWorkers() {
}

func (definer *GoWorkerPoolDefinerMock) WaitWhileAlive() error {
	return nil
}
