package manager

import (
	"fmt"
	"github.com/enriquebris/goworkerpool"
	"github.com/pkg/errors"
	"go-workers-multipool/pool"
	"strings"
	"sync"
)

//Manager takes care of the different existing pools
type Manager struct {
	poolsInitializer map[string]int
	pools map[string]pool.Descriptor
}

//AddPool creates a new pool in the map of pools and returns the success of the operation
func (manager *Manager) AddPool(poolID string, initialWorkers int, maxJobsInQueue int, verbose bool) error {
	if poolID == "" || strings.Trim(poolID, " ") == "" {
		return errors.New("PoolId cannot be empty")
	}
	if maxJobsInQueue < 1 {
		return errors.New("maxJobsInQueue has to be greater than 0")
	}
	if manager.isPoolDefined(poolID) {
		return errors.New(fmt.Sprintf("A pool with `%s` id already exist", poolID))
	}
	manager.pools[poolID] = &pool.GoWorkerPoolAdapter{Pool: goworkerpool.NewPool(0, maxJobsInQueue, verbose)}
	manager.poolsInitializer[poolID] = initialWorkers
	return nil
}

//StartPool makes the workers to start taking care of jobs
func (manager *Manager) StartPool(poolID string) error {
	if !manager.isPoolDefined(poolID) {
		return errors.New(fmt.Sprintf("Pool with `%s` id does not exist", poolID))
	}
	if value, ok := manager.poolsInitializer[poolID]; ok {
		pool, _ := manager.pools[poolID]
		errEditing := pool.EditWorkersAmount(value)
		if errEditing != nil {
			return errEditing
		}
	} else {
		return errors.New(fmt.Sprintf("error initializing pool with id `%s`", poolID))
	}
	return nil
}

//SetFunc defines the function to be executed by an specific pool
func (manager *Manager) SetFunc(poolID string, workerFunc func(interface{}) bool) error {
	if manager.isPoolDefined(poolID) {
		pool, _ := manager.pools[poolID]
		pool.SetWorkerFunc(workerFunc)
	} else {
		return errors.New(fmt.Sprintf("Pool with `%s` id does not exist", poolID))
	}
	return nil
}

//AddTaskToPool enqueues a new task to be accomplished by the desired pool
func (manager *Manager) AddTaskToPool(poolID string, data interface{}) error {
	if data == nil {
		return errors.New("data cannot be nil")
	}
	if manager.isPoolDefined(poolID) {
		pool, _ := manager.pools[poolID]
		pool.AddTask(data)
	} else {
		return errors.New(fmt.Sprintf("No pool exists for poolID: %s", poolID))
	}
	return nil
}

//AddWorkersToPool increments the workers amount in {poolID} by {workersAmount} elements
func (manager *Manager) AddWorkersToPool(poolID string, amount int) error {
	if amount == 0 {
		return errors.New("amount cannot be 0")
	}
	if !manager.isPoolDefined(poolID) {
		return errors.New(fmt.Sprintf("No pool exists for poolID: %s", poolID))
	}
	pool, _ := manager.pools[poolID]
	return pool.AddWorkers(amount)
}

//KillWorkersFromPool decrements the workers amount in {poolID} by {workersAmount} elements
func (manager *Manager) KillWorkersFromPool(poolID string, amount int) error {
	if !manager.isPoolDefined(poolID) {
		return errors.New(fmt.Sprintf("No pool id defined for %s id", poolID))
	}
	if amount == 0 {
		return errors.New("Workers amount cannot be 0")
	}

	pool, _ := manager.pools[poolID]
	return pool.KillWorkers(amount)
}

//EditPoolWorkersAmount set a fixed amount {workersAmount} of workers for poolID
func (manager *Manager) EditPoolWorkersAmount(poolID string, amount int) error {
	if amount < 0 {
		return errors.New("amount has to be greater or equal to 0")
	}
	if !manager.isPoolDefined(poolID) {
		return errors.New(fmt.Sprintf("pool with %s id is not defined", poolID))
	}
	pool, _ := manager.pools[poolID]
	return pool.EditWorkersAmount(amount)
}

//PauseWorkersFromPool pause the work for all the workers from {poolID}
func (manager *Manager) PauseWorkersFromPool(poolID string) error {
	if !manager.isPoolDefined(poolID) {
		return errors.New(fmt.Sprintf("pool with %s id is not defined", poolID))
	}
	pool, _ := manager.pools[poolID]
	pool.PauseAllWorkers()
	return nil
}

//ResumeWorkersFromPool resume the works for all the workers from {poolID}
func (manager *Manager) ResumeWorkersFromPool(poolID string) error {
	if !manager.isPoolDefined(poolID) {
		return errors.New(fmt.Sprintf("pool with %s id is not defined", poolID))
	}
	pool, _ := manager.pools[poolID]
	pool.ResumeAllWorkers()
	return nil
}

//WaitForPool blocks while at least a worker from poolID is alive
func (manager *Manager) WaitForPool(poolID string) error {
	if !manager.isPoolDefined(poolID) {
		return errors.New(fmt.Sprintf("pool with %s id is not defined", poolID))
	}
	pool, _ := manager.pools[poolID]
	pool.Wait()
	return nil
}

//WaitForAllPools blocks while at least a worker from all the pools is alive
func (manager *Manager) WaitForAllPools() error {
	if len(manager.pools) == 0 {
		return errors.New("No pool has been declared")
	}
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(len(manager.pools))
	for _,poolValue := range manager.pools{
		go func(wg *sync.WaitGroup, pool pool.Descriptor) {
			pool.Wait()
			wg.Done()
		}(waitGroup, poolValue)
	}
	waitGroup.Wait()
	return nil
}

func (manager *Manager) isPoolDefined(poolID string) bool {
	_, isElementInMap := manager.pools[poolID]
	return isElementInMap
}
