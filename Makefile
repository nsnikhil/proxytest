APP=proxytest
APP_VERSION:=0.1
APP_EXECUTABLE="./out/$(APP)"
ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")

LOCAL_CONFIG_FILE=local.env
HTTP_SERVE_COMMAND=http-serve

deps:
	go mod download

tidy:
	go mod tidy

check: fmt vet

fmt:
	go fmt $(ALL_PACKAGES)

vet:
	go vet $(ALL_PACKAGES)

copy-config:
	cp .env.sample local.env

test:
	go clean -testcache
	go test ./...

test-cover-html:
	go clean -testcache
	mkdir -p out/
	go test ./... -coverprofile=out/coverage.out
	go tool cover -html=out/coverage.out

compile:
	mkdir -p out/
	go build -o $(APP_EXECUTABLE) cmd/*.go

build: deps compile

local-http-serve: build
	$(APP_EXECUTABLE) -configFile=$(LOCAL_CONFIG_FILE) $(HTTP_SERVE_COMMAND)

http-serve: build
	$(APP_EXECUTABLE) -configFile=$(configFile) $(HTTP_SERVE_COMMAND)

clean:
	rm -rf out/

