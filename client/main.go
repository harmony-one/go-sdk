package main

import (
	"fmt"

	"github.com/harmony-one/go-sdk/common"
)

var (
	version string
	commit  string
	builtAt string
	builtBy string
)

func main() {
	fmt.Printf("%s %s %s %s", version, commit, builtAt, builtBy)
	common.Speak()
}
