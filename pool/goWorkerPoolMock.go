package pool

type GoWorkerPoolMock struct {
	totalWorkers                  int
	SetWorkerFuncHasBeenCalled    bool
	AddTaskFuncHasBeenCalled      bool
	AddWorkersHasBeenCalled       bool
	KillWorkersHasBeenCalled      bool
	EditWorkersHasBeenCalled      bool
	PauseAllWorkersHasBeenCalled  bool
	ResumeAllWorkersHasBeenCalled bool
	WaitHasBeenCalled             bool
}

func (definer *GoWorkerPoolMock) SetWorkerFunc(fn func(interface{}) bool) {
	definer.SetWorkerFuncHasBeenCalled = true
}

func (definer *GoWorkerPoolMock) AddTask(data interface{}) error {
	definer.AddTaskFuncHasBeenCalled = true
	return nil
}

func (definer *GoWorkerPoolMock) AddWorkers(amount int) error {
	definer.AddWorkersHasBeenCalled = true
	return nil
}

func (definer *GoWorkerPoolMock) KillWorkers(amount int) error {
	definer.KillWorkersHasBeenCalled = true
	return nil
}

func (definer *GoWorkerPoolMock) EditWorkersAmount(workersAmount int) error {
	definer.EditWorkersHasBeenCalled = true
	return nil
}

func (definer *GoWorkerPoolMock) PauseAllWorkers() {
	definer.PauseAllWorkersHasBeenCalled = true
}

func (definer *GoWorkerPoolMock) ResumeAllWorkers() {
	definer.ResumeAllWorkersHasBeenCalled = true
}

func (definer *GoWorkerPoolMock) Wait() error {
	definer.WaitHasBeenCalled = true
	return nil
}
