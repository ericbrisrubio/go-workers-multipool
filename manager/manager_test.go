package manager

import (
	"bou.ke/monkey"
	"fmt"
	"github.com/enriquebris/goworkerpool"
	"github.com/pkg/errors"
	"go-workers-multipool/pool"
	"reflect"
	"testing"
)

func TestManager_AddPool(t *testing.T) {
	//type fields struct {
	pools := make(map[string]pool.Descriptor, 1)
	poolsInitializer := make(map[string]int, 1)
	//}
	type args struct {
		poolID         string
		initialWorkers int
		maxJobsInQueue int
		verbose        bool
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			"Adds slow processing queue",
			args{
				poolID:         "slowProcessing",
				initialWorkers: 2,
				maxJobsInQueue: 2,
				verbose:        false,
			},
			nil,
		},
		{
			"Adds a queue with an existing queueId must fail",
			args{
				poolID:         "slowProcessing",
				initialWorkers: 2,
				maxJobsInQueue: 2,
				verbose:        false,
			},
			errors.New("A pool with `slowProcessing` id already exist"),
		},
		{
			"Returns error when pool with empty id tried to be created",
			args{
				poolID:         "",
				initialWorkers: 2,
				maxJobsInQueue: 2,
				verbose:        false,
			},
			errors.New("PoolId cannot be empty"),
		},
		{
			"Returns error when pool with empty tasks queue tried to bre created",
			args{
				poolID:         "processQueue",
				initialWorkers: 2,
				maxJobsInQueue: 0,
				verbose:        false,
			},
			errors.New("maxJobsInQueue has to be greater than 0"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &Manager{
				poolsInitializer: poolsInitializer,
				pools: pools,
			}
			got := manager.AddPool(tt.args.poolID, tt.args.initialWorkers, tt.args.maxJobsInQueue, tt.args.verbose)
			if got == nil && tt.want != nil {
				t.Errorf("AddPool() = %v, want %v", got, tt.want)
			}
			if got != nil && got.Error() != tt.want.Error() {
				t.Errorf("AddPool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_VerifyAddedPool(t *testing.T) {
	newManager := createManagerMock(2)
	if newManager.AddPool("slowProcessing", 2, 2, false) != nil {
		t.Fatalf("AddPool has failed creating the new pool")
	}
	if _, ok := newManager.pools["slowProcessing"]; !ok {
		t.Fail()
	}
	if initialWorkersAmount, ok := newManager.poolsInitializer["slowProcessing"]; !ok && initialWorkersAmount == 0 {
		t.Fail()
	}

}

func TestManager_StartPool(t *testing.T) {
	manager := createManagerMock(1)
	manager.AddPool("slowProcessing", 2, 2, false)
	type fields struct {
		poolsInitializer map[string]int
		pools            map[string]pool.Descriptor
	}
	type args struct {
		poolID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Returns error if pool with {poolId} id does not exist",
			fields{poolsInitializer: make(map[string]int,1), pools: make(map[string]pool.Descriptor,1)},
			args{poolID: "slowProcessing"},
			true,
		},
		{
			"Returns error if difference between pool and poolInitializer",
			fields{poolsInitializer: make(map[string]int,1), pools: manager.pools},
			args{poolID: "slowProcessing"},
			true,
		},
		{
			"Success if pool and poolInitializer are in sync",
			fields{poolsInitializer: manager.poolsInitializer, pools: manager.pools},
			args{poolID: "slowProcessing"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &Manager{
				poolsInitializer: tt.fields.poolsInitializer,
				pools:            tt.fields.pools,
			}
			if err := manager.StartPool(tt.args.poolID); (err != nil) != tt.wantErr {
				t.Errorf("StartPool() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_StartPoolWithEditing(t *testing.T) {
	manager := createManagerMock(1)
	manager.AddPool("slowProcessing", 2, 2, false)
	poolMock := &pool.GoWorkerPoolMock{}
	manager.pools["slowProcessing"] = poolMock
	manager.StartPool("slowProcessing")
	if !poolMock.EditWorkersHasBeenCalled {
		t.Fail()
	}
}

func TestManager_StartPoolWithEditingErr(t *testing.T) {
	manager := createManagerMock(1)
	manager.AddPool("slowProcessing", 2, 2, false)
	poolMock := &pool.GoWorkerPoolMock{}
	monkey.PatchInstanceMethod(reflect.TypeOf(poolMock), "EditWorkersAmount", func(poolMock *pool.GoWorkerPoolMock, amount int) error {
		return errors.New("error from patch")
	})
	defer monkey.UnpatchAll()
	manager.pools["slowProcessing"] = poolMock
	errStart := manager.StartPool("slowProcessing")
	if errStart == nil {
		t.Fail()
	}
}

func TestManager_SetFunc(t *testing.T) {
	manager := createManagerMock(1)
	manager.AddPool("slowProcessing", 2, 2, false)
	type fields struct {
		pools map[string]pool.Descriptor
	}
	type args struct {
		poolID     string
		workerFunc goworkerpool.PoolFunc
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Sets function for an existing pool",
			fields{manager.pools},
			args{"slowProcessing",
				func(data interface{}) bool {
					return true
				},
			},
			false,
		},
		{
			"Fails setting function for a non existing pool",
			fields{make(map[string]pool.Descriptor)},
			args{"slowProcessing",
				func(data interface{}) bool {
					return true
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager = &Manager{
				pools: tt.fields.pools,
			}
			if err := manager.SetFunc(tt.args.poolID, tt.args.workerFunc); (err != nil) != tt.wantErr {
				t.Errorf("SetFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_SetFuncForPool(t *testing.T) {
	manager := createManagerMock(1)
	manager.AddPool("slowProcessing", 2, 2, false)
	poolMock := &pool.GoWorkerPoolMock{}
	manager.pools["slowProcessing"] = poolMock
	manager.SetFunc("slowProcessing", func(i interface{}) bool {
		fmt.Println("executing function")
		return true
	})
	if !poolMock.SetWorkerFuncHasBeenCalled {
		t.Fatal("SetWorkerFunc has not been set correctly")
	}

}

func TestManager_AddTask(t *testing.T) {
	manager := createManagerMock(1)
	manager.AddPool("slowProcessing", 2, 2, false)
	type fields struct {
		pools map[string]pool.Descriptor
	}
	type args struct {
		poolID string
		data   interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Returns error if poolID id is empty",
			fields{pools: make(map[string]pool.Descriptor, 1)},
			args{poolID: "", data: "testdata"},
			true,
		},
		{
			"Returns error if does not exist a pool with poolID id",
			fields{pools: make(map[string]pool.Descriptor, 1)},
			args{poolID: "nonExistingId", data: "testdata"},
			true,
		},
		{
			"Returns error if data is nil",
			fields{pools: manager.pools},
			args{poolID: "slowProcessing", data: nil},
			true,
		},
		{
			"Does not return error if data is as expected",
			fields{pools: manager.pools},
			args{poolID: "slowProcessing", data: struct {
				testField string
			}{
				"filledField",
			}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &Manager{
				pools: tt.fields.pools,
			}
			if err := manager.AddTaskToPool(tt.args.poolID, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("AddTaskToPool() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_AddTaskToPool(t *testing.T) {
	manager := createManagerMock(1)
	manager.AddPool("slowProcessing", 2, 2, false)
	poolMock := &pool.GoWorkerPoolMock{}
	manager.pools["slowProcessing"] = poolMock
	manager.AddTaskToPool("slowProcessing", "task test")
	if !poolMock.AddTaskFuncHasBeenCalled {
		t.Error()
	}
}

func TestManager_AddWorkers(t *testing.T) {
	manager := createManagerMock(1)
	manager.AddPool("slowProcessing", 2, 2, false)
	type fields struct {
		pools map[string]pool.Descriptor
	}
	type args struct {
		poolID        string
		workersAmount int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Returns error if {poolID} id is empty",
			fields{pools: make(map[string]pool.Descriptor, 1)},
			args{poolID: "", workersAmount: 1},
			true,
		},
		{
			"Returns error if the added workers amount is 0",
			fields{pools: make(map[string]pool.Descriptor, 1)},
			args{poolID: "poolTest", workersAmount: 0},
			true,
		},
		{
			"Returns error if pool does not exist with {poolId} id",
			fields{pools: manager.pools},
			args{poolID: "poolTest", workersAmount: 2},
			true,
		},
		{
			"Success if pool with {poolId} id exists and workers amount is >= 1",
			fields{pools: manager.pools},
			args{poolID: "slowProcessing", workersAmount: 2},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &Manager{
				pools: tt.fields.pools,
			}
			if err := manager.AddWorkersToPool(tt.args.poolID, tt.args.workersAmount); (err != nil) != tt.wantErr {
				t.Errorf("AddWorkersToPool() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_AddWorkersToPool(t *testing.T) {
	manager := createManagerMock(1)
	manager.AddPool("slowProcessing", 2, 2, false)
	poolMock := &pool.GoWorkerPoolMock{}
	manager.pools["slowProcessing"] = poolMock
	manager.AddWorkersToPool("slowProcessing", 2)
	if !poolMock.AddWorkersHasBeenCalled {
		t.Fail()
	}
}

func TestManager_KillWorkers(t *testing.T) {
	manager := createManagerMock(1)
	manager.AddPool("slowProcessing", 2, 2, false)
	manager.StartPool("slowProcessing")
	type fields struct {
		pools map[string]pool.Descriptor
		poolsInitializer map[string]int
	}
	type args struct {
		poolID        string
		workersAmount int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Returns error if pool with {poolId} id does not exist",
			fields{pools: manager.pools},
			args{poolID: "", workersAmount: 1},
			true,
		},
		{
			"Returns error if workers amount equals 0",
			fields{pools: manager.pools},
			args{poolID: "slowProcessing", workersAmount: 0},
			true,
		},
		{
			"Returns error if amount to kill > existing workers",
			fields{pools: manager.pools},
			args{poolID: "slowProcessing", workersAmount: 3},
			true,
		},
		{  // Workers need to be started first
			"Success if pool exists and workers amount >= 1",
			fields{pools: manager.pools, poolsInitializer: manager.poolsInitializer},
			args{poolID: "slowProcessing", workersAmount: 2},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &Manager{
				pools: tt.fields.pools,
				poolsInitializer: tt.fields.poolsInitializer,
			}
			if err := manager.KillWorkersFromPool(tt.args.poolID, tt.args.workersAmount); (err != nil) != tt.wantErr {
				t.Errorf("KillWorkersFromPool() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_KillWorkersFromPool(t *testing.T) {
	pools := make(map[string]pool.Descriptor, 1)
	poolsInitializer := make(map[string]int, 1)
	manager := &Manager{
		poolsInitializer: poolsInitializer,
		pools: pools,
	}
	manager.AddPool("slowProcessing", 2, 2, false)
	poolMock := &pool.GoWorkerPoolMock{}
	manager.pools["slowProcessing"] = poolMock
	manager.KillWorkersFromPool("slowProcessing", 2)
	if !poolMock.KillWorkersHasBeenCalled {
		t.Fail()
	}
}

func TestManager_EditWorkersAmount(t *testing.T) {
	pools := make(map[string]pool.Descriptor, 1)
	poolsInitializer := make(map[string]int, 1)
	manager := &Manager{
		poolsInitializer: poolsInitializer,
		pools: pools,
	}
	manager.AddPool("slowProcessing", 2, 2, false)
	type fields struct {
		pools map[string]pool.Descriptor
	}
	type args struct {
		poolID        string
		workersAmount int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Returns error if pool with {poolId} id is not defined",
			fields{pools},
			args{poolID: "slowProcessing1", workersAmount: 2},
			true,
		},
		{
			"Returns error if workers amount is less than 0",
			fields{pools},
			args{poolID: "slowProcessing", workersAmount: -1},
			true,
		},
		{
			"Success if pool is defined and workers amount >= 0",
			fields{pools},
			args{poolID: "slowProcessing", workersAmount: 0},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &Manager{
				pools: tt.fields.pools,
			}
			if err := manager.EditPoolWorkersAmount(tt.args.poolID, tt.args.workersAmount); (err != nil) != tt.wantErr {
				t.Errorf("EditPoolWorkersAmount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_EditWorkersAmountToPool(t *testing.T) {
	pools := make(map[string]pool.Descriptor, 1)
	poolsInitializer := make(map[string]int, 1)
	manager := &Manager{
		poolsInitializer: poolsInitializer,
		pools: pools,
	}
	manager.AddPool("slowProcessing", 2, 2, false)
	poolMock := &pool.GoWorkerPoolMock{}
	manager.pools["slowProcessing"] = poolMock
	manager.EditPoolWorkersAmount("slowProcessing", 5)
	if poolMock.EditWorkersHasBeenCalled == false {
		t.Fail()
	}
}

func TestManager_PauseWorkers(t *testing.T) {
	pools := make(map[string]pool.Descriptor, 1)
	poolsInitializer := make(map[string]int, 1)
	manager := &Manager{
		poolsInitializer: poolsInitializer,
		pools: pools,
	}
	manager.AddPool("slowProcessing", 2, 2, false)
	type fields struct {
		pools map[string]pool.Descriptor
	}
	type args struct {
		poolID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Returns error if pool with {poolId} id not defined",
			fields{pools: pools},
			args{"slowProcessing1"},
			true,
		},
		{
			"Success if pool with {poolId} id is defined",
			fields{pools: pools},
			args{"slowProcessing"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &Manager{
				pools: tt.fields.pools,
			}
			if err := manager.PauseWorkersFromPool(tt.args.poolID); (err != nil) != tt.wantErr {
				t.Errorf("PauseWorkersFromPool() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_PauseWorkersFromPool(t *testing.T) {
	pools := make(map[string]pool.Descriptor, 1)
	poolsInitializer := make(map[string]int, 1)
	manager := &Manager{
		poolsInitializer: poolsInitializer,
		pools: pools,
	}
	manager.AddPool("slowProcessing", 2, 2, false)
	poolMock := &pool.GoWorkerPoolMock{}
	manager.pools["slowProcessing"] = poolMock
	manager.PauseWorkersFromPool("slowProcessing")
	if !poolMock.PauseAllWorkersHasBeenCalled {
		t.Fail()
	}
}

func TestManager_ResumeWorkers(t *testing.T) {
	pools := make(map[string]pool.Descriptor, 1)
	poolsInitializer := make(map[string]int, 1)
	manager := &Manager{
		poolsInitializer: poolsInitializer,
		pools: pools,
	}
	manager.AddPool("slowProcessing", 2, 2, false)
	type fields struct {
		pools map[string]pool.Descriptor
	}
	type args struct {
		poolID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Returns error if pool with {poolId} id not defined",
			fields{pools: pools},
			args{"slowProcessing1"},
			true,
		},
		{
			"Success if pool with {poolId} id is defined",
			fields{pools: pools},
			args{"slowProcessing"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &Manager{
				pools: tt.fields.pools,
			}
			if err := manager.ResumeWorkersFromPool(tt.args.poolID); (err != nil) != tt.wantErr {
				t.Errorf("ResumeWorkersFromPool() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_ResumeWorkersFromPool(t *testing.T) {
	pools := make(map[string]pool.Descriptor, 1)
	poolsInitializer := make(map[string]int, 1)
	manager := &Manager{
		poolsInitializer: poolsInitializer,
		pools: pools,
	}
	manager.AddPool("slowProcessing", 2, 2, false)
	poolMock := &pool.GoWorkerPoolMock{}
	manager.pools["slowProcessing"] = poolMock
	manager.ResumeWorkersFromPool("slowProcessing")
	if !poolMock.ResumeAllWorkersHasBeenCalled {
		t.Fail()
	}
}

func TestManager_Wait(t *testing.T) {
	pools := make(map[string]pool.Descriptor, 1)
	poolsInitializer := make(map[string]int, 1)
	manager := &Manager{
		poolsInitializer: poolsInitializer,
		pools: pools,
	}
	manager.AddPool("slowProcessing", 2, 2, false)
	poolMock := &pool.GoWorkerPoolMock{}
	manager.pools["slowProcessing"] = poolMock
	type fields struct {
		pools map[string]pool.Descriptor
	}
	type args struct {
		poolID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Returns error if pool with {poolId} id not defined",
			fields{pools: pools},
			args{"slowProcessing1"},
			true,
		},
		{
			"Success if pool with {poolId} id is defined",
			fields{pools: pools},
			args{"slowProcessing"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &Manager{
				pools: tt.fields.pools,
			}
			if err := manager.WaitForPool(tt.args.poolID); (err != nil) != tt.wantErr {
				t.Errorf("WaitForPool() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_WaitForPool(t *testing.T) {
	pools := make(map[string]pool.Descriptor, 1)
	poolsInitializer := make(map[string]int, 1)
	manager := &Manager{
		poolsInitializer: poolsInitializer,
		pools: pools,
	}
	manager.AddPool("slowProcessing", 2, 2, false)
	poolMock := &pool.GoWorkerPoolMock{}
	manager.pools["slowProcessing"] = poolMock
	manager.WaitForPool("slowProcessing")
	if !poolMock.WaitHasBeenCalled {
		t.Fail()
	}
}

func TestManager_WaitForAllPools(t *testing.T) {
	type fields struct {
		pools map[string]pool.Descriptor
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"Returns error if no pool is declared",
			fields{make(map[string]pool.Descriptor,2)},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &Manager{
				pools: tt.fields.pools,
			}
			if err := manager.WaitForAllPools(); (err != nil) != tt.wantErr {
				t.Errorf("WaitForAllPools() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_WaitForAllPoolsGroup(t *testing.T) {
	pools := make(map[string]pool.Descriptor, 1)
	poolsInitializer := make(map[string]int, 1)
	manager := &Manager{
		poolsInitializer: poolsInitializer,
		pools: pools,
	}
	manager.AddPool("slowProcessing", 2, 2, false)
	manager.AddPool("fastProcessing", 5, 2, false)
	poolMockSlow := &pool.GoWorkerPoolMock{}
	poolMockFast := &pool.GoWorkerPoolMock{}
	manager.pools["slowProcessing"] = poolMockSlow
	manager.pools["fastProcessing"] = poolMockFast
	if manager.WaitForAllPools() != nil {
		t.Fail()
	}

}

func createManagerMock(elementAmountInMap int) *Manager {
	pools := make(map[string]pool.Descriptor, elementAmountInMap)
	poolsInitializer := make(map[string]int, elementAmountInMap)
	manager := &Manager{
		poolsInitializer: poolsInitializer,
		pools: pools,
	}
	return manager
}