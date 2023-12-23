OS=$(shell go env GOOS)
ARCH=$(shell go env GOARCH)
FILENAME=aoc-$(OS)-$(ARCH)
FILE_LOCATION=./bin/ng-dfs-notifier

build:
	CGO_ENABLED=1 go build -o ./bin/ng-dfs-notifier ./main.go

buildWithArch:
	go build -o $(FILE_LOCATION) ./

buildLocation:
	@echo $(FILE_LOCATION)

buildName:
	@echo $(FILENAME)

test:
	go test ./...