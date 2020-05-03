package pool

import (
	"errors"
	"github.com/enriquebris/goworkerpool"
)


type GoWorkerPoolAdapter struct {
	*goworkerpool.Pool
}

//SetWorkerFunc sets the function to be executed by the workers on this pool
func (definer *GoWorkerPoolAdapter) SetWorkerFunc(fn func(interface{})bool) {
	definer.Pool.SetWorkerFunc(fn)
}

//AddTask adds task to be executed
func (definer *GoWorkerPoolAdapter) AddTask(data interface{}) error {
	return definer.Pool.AddTask(data)
}

//AddWorkers adds workers on the fly to the pool
func (definer *GoWorkerPoolAdapter) AddWorkers(amount int) error {
	return definer.Pool.AddWorkers(amount)
}

//KillWorkers kills the desired amount of workers after they  finish their current job
func (definer *GoWorkerPoolAdapter) KillWorkers(amount int) error {
	if amount > definer.GetTotalWorkers() {
		return errors.New("cannot kill an amount bigger than the existing workers")
	}
	return definer.Pool.KillWorkers(amount)
}

//EditWorkersAmount changes the amount of workers on the fly
func (definer *GoWorkerPoolAdapter) EditWorkersAmount(workersAmount int) error {
	return definer.Pool.SetTotalWorkers(workersAmount)
}

//PauseAllWorkers stops the workers from doing any work
func (definer *GoWorkerPoolAdapter) PauseAllWorkers() {
	definer.Pool.PauseAllWorkers()
}

//ResumeAllWorkers puts workers to work after being paused
func (definer *GoWorkerPoolAdapter) ResumeAllWorkers() {
	definer.Pool.ResumeAllWorkers()
}

//Wait wait while there is at least one worker doing some work
func (definer *GoWorkerPoolAdapter) Wait() error{
	return definer.Pool.Wait()
}