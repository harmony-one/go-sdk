package console

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	ethereum_rpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/harmony-one/go-sdk/pkg/account"
	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/console/jsre"
	"github.com/harmony-one/go-sdk/pkg/console/jsre/deps"
	"github.com/harmony-one/go-sdk/pkg/console/prompt"
	"github.com/harmony-one/go-sdk/pkg/console/web3ext"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/go-sdk/pkg/store"
	"github.com/harmony-one/go-sdk/pkg/transaction"
	"github.com/harmony-one/harmony/accounts"

	"github.com/dop251/goja"
	"github.com/mattn/go-colorable"
	"github.com/peterh/liner"
)

var (
	// u: unlock, s: signXX, sendXX, n: newAccount, i: importXX
	passwordRegexp = regexp.MustCompile(`personal.[nusi]`)
	onlyWhitespace = regexp.MustCompile(`^\s*$`)
	exit           = regexp.MustCompile(`^\s*exit\s*;*\s*$`)
)

// HistoryFile is the file within the data directory to store input scrollback.
const HistoryFile = "history"

// DefaultPrompt is the default prompt line prefix to use for user input querying.
const DefaultPrompt = "> "

// Config is the collection of configurations to fine tune the behavior of the
// JavaScript console.
type Config struct {
	DataDir  string               // Data directory to store the console history at
	DocRoot  string               // Filesystem path from where to load JavaScript files from
	Client   *ethereum_rpc.Client // RPC client to execute Ethereum requests through
	Prompt   string               // Input prompt prefix string (defaults to DefaultPrompt)
	Prompter prompt.UserPrompter  // Input prompter to allow interactive user feedback (defaults to TerminalPrompter)
	Printer  io.Writer            // Output writer to serialize any display strings to (defaults to os.Stdout)
	Preload  []string             // Absolute paths to JavaScript files to preload
	NodeUrl  string               // Hmy Node url
	ShardId  int                  // Hmy Shard ID
	Net      string               // Hmy  Network
}

// Console is a JavaScript interpreted runtime environment. It is a fully fledged
// JavaScript console attached to a running node via an external or in-process RPC
// client.
type Console struct {
	client   *ethereum_rpc.Client // RPC client to execute Ethereum requests through
	jsre     *jsre.JSRE           // JavaScript runtime environment running the interpreter
	prompt   string               // Input prompt prefix string
	prompter prompt.UserPrompter  // Input prompter to allow interactive user feedback
	histPath string               // Absolute path to the console scrollback history
	history  []string             // Scroll history maintained by the console
	printer  io.Writer            // Output writer to serialize any display strings to
	nodeUrl  string               // Hmy Node url
	shardId  int                  // Hmy Shard ID
	net      string               // Hmy  Network
}

// New initializes a JavaScript interpreted runtime environment and sets defaults
// with the config struct.
func New(config Config) (*Console, error) {
	// Handle unset config values gracefully
	if config.Prompter == nil {
		config.Prompter = prompt.Stdin
	}
	if config.Prompt == "" {
		config.Prompt = DefaultPrompt
	}
	if config.Printer == nil {
		config.Printer = colorable.NewColorableStdout()
	}

	// Initialize the console and return
	console := &Console{
		client:   config.Client,
		jsre:     jsre.New(config.DocRoot, config.Printer),
		prompt:   config.Prompt,
		prompter: config.Prompter,
		printer:  config.Printer,
		histPath: filepath.Join(config.DataDir, HistoryFile),
		nodeUrl:  config.NodeUrl,
		shardId:  config.ShardId,
		net:      config.Net,
	}
	if err := os.MkdirAll(config.DataDir, 0700); err != nil {
		return nil, err
	}
	if err := console.init(config.Preload); err != nil {
		return nil, err
	}
	return console, nil
}

