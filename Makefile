  run-pipeline:
	@echo ******RUNNING BUILD******
	go build
	@echo ******MAKING SURE LINT IS CORRECT******
	go get -u golang.org/x/lint/golint
	golint -set_exit_status manager/... ./
	@echo ******STARTING TESTS******
	go test -gcflags=-l ./...
	@echo ******DONE******