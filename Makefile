# Image to build and push
KO_DOCKER_REPO := rchalumeau/tfmodules
VERSION := $(shell cat VERSION)

# helpers
COMMAND := cmd/tfmodules
PACKAGE := modules

.PHONY: vendor
vendor:
	GO111MODULE=on go mod vendor
	GO111MODULE=on go mod tidy

.PHONY: lint
lint:
	golangci-lint version
	GL_DEBUG=linters_output GO111MODULE=on golangci-lint run

.PHONY: generate
generate:
	# install the generator with go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen
	oapi-codegen \
		--package=${PACKAGE}  \
		--generate=types,chi-server,spec \
		-o pkg/${PACKAGE}/${PACKAGE}.gen.go \
		api/${PACKAGE}.yaml

.PHONY: server
server:
	VERBOSE=1 go run ${COMMAND}/main.go

.PHONY: local
local:
	BACKEND=fake \
	VERBOSE=1 \
	go run ${COMMAND}/main.go

.PHONY: doc
doc:
	openapi-generator generate -i api/modules.yaml -g markdown --skip-validate-spec -o docs

.PHONY: test
test:
	go test -cover ./... -v

.PHONY: prepare-test-module
prepare-test-module:
	tar -czvf pkg/backends/fake/fake_storage/testModule.tar.gz -C test/testModule .

.PHONY: build
build:
	GOFLAGS="-ldflags=-X=main.version=${VERSION}" \
	KO_DOCKER_REPO=${KO_DOCKER_REPO} \
	ko publish ./${COMMAND} --bare

.PHONY: push
push:
	GOFLAGS="-ldflags=-X=main.version=${VERSION}" \
	KO_DOCKER_REPO=${KO_DOCKER_REPO} \
	ko publish ./${COMMAND} --bare --push -t ${VERSION}