// init retrieves the available APIs from the remote RPC provider and initializes
// the console's JavaScript namespaces based on the exposed modules.
func (c *Console) init(preload []string) error {
	c.initConsoleObject()

	// Initialize the JavaScript <-> Go RPC bridge.
	bridge := newBridge(c)
	if err := c.initWeb3(bridge); err != nil {
		return err
	}
	if err := c.initExtensions(); err != nil {
		return err
	}

	// Add bridge overrides for web3.js functionality.
	c.jsre.Do(func(vm *goja.Runtime) {
		c.initPersonal(vm, bridge)
		c.initEth(vm, bridge)
	})

	// Preload JavaScript files.
	for _, path := range preload {
		if err := c.jsre.Exec(path); err != nil {
			failure := err.Error()
			if gojaErr, ok := err.(*goja.Exception); ok {
				failure = gojaErr.String()
			}
			return fmt.Errorf("%s: %v", path, failure)
		}
	}

	// Configure the input prompter for history and tab completion.
	if c.prompter != nil {
		if content, err := ioutil.ReadFile(c.histPath); err != nil {
			c.prompter.SetHistory(nil)
		} else {
			c.history = strings.Split(string(content), "\n")
			c.prompter.SetHistory(c.history)
		}
		c.prompter.SetWordCompleter(c.AutoCompleteInput)
	}
	return nil
}

func (c *Console) initConsoleObject() {
	c.jsre.Do(func(vm *goja.Runtime) {
		console := vm.NewObject()
		console.Set("log", c.consoleOutput)
		console.Set("error", c.consoleOutput)
		vm.Set("console", console)
	})
}

func (c *Console) initWeb3(bridge *bridge) error {
	var err error

	bnJS, err := deps.Asset("bignumber.js")
	if err != nil {
		return err
	}

	web3JS, err := deps.Asset("web3.js")
	if err != nil {
		return err
	}

	if err := c.jsre.Compile("bignumber.js", string(bnJS)); err != nil {
		return fmt.Errorf("bignumber.js: %v", err)
	}
	if err := c.jsre.Compile("web3.js", string(web3JS)); err != nil {
		return fmt.Errorf("web3.js: %v", err)
	}
	if _, err := c.jsre.Run("var Web3 = require('web3');"); err != nil {
		return fmt.Errorf("web3 require: %v", err)
	}

	c.jsre.Do(func(vm *goja.Runtime) {
		transport := vm.NewObject()
		transport.Set("send", jsre.MakeCallback(vm, bridge.Send))
		transport.Set("sendAsync", jsre.MakeCallback(vm, bridge.Send))
		vm.Set("_consoleWeb3Transport", transport)
		_, err = vm.RunString("var web3 = new Web3(_consoleWeb3Transport)")
	})
	return err
}

// initExtensions loads and registers web3.js extensions.
func (c *Console) initExtensions() error {
	// Compute aliases from server-provided modules.
	apis, err := c.client.SupportedModules()
	if err != nil {
		return fmt.Errorf("api modules: %v", err)
	}
	aliases := map[string]struct{}{"eth": {}, "personal": {}}
	for api := range apis {
		if api == "web3" {
			continue
		}
		aliases[api] = struct{}{}
		if file, ok := web3ext.Modules[api]; ok {
			if err = c.jsre.Compile(api+".js", file); err != nil {
				return fmt.Errorf("%s.js: %v", api, err)
			}
		}
	}

	// Apply aliases.
	c.jsre.Do(func(vm *goja.Runtime) {
		web3 := getObject(vm, "web3")
		for name := range aliases {
			if v := web3.Get(name); v != nil {
				vm.Set(name, v)
			}
		}
	})
	return nil
}

// initPersonal redirects account-related API methods through the bridge.
//
// If the console is in interactive mode and the 'personal' API is available, override
// the openWallet, unlockAccount, newAccount and sign methods since these require user
// interaction. The original web3 callbacks are stored in 'jeth'. These will be called
// by the bridge after the prompt and send the original web3 request to the backend.
func (c *Console) initPersonal(vm *goja.Runtime, bridge *bridge) {
	personal := getObject(vm, "personal")
	if personal == nil || c.prompter == nil {
		return
	}
	personal.Set("getListAccounts", jsre.MakeCallback(vm, bridge.HmyGetListAccounts))
	personal.Set("signTransaction", jsre.MakeCallback(vm, bridge.callbackProtected(bridge.HmySignTransaction)))
	personal.Set("sendTransaction", jsre.MakeCallback(vm, bridge.callbackProtected(bridge.HmySendTransaction)))
	personal.Set("lockAccount", jsre.MakeCallback(vm, bridge.callbackProtected(bridge.HmyLockAccount)))
	personal.Set("importRawKey", jsre.MakeCallback(vm, bridge.HmyImportRawKey))
	personal.Set("unlockAccount", jsre.MakeCallback(vm, bridge.callbackProtected(bridge.HmyUnlockAccount)))
	personal.Set("newAccount", jsre.MakeCallback(vm, bridge.callbackProtected(bridge.HmyNewAccount)))
	personal.Set("sign", jsre.MakeCallback(vm, bridge.callbackProtected(bridge.HmySign)))

	_, err := vm.RunString(`Object.defineProperty(personal, "listAccounts", {get: personal.getListAccounts});`)
	if err != nil {
		panic(err)
	}
}

