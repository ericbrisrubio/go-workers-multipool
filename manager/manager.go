package manager

import (
	"fmt"
	"github.com/enriquebris/goworkerpool"
	"github.com/pkg/errors"
)

//Manager takes care of the different existing pools
type Manager struct {
	pools map[string]*goworkerpool.Pool
}

//AddPool creates a new pool in the map of pools and returns the success of the operation
func (manager *Manager) AddPool(poolID string, initialWorkers int, maxJobsInQueue int, verbose bool) error{
	if poolID == ""{
		return errors.New("PoolId cannot be empty")
	}
	if maxJobsInQueue < 1{
		return errors.New("maxJobsInQueue has to be greater than 0")
	}
	if _, ok := manager.pools[poolID]; ok {
		return errors.New(fmt.Sprintf("A pool with `%s` id already exist", poolID))
	}
	manager.pools[poolID] = goworkerpool.NewPool(initialWorkers, maxJobsInQueue, verbose)
	return nil
}

//SetFunc defines the function to be executed by an specific pool
func (manager *Manager) SetFunc(poolID string, workerFunc goworkerpool.PoolFunc) error{
	pool, isElementInMap := manager.pools[poolID]
	if !isElementInMap {
		return errors.New(fmt.Sprintf("Pool with `%s` id does not exist", poolID))
	}
	if isElementInMap {
		pool.SetWorkerFunc(workerFunc)
	}
	return nil
}

//AddTaskToPool enqueues a new task to be accomplished by the desired pool
func (manager *Manager) AddTaskToPool(poolID string, data interface{}) error{
	return nil
}

//AddWorkersToPool increments the workers amount in {poolID} by {workersAmount} elements
func (manager *Manager) AddWorkersToPool(poolID string, workersAmount int) error {
	return nil
}

//KillWorkersFromPool decrements the workers amount in {poolID} by {workersAmount} elements
func (manager *Manager) KillWorkersFromPool(poolID string, workersAmount int) error {
	return nil
}

//EditPoolWorkersAmount set a fixed amount {workersAmount} of workers for poolID
func (manager *Manager) EditPoolWorkersAmount(poolID string, workersAmount int) error {
	return nil
}

//PauseWorkersFromPool pause the work for all the workers from {poolID}
func (manager *Manager) PauseWorkersFromPool(poolID string) error {
	return nil
}

//ResumeWorkersFromPool resume the works for all the workers from {poolID}
func (manager *Manager) ResumeWorkersFromPool(poolID string) error {
	return nil
}

//WaitForPool blocks while at least a worker from poolID is alive
func (manager *Manager) WaitForPool(poolID string) error {
	return nil
}

//WaitForAllPools blocks while at least a worker from all the pools is alive
func (manager *Manager) WaitForAllPools() error {
	return nil
}