.PHONY: build build_proto fmt lint vet format build_run_grpc test mock

GO_PACKAGES = $(shell go list ./... )
GO_FILES = $(shell find . -name "*.go" | uniq)

build:
	go build

lint:
	@go list ./...  | xargs -L1 golint
fmt:
	@goimports -w $(GO_FILES)
vet:
	@go vet $(GO_PACKAGES)

format: fmt lint vet

test:
	go test $(shell go list ./... | grep -v -e /example -e /runner/* -e /mock )

build_example:
	mkdir -p remove && find ./example  -name "*.go" -print0 |xargs -0 -I {} -n1 go build -o remove/{} {} && rm -rf remove

mock:
	rm -rf mock
	mkdir -p mock
	mockgen -source taskor.go -destination mock/taskor_mock.go

	rm -rf runner/mock
	mkdir runner/mock
	mockgen --package mock github.com/scaleway/taskor/runner  Runner > runner/mock/runner.go
