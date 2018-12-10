GOBIN := $(CURDIR)/bin/
ENV:=GOBIN=$(GOBIN)
DIR:=FILE_DIR=$(CURDIR)/testfiles
GODEBUG:=GODEBUG=gocacheverify=1

##
## List of commands:
##

## default:
all: mod deps fmt lint test run-server

all-deps: mod run-server

tests: fmt deps lint test

test: test-filter test-scalable

deps:
	@echo "======================================================================"
	@echo 'MAKE: install...'
	@mkdir -p $(GOBIN)
	@$(ENV) go get -u golang.org/x/lint/golint

bench:
	@echo "======================================================================"
	@echo "Run bench test for ./bloomfilter/"
	@$(DIR) cd ./bloomfilter/ && $(DIR) go test -bench=.


test-scalable:
	@echo "======================================================================"
	@echo "Run race test for ./bloomfilter/scalable"
	@$(DIR) $(GODEBUG) go test -cover -race ./bloomfilter/scalable/

test-filter:
	@echo "======================================================================"
	@echo "Run race test for ./bloomfilter/"
	@$(DIR) $(GODEBUG) go test -cover -race ./bloomfilter/

lint:
	@echo "======================================================================"
	@echo "Run golint..."
	$(GOBIN)golint ./*.go
	$(GOBIN)golint ./bloomfilter/array/*.go
	$(GOBIN)golint ./bloomfilter/scalable/*.go
	$(GOBIN)golint ./bloomfilter/*.go

fmt:
	@echo "======================================================================"
	@echo "Run go fmt..."
	@go fmt ./*.go
	@go fmt ./bloomfilter/array/*.go
	@go fmt ./bloomfilter/scalable/*.go
	@go fmt ./bloomfilter/*.go

mod:
	@echo "======================================================================"
	@echo "Run MOD"
	@go mod verify
	@go mod tidy -v
	@go mod vendor -v
	@go mod download
	@go mod verify