func (c *Console) initEth(vm *goja.Runtime, bridge *bridge) {
	eth := getObject(vm, "eth")
	if eth == nil || c.prompter == nil {
		return
	}

	eth.Set("sendTransaction", jsre.MakeCallback(vm, bridge.callbackProtected(bridge.HmySendTransaction)))
	eth.Set("signTransaction", jsre.MakeCallback(vm, bridge.callbackProtected(bridge.HmySignTransaction)))
}

func (c *Console) clearHistory() {
	c.history = nil
	c.prompter.ClearHistory()
	if err := os.Remove(c.histPath); err != nil {
		fmt.Fprintln(c.printer, "can't delete history file:", err)
	} else {
		fmt.Fprintln(c.printer, "history file deleted")
	}
}

// consoleOutput is an override for the console.log and console.error methods to
// stream the output into the configured output stream instead of stdout.
func (c *Console) consoleOutput(call goja.FunctionCall) goja.Value {
	var output []string
	for _, argument := range call.Arguments {
		output = append(output, fmt.Sprintf("%v", argument))
	}
	fmt.Fprintln(c.printer, strings.Join(output, " "))
	return goja.Null()
}

// AutoCompleteInput is a pre-assembled word completer to be used by the user
// input prompter to provide hints to the user about the methods available.
func (c *Console) AutoCompleteInput(line string, pos int) (string, []string, string) {
	// No completions can be provided for empty inputs
	if len(line) == 0 || pos == 0 {
		return "", nil, ""
	}
	// Chunck data to relevant part for autocompletion
	// E.g. in case of nested lines eth.getBalance(eth.coinb<tab><tab>
	start := pos - 1
	for ; start > 0; start-- {
		// Skip all methods and namespaces (i.e. including the dot)
		if line[start] == '.' || (line[start] >= 'a' && line[start] <= 'z') || (line[start] >= 'A' && line[start] <= 'Z') {
			continue
		}
		// Handle web3 in a special way (i.e. other numbers aren't auto completed)
		if start >= 3 && line[start-3:start] == "web3" {
			start -= 3
			continue
		}
		// We've hit an unexpected character, autocomplete form here
		start++
		break
	}
	return line[:start], c.jsre.CompleteKeywords(line[start:pos]), line[pos:]
}

// Welcome show summary of current Geth instance and some metadata about the
// console's available modules.
func (c *Console) Welcome() {
	message := "Welcome to the Hmy JavaScript console!\n\n"

	// Print some generic Geth metadata
	if res, err := c.jsre.Run(`
		var message = "instance: " + web3.version.node + "\n";
		try {
			message += "coinbase: " + eth.coinbase + "\n";
		} catch (err) {}
		message += "at shard: " + hmy.shardID + "\n";
		message += "at block: " + eth.blockNumber + " (" + new Date(1000 * eth.getBlock(eth.blockNumber).timestamp) + ")\n";
		try {
			message += " datadir: " + admin.datadir + "\n";
		} catch (err) {}
		message
	`); err == nil {
		message += res.String()
	}
	// List all the supported modules for the user to call
	if apis, err := c.client.SupportedModules(); err == nil {
		modules := make([]string, 0, len(apis))
		for api, version := range apis {
			modules = append(modules, fmt.Sprintf("%s:%s", api, version))
		}
		sort.Strings(modules)
		message += " modules: " + strings.Join(modules, " ") + "\n"
	}
	message += "\nTo exit, press ctrl-d"
	fmt.Fprintln(c.printer, message)
}

