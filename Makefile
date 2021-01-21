APP=proxytest
APP_VERSION:=0.1
ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")

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

