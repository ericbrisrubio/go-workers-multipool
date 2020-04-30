package manager

import (
	"fmt"
	"github.com/enriquebris/goworkerpool"
	"github.com/pkg/errors"
	"go-workers-multipool/pool"
	"testing"
)

func TestManager_AddPool(t *testing.T) {
	//type fields struct {
		pools := make(map[string]pool.Descriptor,1)
	//}
	type args struct {
		poolID         string
		initialWorkers int
		maxJobsInQueue int
		verbose        bool
	}
	tests := []struct {
		name   string
		args   args
		want   error
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
				pools: pools,
			}
			got := manager.AddPool(tt.args.poolID, tt.args.initialWorkers, tt.args.maxJobsInQueue, tt.args.verbose);
			if got == nil && tt.want != nil {
				t.Errorf("AddPool() = %v, want %v", got, tt.want)
			}
			if  got != nil && got.Error() != tt.want.Error() {
				t.Errorf("AddPool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_VerifyAddedPool(t *testing.T) {
	newManager := Manager{
					pools: make(map[string]pool.Descriptor,1),
				}
	if newManager.AddPool("slowProcessing", 2, 2, false) != nil {
		t.Fatalf("AddPool has failed creating the new pool")
	}
	if _, ok := newManager.pools["slowProcessing"]; !ok{
		t.Fail()
	}
}

func TestManager_SetFunc(t *testing.T) {
	pools := make(map[string]pool.Descriptor,1)
	manager := &Manager{
		pools: pools,
	}
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
		fields fields
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
	pools := make(map[string]pool.Descriptor,1)
	manager := &Manager{
		pools: pools,
	}
	manager.AddPool("slowProcessing", 2, 2, false)
    poolMock := &pool.GoWorkerPoolDefinerMock{}
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
	pools := make(map[string]pool.Descriptor,1)
	manager := &Manager{
		pools: pools,
	}
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
			fields{pools:make(map[string]pool.Descriptor,1)},
			args{poolID:"",data:"testdata"},
			true,
		},
		{
			"Returns error if does not exist a pool with poolID id",
			fields{pools:make(map[string]pool.Descriptor,1)},
			args{poolID:"nonExistingId",data:"testdata"},
			true,
		},
		{
			"Returns error if data is nil",
			fields{pools:manager.pools},
			args{poolID:"slowProcessing",data:nil},
			true,
		},
		{
			"Does not return error if data is as expected",
			fields{pools:manager.pools},
			args{poolID:"slowProcessing",data: struct {
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
	pools := make(map[string]pool.Descriptor,1)
	manager := &Manager{
		pools: pools,
	}
	manager.AddPool("slowProcessing", 2, 2, false)
	poolMock := &pool.GoWorkerPoolDefinerMock{}
	manager.pools["slowProcessing"] = poolMock
	manager.AddTaskToPool("slowProcessing", "task test")
	if !poolMock.AddTaskFuncHasBeenCalled {
		t.Error()
	}
}

func TestManager_AddWorkers(t *testing.T) {
	pools := make(map[string]pool.Descriptor,1)
	manager := &Manager{
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
			fields{pools: pools},
			args{poolID: "poolTest", workersAmount: 2},
			true,
		},
		{
			"Success if pool with {poolId} id exists and workers amount is >= 1",
			fields{pools: pools},
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
	pools := make(map[string]pool.Descriptor,1)
	manager := &Manager{
		pools: pools,
	}
	manager.AddPool("slowProcessing", 2, 2, false)
	poolMock := &pool.GoWorkerPoolDefinerMock{}
	manager.pools["slowProcessing"] = poolMock
	manager.AddWorkersToPool("slowProcessing", 2)
	if !poolMock.AddWorkersHasBeenCalled {
		t.Fail()
	}
}

func TestManager_KillWorkers(t *testing.T) {
	pools := make(map[string]pool.Descriptor,1)
	manager := &Manager{
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
			"Returns error if pool with {poolId} id does not exist",
			fields{pools: pools},
			args{poolID: "", workersAmount: 1},
			true,
		},
		{
			"Returns error if workers amount equals 0",
			fields{pools: pools},
			args{poolID: "slowProcessing", workersAmount: 0},
			true,
		},
		{
			"Success if pool exists and workers amount >= 1",
			fields{pools: pools},
			args{poolID: "slowProcessing", workersAmount: 2},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &Manager{
				pools: tt.fields.pools,
			}
			if err := manager.KillWorkersFromPool(tt.args.poolID, tt.args.workersAmount); (err != nil) != tt.wantErr {
				t.Errorf("KillWorkersFromPool() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_KillWorkersFromPool(t *testing.T) {
	pools := make(map[string]pool.Descriptor,1)
	manager := &Manager{
		pools: pools,
	}
	manager.AddPool("slowProcessing", 2, 2, false)
	poolMock := &pool.GoWorkerPoolDefinerMock{}
	manager.pools["slowProcessing"] = poolMock
	manager.KillWorkersFromPool("slowProcessing", 2)
	if !poolMock.KillWorkersHasBeenCalled {
		t.Fail()
	}
}