// Evaluate executes code and pretty prints the result to the specified output
// stream.
func (c *Console) Evaluate(statement string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(c.printer, "[native] error: %v\n", r)
		}
	}()
	c.jsre.Evaluate(statement, c.printer)
}

// Interactive starts an interactive user session, where input is propted from
// the configured user prompter.
func (c *Console) Interactive() {
	var (
		prompt      = c.prompt             // the current prompt line (used for multi-line inputs)
		indents     = 0                    // the current number of input indents (used for multi-line inputs)
		input       = ""                   // the current user input
		inputLine   = make(chan string, 1) // receives user input
		inputErr    = make(chan error, 1)  // receives liner errors
		requestLine = make(chan string)    // requests a line of input
		interrupt   = make(chan os.Signal, 1)
	)

	// Monitor Ctrl-C. While liner does turn on the relevant terminal mode bits to avoid
	// the signal, a signal can still be received for unsupported terminals. Unfortunately
	// there is no way to cancel the line reader when this happens. The readLines
	// goroutine will be leaked in this case.
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	// The line reader runs in a separate goroutine.
	go c.readLines(inputLine, inputErr, requestLine)
	defer close(requestLine)

	for {
		// Send the next prompt, triggering an input read.
		requestLine <- prompt

		select {
		case <-interrupt:
			fmt.Fprintln(c.printer, "caught interrupt, exiting")
			return

		case err := <-inputErr:
			if err == liner.ErrPromptAborted {
				// When prompting for multi-line input, the first Ctrl-C resets
				// the multi-line state.
				prompt, indents, input = c.prompt, 0, ""
				continue
			}
			return

		case line := <-inputLine:
			// User input was returned by the prompter, handle special cases.
			if indents <= 0 && exit.MatchString(line) {
				return
			}
			if onlyWhitespace.MatchString(line) {
				continue
			}
			// Append the line to the input and check for multi-line interpretation.
			input += line + "\n"
			indents = countIndents(input)
			if indents <= 0 {
				prompt = c.prompt
			} else {
				prompt = strings.Repeat(".", indents*3) + " "
			}
			// If all the needed lines are present, save the command and run it.
			if indents <= 0 {
				if len(input) > 0 && input[0] != ' ' && !passwordRegexp.MatchString(input) {
					if command := strings.TrimSpace(input); len(c.history) == 0 || command != c.history[len(c.history)-1] {
						c.history = append(c.history, command)
						if c.prompter != nil {
							c.prompter.AppendHistory(command)
						}
					}
				}
				c.Evaluate(input)
				input = ""
			}
		}
	}
}

// readLines runs in its own goroutine, prompting for input.
func (c *Console) readLines(input chan<- string, errc chan<- error, prompt <-chan string) {
	for p := range prompt {
		line, err := c.prompter.PromptInput(p)
		if err != nil {
			errc <- err
		} else {
			input <- line
		}
	}
}

// countIndents returns the number of identations for the given input.
// In case of invalid input such as var a = } the result can be negative.
func countIndents(input string) int {
	var (
		indents     = 0
		inString    = false
		strOpenChar = ' '   // keep track of the string open char to allow var str = "I'm ....";
		charEscaped = false // keep track if the previous char was the '\' char, allow var str = "abc\"def";
	)

	for _, c := range input {
		switch c {
		case '\\':
			// indicate next char as escaped when in string and previous char isn't escaping this backslash
			if !charEscaped && inString {
				charEscaped = true
			}
		case '\'', '"':
			if inString && !charEscaped && strOpenChar == c { // end string
				inString = false
			} else if !inString && !charEscaped { // begin string
				inString = true
				strOpenChar = c
			}
			charEscaped = false
		case '{', '(':
			if !inString { // ignore brackets when in string, allow var str = "a{"; without indenting
				indents++
			}
			charEscaped = false
		case '}', ')':
			if !inString {
				indents--
			}
			charEscaped = false
		default:
			charEscaped = false
		}
	}

	return indents
}

// Execute runs the JavaScript file specified as the argument.
func (c *Console) Execute(path string) error {
	return c.jsre.Exec(path)
}

// Stop cleans up the console and terminates the runtime environment.
func (c *Console) Stop(graceful bool) error {
	if err := ioutil.WriteFile(c.histPath, []byte(strings.Join(c.history, "\n")), 0600); err != nil {
		return err
	}
	if err := os.Chmod(c.histPath, 0600); err != nil { // Force 0600, even if it was different previously
		return err
	}
	c.jsre.Stop(graceful)
	return nil
}

