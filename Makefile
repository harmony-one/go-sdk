version=$(shell git rev-list --count HEAD)
commit=$(shell git describe --always --long --dirty)
built_at=$(shell date +%FT%T%z)
built_by=${USER}

flags := -gcflags="all=-N -l -c 2"
cli := hmy_cli

all:
	printf '%s %s %s %s\n' $(version) $(commit) $(built_at) $(built_by)
	go build $(flags) -o $(cli) client/main.go

clean:
	@rm -f $(cli)
