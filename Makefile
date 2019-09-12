version := $(shell git rev-list --count HEAD)
commit := $(shell git describe --always --long --dirty)
built_at := $(shell date +%FT%T%z)
built_by := ${USER}@harmony.one

flags := -gcflags="all=-N -l -c 2"
ldflags := -X main.version=v${version} -X main.commit=${commit}
ldflags += -X main.builtAt=${built_at} -X main.builtBy=${built_by}
cli := ./dist/hmy_cli

env := GO111MODULE=on

DIR := ${CURDIR}
export CGO_LDFLAGS=-L$(DIR)/dist/lib -Wl,-rpath -Wl,\$ORIGIN/lib

all:prepare-dirs
	$(env) go build -o $(cli) -ldflags="$(ldflags)" cmd/main.go
	cp $(cli) hmy

debug:prepare-dirs
	$(env) go build $(flags) -o $(cli) -ldflags="$(ldflags)" cmd/main.go
	cp $(cli) hmy


run-tests: test-rpc test-key;

test-key:
	go test ./pkg/keys -cover -v

test-rpc:
	go test ./pkg/rpc -cover -v

prepare-dirs:
	mkdir -p dist
	rsync -a $(shell go env GOPATH)/src/github.com/harmony-one/bls/lib/* ./dist/lib/
	rsync -a $(shell go env GOPATH)/src/github.com/harmony-one/mcl/lib/* ./dist/lib/
	rsync -a /usr/local/opt/openssl/lib/* ./dist/lib/

.PHONY:clean run-tests

clean:
	@rm -f $(cli)
	@rm -rf ./dist
