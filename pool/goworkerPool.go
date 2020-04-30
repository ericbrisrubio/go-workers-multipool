package pool

import (
	"errors"
	"github.com/enriquebris/goworkerpool"
)


type GoWorkerPoolDefiner struct {
	*goworkerpool.Pool
}

func (definer *GoWorkerPoolDefiner) SetWorkerFunc(fn func(interface{})bool) {
	definer.Pool.SetWorkerFunc(fn)
}

func (definer *GoWorkerPoolDefiner) AddTask(data interface{}) error {
	return definer.Pool.AddTask(data)
}

func (definer *GoWorkerPoolDefiner) AddWorkers(amount int) error {
	return definer.Pool.AddWorkers(amount)
}

func (definer *GoWorkerPoolDefiner) KillWorkers(amount int) error {
	if amount > definer.GetTotalWorkers() {
		return errors.New("cannot kill an amount bigger than the existing workers")
	}
	return definer.Pool.KillWorkers(amount)
}

func (definer *GoWorkerPoolDefiner) EditWorkersAmount(workersAmount int) error {
	return definer.Pool.SetTotalWorkers(workersAmount)
}

func (definer *GoWorkerPoolDefiner) PauseAllWorkers() {
	definer.Pool.PauseAllWorkers()
}

func (definer *GoWorkerPoolDefiner) ResumeAllWorkers() {
	definer.Pool.ResumeAllWorkers()
}

//WaitWhileAlive wait while there is at least one worker doing some work
func (definer *GoWorkerPoolDefiner) WaitWhileAlive() error{
	return definer.Pool.Wait()
}