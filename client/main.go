package main

import (
	"fmt"

	"github.com/harmony-one/go-sdk/client/cmd"
)

var (
	version string
	commit  string
	builtAt string
	builtBy string
)

func main() {
	fmt.Printf("%s %s %s %s", version, commit, builtAt, builtBy)
	cmd.Execute()
}
