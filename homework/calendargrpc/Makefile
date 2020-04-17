PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
GOBASE := $(shell pwd)
GOPATH := $(GOBASE)/vendor:$(GOBASE)
GOBIN := $(GOBASE)/bin

.PHONY: clean
clean:
	@-rm -fR $(GOBIN)
	@-rm -fR $(GOBASE)/logs
	@echo "  >  Cleaning build cache"
    @GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

.PHONY: build
build:
	@echo "  >  Building binary..."
	@go build -o $(GOBIN)/calendar_api        $(GOBASE)/cmd/calendar_api/calendar_api.go
	@go build -o $(GOBIN)/calendar_api_client $(GOBASE)/cmd/calendar_api_client/calendar_api_client.go
	@go build -o $(GOBIN)/calendar_scheduler  $(GOBASE)/cmd/calendar_scheduler/calendar_scheduler.go
	@go build -o $(GOBIN)/calendar_sender     $(GOBASE)/cmd/calendar_sender/calendar_sender.go

.PHONY: test
test:
	@echo "  >  Making integration tests"
	set -e ; \
	docker-compose -f deployments/docker-compose.yml -f deployments/docker-compose.test.yml up --build -d ; \
	sleep 10 ; \
	exitCode=0 ; \
	docker-compose -f deployments/docker-compose.yml -f deployments/docker-compose.test.yml \
		run -e CGO_ENABLED=0 -e GOOS=linux integration_tests go test || exitCode=$$? ; \
	docker-compose -f deployments/docker-compose.yml -f deployments/docker-compose.test.yml down ; \
	exit $$exitCode

.PHONY: docker-build
docker-build:
	@echo "  >  Building docker images..."
	@docker build -t deployments_builder:latest  $(GOBASE)/build/package/builder
	@docker build -t deployments_api:latest  $(GOBASE)/build/package/api
	@docker build -t deployments_scheduler:latest  $(GOBASE)/build/package/scheduler
	@docker build -t deployments_sender:latest  $(GOBASE)/build/package/sender

.PHONY: monitor-up
monitor-up:
	@echo "  >  Starting service monitoring"
	@docker-compose -f deployments/docker-compose.yml -f deployments/docker-compose.monitor.yml up -d

.PHONY: monitor-down
monitor-down:
	@echo "  >  Stopping service monitoring"
	@docker-compose -f deployments/docker-compose.yml -f deployments/docker-compose.monitor.yml down