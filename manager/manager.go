package manager

import (
	"fmt"
	"github.com/enriquebris/goworkerpool"
	"github.com/pkg/errors"
)

type Manager struct {
	pools map[string]*goworkerpool.Pool
}

//AddPool creates a new pool in the map of pools and returns the success of the operation
func (manager *Manager) AddPool(poolId string, initialWorkers int, maxJobsInQueue int, verbose bool) error{
	if poolId == ""{
		return errors.New("PoolId cannot be empty")
	}
	if maxJobsInQueue < 1{
		return errors.New("maxJobsInQueue has to be greater than 0")
	}
	if _, ok := manager.pools[poolId]; ok {
		return errors.New(fmt.Sprintf("A pool with `%s` id already exist", poolId))
	}
	manager.pools[poolId] = goworkerpool.NewPool(initialWorkers, maxJobsInQueue, verbose)
	return nil
}

//SetFunc defines the function to be executed by an specific pool
func (manager *Manager) SetFunc(poolId string, workerFunc goworkerpool.PoolFunc) error{
	pool, isElementInMap := manager.pools[poolId]
	if !isElementInMap {
		return errors.New(fmt.Sprintf("Pool with `%s` id does not exist", poolId))
	}
	if isElementInMap {
		pool.SetWorkerFunc(workerFunc)
	}
	return nil
}

//AddTaskToPool enqueues a new task to be accomplished by the desired pool
func (manager *Manager) AddTaskToPool(poolId string, data interface{}) error{
	return nil
}

//AddWorkersToPool increments the workers amount in {poolId} by {workersAmount} elements
func (manager *Manager) AddWorkersToPool(poolId string, workersAmount int) error {
	return nil
}

//KillWorkersFromPool decrements the workers amount in {poolId} by {workersAmount} elements
func (manager *Manager) KillWorkersFromPool(poolId string, workersAmount int) error {
	return nil
}

//EditPoolWorkersAmount set a fixed amount {workersAmount} of workers for poolId
func (manager *Manager) EditPoolWorkersAmount(poolId string, workersAmount int) error {
	return nil
}

//PauseWorkersFromPool pause the work for all the workers from {poolId}
func (manager *Manager) PauseWorkersFromPool(poolId string) error {
	return nil
}

//PauseWorkersFromPool resume the works for all the workers from {poolId}
func (manager *Manager) ResumeWorkersFromPool(poolId string) error {
	return nil
}

//WaitForPool blocks while at least a worker from poolId is alive
func (manager *Manager) WaitForPool(poolId string) error {
	return nil
}

//WaitForPool blocks while at least a worker from all the pools is alive
func (manager *Manager) WaitForAllPools() error {
	return nil
}