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
export CGO_LDFLAGS=-L$(DIR)/lib -Wl,-rpath -Wl,\$ORIGIN/lib
#rsync ./lib/. ./dist/

all:
	mkdir dist
	rsync -a ./lib/* ./dist/lib/
	$(env) go build -o $(cli) -ldflags="$(ldflags)" client/main.go

debug:
	$(env) go build $(flags) -o $(cli) -ldflags="$(ldflags)" client/main.go

run-tests: test-rpc test-key;

test-key:
	go test ./pkg/keys -cover -v

test-rpc:
	go test ./pkg/rpc -cover -v

.PHONY:clean run-tests

clean:
	@rm -f $(cli)
	@rm -rf ./dist
