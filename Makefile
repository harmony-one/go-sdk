SHELL := /bin/bash
version := $(shell git rev-list --count HEAD)
commit := $(shell git describe --always --long --dirty)
built_at := $(shell date +%FT%T%z)
built_by := ${USER}@harmony.one

flags := -gcflags="all=-N -l -c 2"
ldflags := -X main.version=v${version} -X main.commit=${commit}
ldflags += -X main.builtAt=${built_at} -X main.builtBy=${built_by}
cli := ./dist/hmy
upload-path-darwin := 's3://pub.harmony.one/release/darwin-x86_64/mainnet/hmy'
upload-path-darwin-version := 's3://pub.harmony.one/release/darwin-x86_64/mainnet/hmy_version'
upload-path-linux := 's3://pub.harmony.one/release/linux-x86_64/mainnet/hmy'
upload-path-linux-version := 's3://pub.harmony.one/release/linux-x86_64/mainnet/hmy_version'

env := GO111MODULE=on

DIR := ${CURDIR}
export CGO_LDFLAGS=-L$(DIR)/dist/lib -Wl,-rpath -Wl,\$ORIGIN/lib

all:
	source $(shell go env GOPATH)/src/github.com/harmony-one/harmony/scripts/setup_bls_build_flags.sh && $(env) go build -o $(cli) -ldflags="$(ldflags)" cmd/main.go
	cp $(cli) hmy

static:
	make -C $(shell go env GOPATH)/src/github.com/harmony-one/mcl
	make -C $(shell go env GOPATH)/src/github.com/harmony-one/bls minimised_static BLS_SWAP_G=1
	source $(shell go env GOPATH)/src/github.com/harmony-one/harmony/scripts/setup_bls_build_flags.sh && $(env) go build -o $(cli) -ldflags="$(ldflags) -w -extldflags \"-static\"" cmd/main.go
	cp $(cli) hmy

debug:
	source $(shell go env GOPATH)/src/github.com/harmony-one/harmony/scripts/setup_bls_build_flags.sh && $(env) go build $(flags) -o $(cli) -ldflags="$(ldflags)" cmd/main.go
	cp $(cli) hmy

install:all
	cp $(cli) ~/.local/bin

run-tests: test-rpc test-key;

test-key:
	go test ./pkg/keys -cover -v

test-rpc:
	go test ./pkg/rpc -cover -v

# Notice assumes you have correct uploading credentials
upload-darwin:all
	aws --profile upload s3 cp ./hmy ${upload-path-darwin}
	./hmy version &> ./hmy_version
	aws --profile upload s3 cp ./hmy_version ${upload-path-darwin-version}

upload-linux:static
	aws --profile upload s3 cp ./hmy ${upload-path-linux}
	./hmy version &> ./hmy_version
	aws --profile upload s3 cp ./hmy_version ${upload-path-linux-version}

.PHONY:clean run-tests upload-darwin upload-linux

clean:
	@rm -f $(cli)
	@rm -rf ./dist
