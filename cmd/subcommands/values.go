package cmd

import (
	"fmt"

	color "github.com/fatih/color"
)

const (
	hmyDocsDir      = "hmycli-docs"
	defaultNodeAddr = "http://localhost:9500"
)

var (
	g           = color.New(color.FgGreen).SprintFunc()
	cookbookDoc = fmt.Sprintf(`
Cookbook of usage, note that every subcommand recognizes a '--help' flag

%s
hmy --node="https://api.s1.b.hmny.io/" --pretty balance <SOME_ONE_ADDRESS>

%s
hmy --node="https://api.s1.b.hmny.io" --pretty blockchain transaction-by-hash <SOME_TRANSACTION_HASH>

%s
hmy keys list
`,
		g("1. Check Balances"),
		g("2. Check local keys"),
		g("3. List local keys"),
	)
)
