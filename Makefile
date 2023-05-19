.PHONY: build run lint testacc test clean website-serve website-build release

vars:=$(shell test -f .env && grep -v '^\#' .env | xargs)
VERSION:=v0.0.0-local

#? build: build binary for current system
build: bin/current_system/terraform-provider-kubeflow-training_$(VERSION)

#? run: run plugin
run: build
	bin/current_system/terraform-provider-kubeflow-training

#? lint: run a meta linter
lint:
	@hash golangci-lint || (echo "Download golangci-lint from https://github.com/golangci/golangci-lint#install" && exit 1)
	golangci-lint run kubeflow-training/...

#? testacc: run acceptance tests
testacc:
	$(vars) TF_ACC=1 go test -cover -v -timeout 30m -failfast ./kubeflow-training

#? test: run unit tests
test:
	$(vars) go test $(shell go list ./... | grep -v /\.tmp/ ) -v

.env:
	cp .env.example .env

#? clean: removes all artificats
clean:
	rm -fr bin/ .tmp/ website/public website/resources

#? website-serve: serve documentation website
website-serve:
	@hugo -s website server

#? website-build: build website
website-build:
	@hugo -s website

bin/current_system/terraform-provider-kubeflow-training_%:  GOARGS =
bin/darwin_amd64/terraform-provider-kubeflow-training_%:  GOARGS = GOOS=darwin GOARCH=amd64
bin/linux_amd64/terraform-provider-kubeflow-training_%:  GOARGS = GOOS=linux GOARCH=amd64
bin/linux_386/terraform-provider-kubeflow-training_%:  GOARGS = GOOS=linux GOARCH=386
bin/linux_arm/terraform-provider-kubeflow-training_%:  GOARGS = GOOS=linux GOARCH=arm
bin/windows_amd64/terraform-provider-kubeflow-training_%:  GOARGS = GOOS=windows GOARCH=amd64
bin/windows_386/terraform-provider-kubeflow-training_%:  GOARGS = GOOS=windows GOARCH=386

bin/%/terraform-provider-kubeflow-training_$(VERSION): clean
	$(GOARGS) CGO_ENABLED=0 go build -o $@ -ldflags="-s -w" .

#? release: make a release for all systems
release: \
	bin/release/terraform-provider-kubeflow-training_darwin_amd64.zip \
	bin/release/terraform-provider-kubeflow-training_linux_amd64.zip \
	bin/release/terraform-provider-kubeflow-training_linux_386.zip \
	bin/release/terraform-provider-kubeflow-training_linux_arm.zip \
	bin/release/terraform-provider-kubeflow-training_windows_amd64.zip \
	bin/release/terraform-provider-kubeflow-training_windows_386.zip

bin/release/terraform-provider-kubeflow-training_%.zip: NAME=terraform-provider-kubeflow-training_$(VERSION)_$*
bin/release/terraform-provider-kubeflow-training_%.zip: DEST=bin/release/$(VERSION)/$(NAME)
bin/release/terraform-provider-kubeflow-training_%.zip: bin/%/terraform-provider-kubeflow-training_$(VERSION)
	mkdir -p $(DEST)
	cp bin/$*/terraform-provider-kubeflow-training_$(VERSION) readme.md $(DEST)
	cd $(DEST) && zip -r ../$(NAME).zip . && cd .. && sha256sum $(NAME).zip > $(NAME).sha256 && rm -rf $(NAME)

#? help: display help
help: Makefile
	@printf "Available make targets:\n\n"
	@sed -n 's/^#?//p' $< | column -t -s ':' |  sed -e 's/^/ /'