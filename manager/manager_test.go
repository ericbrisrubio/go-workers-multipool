package manager

import (
	"github.com/enriquebris/goworkerpool"
	"github.com/pkg/errors"
	"testing"
)

func TestManager_AddPool(t *testing.T) {
	//type fields struct {
		pools := make(map[string]*goworkerpool.Pool,1)
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
					pools: make(map[string]*goworkerpool.Pool,1),
				}
	if newManager.AddPool("slowProcessing", 2, 2, false) != nil {
		t.Fatalf("AddPool has failed creating the new pool")
	}
	if _, ok := newManager.pools["slowProcessing"]; !ok{
		t.Fail()
	}
}

func TestManager_SetFunc(t *testing.T) {
	pools := make(map[string]*goworkerpool.Pool,1)
	manager := &Manager{
		pools: pools,
	}
	manager.AddPool("slowProcessing", 2, 2, false)
	type fields struct {
		pools map[string]*goworkerpool.Pool
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
			fields{make(map[string]*goworkerpool.Pool)},
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

func TestManager_SetFuncExecution(t *testing.T) {
	pools := make(map[string]*goworkerpool.Pool,1)
	manager := &Manager{
		pools: pools,
	}
	manager.AddPool("slowProcessing", 2, 2, false)

}