func (b *bridge) callbackProtected(protectedFunc func(call jsre.Call) (goja.Value, error)) func(call jsre.Call) (goja.Value, error) {
	return func(call jsre.Call) (goja.Value, error) {
		var availableCB goja.Callable = nil
		for i, args := range call.Arguments {
			if cb, ok := goja.AssertFunction(args); ok {
				availableCB = cb
				call.Arguments = call.Arguments[:i] // callback must be last
				break
			}
		}

		value, err := protectedFunc(call)
		jsErr := goja.Undefined()
		if err != nil {
			jsErr = call.VM.NewGoError(err)
		}
		if availableCB != nil {
			_, _ = availableCB(nil, jsErr, value)
		}

		return value, err
	}
}

func (b *bridge) HmyGetListAccounts(call jsre.Call) (goja.Value, error) {
	var accounts = []string{}

	for _, name := range store.LocalAccounts() {
		ks := store.FromAccountName(name)
		allAccounts := ks.Accounts()
		for _, account := range allAccounts {
			accounts = append(accounts, account.Address.String())
		}
	}

	return call.VM.ToValue(accounts), nil
}

func (b *bridge) HmySignTransaction(call jsre.Call) (goja.Value, error) {
	txObj := call.Arguments[0].ToObject(call.VM)
	password := call.Arguments[1].String()

	from := getStringFromJsObjWithDefault(txObj, "from", "")
	to := getStringFromJsObjWithDefault(txObj, "to", "")
	gasLimit := getStringFromJsObjWithDefault(txObj, "gas", "1000000")
	amount := getStringFromJsObjWithDefault(txObj, "value", "0")
	gasPrice := getStringFromJsObjWithDefault(txObj, "gasPrice", "1")
	input, err := transaction.StringToByte(getStringFromJsObjWithDefault(txObj, "data", ""))
	if err != nil {
		return nil, err
	}

	networkHandler := rpc.NewHTTPHandler(b.console.nodeUrl)
	chanId, err := common.StringToChainID(b.console.net)
	if err != nil {
		return nil, err
	}

	ks, acct, err := store.UnlockedKeystore(from, password)
	if err != nil {
		return nil, err
	}
	ctrlr := transaction.NewController(networkHandler, ks, acct, *chanId, func(controller *transaction.Controller) {
		// nop
	})

	tempLimit, err := strconv.ParseInt(gasLimit, 10, 64)
	if err != nil {
		return nil, err
	}
	if tempLimit < 0 {
		return nil, errors.New(fmt.Sprintf("gas-limit can not be negative: %s", gasLimit))
	}
	gLimit := uint64(tempLimit)

	amt, err := common.NewDecFromString(amount)
	if err != nil {
		return nil, fmt.Errorf("amount %w", err)
	}

	gPrice, err := common.NewDecFromString(gasPrice)
	if err != nil {
		return nil, fmt.Errorf("gas-price %w", err)
	}

	toP := &to
	if to == "" {
		toP = nil
	}

	nonce := transaction.GetNextPendingNonce(from, networkHandler)
	err = ctrlr.SignTransaction(
		nonce, gLimit,
		toP,
		uint32(b.console.shardId), uint32(b.console.shardId),
		amt, gPrice,
		input,
	)
	if err != nil {
		return nil, err
	}

	info := ctrlr.TransactionInfo()

	return call.VM.ToValue(map[string]interface{}{
		"raw": ctrlr.RawTransaction(),
		"tx": map[string]string{
			"nonce":    "0x" + big.NewInt(int64(info.Nonce())).Text(16),
			"gasPrice": "0x" + info.GasPrice().Text(16),
			"gas":      "0x" + big.NewInt(int64(info.GasLimit())).Text(16),
			"to":       info.To().Hex(),
			"value":    "0x" + info.Value().Text(16),
			"input":    "0x" + hex.EncodeToString(info.Data()),
			"v":        "0x" + info.V().Text(16),
			"r":        "0x" + info.R().Text(16),
			"s":        "0x" + info.S().Text(16),
			"hash":     info.Hash().Hex(),
		},
	}), nil
}

