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
	@echo "  >  Making unit tests"
	@go test ./...

.PHONY: docker-build
docker-build:
	@echo "  >  Building docker images..."
	@docker build -t otus-go/calendar_api:1.0  $(GOBASE)/build/package/api
	@docker build -t otus-go/calendar_scheduler:1.0  $(GOBASE)/build/package/scheduler
	@docker build -t otus-go/calendar_sender:1.0  $(GOBASE)/build/package/sender
