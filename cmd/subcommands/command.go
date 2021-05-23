package cmd

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/harmony-one/go-sdk/pkg/console"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
)

func init() {
	shardNum := 0
	net := "mainnet"

	cmdCommand := &cobra.Command{
		Use:   "command",
		Short: "Start an interactive JavaScript environment (connect to node)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return openConsole(net, shardNum)
		},
	}

	cmdCommand.Flags().IntVar(&shardNum, "shard", 0, "used shard(default: 0)")
	cmdCommand.Flags().StringVar(&net, "net", "mainnet", "used net(default: mainnet, choose: mainnet or testnet)")

	RootCmd.AddCommand(cmdCommand)
}

func checkAndMakeDirIfNeeded() string {
	userDir, _ := homedir.Dir()
	hmyCLIDir := path.Join(userDir, ".hmy_cli", "command")
	if _, err := os.Stat(hmyCLIDir); os.IsNotExist(err) {
		// Double check with Leo what is right file persmission
		os.Mkdir(hmyCLIDir, 0700)
	}

	return hmyCLIDir
}

// remoteConsole will connect to a remote geth instance, attaching a JavaScript
// console to it.
func openConsole(net string, shard int) error {
	var netMap = map[string]string{
		"mainnet-0": "https://api.harmony.one",
		"mainnet-1": "https://s1.api.harmony.one",
		"mainnet-2": "https://s2.api.harmony.one",
		"mainnet-3": "https://s3.api.harmony.one",
		"testnet-0": "https://api.s0.b.hmny.io",
		"testnet-1": "https://api.s1.b.hmny.io",
		"testnet-2": "https://api.s2.b.hmny.io",
		"testnet-3": "https://api.s3.b.hmny.io",
	}

	endpoint := ""
	if findedEndpoint, ok := netMap[fmt.Sprintf("%s-%d", net, shard)]; ok {
		endpoint = findedEndpoint
	} else {
		log.Fatalf("Unknown network `%s` or shardId `%d`", net, shard)
	}

	client, err := rpc.Dial(endpoint)
	if err != nil {
		log.Fatalf("Unable to attach to remote geth: %v", err)
	}
	config := console.Config{
		DataDir: checkAndMakeDirIfNeeded(),
		DocRoot: ".",
		Client:  client,
		Preload: nil,
		NodeUrl: endpoint,
		ShardId: shard,
		Net:     net,
	}

	consoleInstance, err := console.New(config)
	if err != nil {
		log.Fatalf("Failed to start the JavaScript console insatnce: %v", err)
	}
	defer consoleInstance.Stop(false)

	// Otherwise print the welcome screen and enter interactive mode
	consoleInstance.Welcome()
	consoleInstance.Interactive()

	return nil
}