func (b *bridge) HmySendTransaction(call jsre.Call) (goja.Value, error) {
	txObj := call.Arguments[0].ToObject(call.VM)
	password := ""
	if len(call.Arguments) > 1 {
		password = call.Arguments[1].String()
	}

	from := getStringFromJsObjWithDefault(txObj, "from", "")
	to := getStringFromJsObjWithDefault(txObj, "to", "")
	gasLimit := getStringFromJsObjWithDefault(txObj, "gas", "1000000")
	amount := getStringFromJsObjWithDefault(txObj, "value", "0")
	gasPrice := getStringFromJsObjWithDefault(txObj, "gasPrice", "1")
	input, err := transaction.StringToByte(getStringFromJsObjWithDefault(txObj, "data", ""))
	if err != nil {
		return nil, err
	}

	networkHandler := rpc.NewHTTPHandler(b.console.nodeUrl)
	chanId, err := common.StringToChainID(b.console.net)
	if err != nil {
		return nil, err
	}

	ks, acct, err := store.UnlockedKeystore(from, password)
	if err != nil {
		return nil, err
	}
	ctrlr := transaction.NewController(networkHandler, ks, acct, *chanId, func(controller *transaction.Controller) {
		// nop
	})

	tempLimit, err := strconv.ParseInt(gasLimit, 10, 64)
	if err != nil {
		return nil, err
	}
	if tempLimit < 0 {
		return nil, errors.New(fmt.Sprintf("gas-limit can not be negative: %s", gasLimit))
	}
	gLimit := uint64(tempLimit)

	amt, err := common.NewDecFromString(amount)
	if err != nil {
		return nil, fmt.Errorf("amount %w", err)
	}

	gPrice, err := common.NewDecFromString(gasPrice)
	if err != nil {
		return nil, fmt.Errorf("gas-price %w", err)
	}

	toP := &to
	if to == "" {
		toP = nil
	}

	nonce := transaction.GetNextPendingNonce(from, networkHandler)
	err = ctrlr.ExecuteTransaction(
		nonce, gLimit,
		toP,
		uint32(b.console.shardId), uint32(b.console.shardId),
		amt, gPrice,
		input,
	)
	if err != nil {
		return nil, err
	}

	return call.VM.ToValue(*ctrlr.TransactionHash()), nil
}

func (b *bridge) HmyLockAccount(call jsre.Call) (goja.Value, error) {
	address := call.Arguments[0].String()

	_, _, err := store.LockKeystore(address)
	if err != nil {
		return nil, err
	}

	return goja.Null(), nil
}

func (b *bridge) HmyImportRawKey(call jsre.Call) (goja.Value, error) {
	privateKey := call.Arguments[0].String()
	password := call.Arguments[1].String()

	name, err := account.ImportFromPrivateKey(privateKey, "", password)
	if err != nil {
		return nil, err
	}

	return call.VM.ToValue(name), nil
}

func (b *bridge) HmyUnlockAccount(call jsre.Call) (goja.Value, error) {
	if len(call.Arguments) < 3 {
		return nil, errors.New("arguments < 3")
	}
	address := call.Arguments[0].String()
	password := call.Arguments[1].String()
	unlockDuration := call.Arguments[2].ToInteger()

	_, _, err := store.UnlockedKeystoreTimeLimit(address, password, time.Duration(unlockDuration)*time.Second)
	if err != nil {
		return nil, err
	}

	return goja.Null(), nil
}

func (b *bridge) HmyNewAccount(call jsre.Call) (goja.Value, error) {
	return goja.Null(), nil
}

func (b *bridge) HmySign(call jsre.Call) (goja.Value, error) {
	dataToSign := call.Arguments[0].String()
	addressStr := call.Arguments[1].String()
	password := call.Arguments[2].String()

	ks := store.FromAddress(addressStr)
	if ks == nil {
		return nil, fmt.Errorf("could not open local keystore for %s", addressStr)
	}

	acc, err := ks.Find(accounts.Account{Address: address.Parse(addressStr)})
	if err != nil {
		return nil, err
	}

	message, err := signMessageWithPassword(ks, acc, password, []byte(dataToSign))
	if err != nil {
		return nil, err
	}

	return call.VM.ToValue(hex.EncodeToString(message)), nil
}
