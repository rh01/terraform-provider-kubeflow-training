PKG_NAME=kubeflowtraining
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

default: build test

build:
	go install

test:
	go test -v ./...  

testacc:
	TF_ACC=1 KUBEFLOWTRAINING_HOST=http://localhost:8080 go test -v ./... -timeout 120m -covermode=count -coverprofile=coverage.out


#? lint: run a meta linter
lint:
	@hash golangci-lint || (echo "Download golangci-lint from https://github.com/golangci/golangci-lint#install" && exit 1)
	golangci-lint run ./...


fmt:
	gofmt -w $(GOFMT_FILES)