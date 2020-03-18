PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
GOBASE := $(shell pwd)
GOPATH := $(GOBASE)/vendor:$(GOBASE)
GOBIN := $(GOBASE)/bin

clean:
	@-rm -fR $(GOBIN)
	@echo "  >  Cleaning build cache"
    @GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

build:
	@echo "  >  Building binary..."
	@go build -o $(GOBIN)/calendar_api        $(GOBASE)/cmd/calendar_api/calendar_api.go
	@go build -o $(GOBIN)/calendar_api_client $(GOBASE)/cmd/calendar_api_client/calendar_api_client.go
	@go build -o $(GOBIN)/calendar_scheduler  $(GOBASE)/cmd/calendar_scheduler/calendar_scheduler.go
	@go build -o $(GOBIN)/calendar_sender     $(GOBASE)/cmd/calendar_sender/calendar_sender.go

test:
	@echo "  >  Making unit tests"
	@go test ./...
