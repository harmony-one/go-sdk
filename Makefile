version := $(shell git rev-list --count HEAD)
commit := $(shell git describe --always --long --dirty)
built_at := $(shell date +%FT%T%z)
built_by := ${USER}@harmony.one

flags := -gcflags="all=-N -l -c 2"
ldflags := -X main.version=v${version} -X main.commit=${commit}
ldflags += -X main.builtAt=${built_at} -X main.builtBy=${built_by}
cli := hmy_cli

all:
	go build -o $(cli) -ldflags="$(ldflags)" client/main.go

debug:
	go build $(flags) -o $(cli) -ldflags="$(ldflags)" client/main.go

run: all
	./$(cli)

.PHONY:clean

clean:
	@rm -f $(cli)
