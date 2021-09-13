package governance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/harmony/accounts"
	"github.com/harmony-one/harmony/accounts/keystore"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func PrintListSpace() error {
	spaces, err := listSpaces()
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetHeader([]string{"Key", "Name"})

	for key, space := range spaces {
		table.Append([]string{key, space.Name})
	}

	table.Render()
	return nil
}

func PrintListProposals(spaceName string) error {
	proposals, err := listProposalsBySpace(spaceName)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetHeader([]string{"Key", "Name", "Start Date", "End Date"})

	for key, proposal := range proposals {
		table.Append([]string{
			key,
			proposal.Msg.Payload.Name,
			timestampToDateString(proposal.Msg.Payload.Start),
			timestampToDateString(proposal.Msg.Payload.End),
		})
	}

	table.Render()
	return nil
}

func PrintViewProposal(proposalHash string) error {
	proposals, err := viewProposalsByProposalHash(proposalHash)
	if err != nil {
		return err
	}

	var validators map[string]ValidatorsItem
	switch proposals.parsedMsg.Space {
	case "staking-mainnet":
		validators = getValidators(urlGetValidatorsInMainNet)
	case "staking-testnet":
		validators = getValidators(urlGetValidatorsInTestNet)
	}

	fmt.Printf("Author : %s\n", address.ToBech32(address.Parse(proposals.Address)))
	fmt.Printf("IPFS   : %s\n", proposalHash)
	fmt.Printf("Space  : %s\n", proposals.parsedMsg.Space)
	fmt.Printf("Start  : %s\n", timestampToDateString(proposals.parsedMsg.Payload.Start))
	fmt.Printf("End    : %s\n", timestampToDateString(proposals.parsedMsg.Payload.End))
	fmt.Printf("Choose : %s\n", strings.Join(proposals.parsedMsg.Payload.Choices, " / "))
	fmt.Printf("Content: \n")

	linePaddingPrint(proposals.parsedMsg.Payload.Body, true)

	fmt.Printf("\n")
	fmt.Printf("Votes: \n")

	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)
	table.SetBorder(false)
	table.SetHeader([]string{"Address", "Choose", "Stack"})

	for _, vote := range proposals.votes {
		stack := "0"
		addr := address.ToBech32(address.Parse(vote.Address))
		if v, ok := validators[addr]; ok {
			float, err := strconv.ParseFloat(v.TotalStake, 64)
			if err == nil {
				stack = fmt.Sprintf("%.2f", float/1e18)
			}
		}

		choices := make([]string, 0)
		for _, choice := range vote.Msg.Payload.choices() {
			choices = append(choices, proposals.parsedMsg.Payload.Choices[choice-1])
		}

		table.Append([]string{
			addr,
			strings.Join(choices, ", "),
			stack,
		})
	}

	table.Render()
	linePaddingPrint(buf.String(), false)
	return nil
}

type NewProposalYaml struct {
	Space   string    `yaml:"space"`
	Start   time.Time `yaml:"start"`
	End     time.Time `yaml:"end"`
	Choices []string  `yaml:"choices"`
	Title   string    `yaml:"title"`
	Body    string    `yaml:"body"`
}

var proposalTemplate = []byte(`{
  "version": "0.2.0",
  "type": "proposal",
  "payload": {
    "metadata": {
      "strategies": [
        {
          "name": "erc20-balance-of",
          "params": {
            "address": "0x00eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
            "symbol": "ONE",
            "decimals": 18
          }
        }
      ]
    },
    "maxCanSelect": 1
  }
}`)

func NewProposal(keyStore *keystore.KeyStore, account accounts.Account, proposalYamlPath string) error {
	proposalYamlFile, err := os.Open(proposalYamlPath)
	if err != nil {
		return err
	}
	defer proposalYamlFile.Close()

	proposalYaml := &NewProposalYaml{}
	err = yaml.NewDecoder(proposalYamlFile).Decode(proposalYaml)
	if err != nil {
		return err
	}

	rand.Seed(time.Now().Unix())
	proposalJson := &NewProposalJson{}
	err = json.Unmarshal(proposalTemplate, proposalJson)
	if err != nil {
		return err
	}

	proposalJson.Space = proposalYaml.Space
	proposalJson.Timestamp = fmt.Sprintf("%d", time.Now().Unix())
	proposalJson.Payload.Name = proposalYaml.Title
	proposalJson.Payload.Body = proposalYaml.Body
	proposalJson.Payload.Choices = proposalYaml.Choices
	proposalJson.Payload.Start = float64(proposalYaml.Start.Unix())
	proposalJson.Payload.End = float64(proposalYaml.End.Unix())
	proposalJson.Payload.Snapshot = rand.Intn(9999999) + 1

	if !checkPermission(proposalJson.Space, account) {
		return fmt.Errorf("no permission!")
	}

	proposalJsonData, err := json.Marshal(proposalJson)
	if err != nil {
		return err
	}

	sign, err := signMessage(keyStore, account, proposalJsonData)
	if err != nil {
		return err
	}

	proposal, err := submitMessage(account.Address.String(), string(proposalJsonData), fmt.Sprintf("0x%x", sign))
	if err != nil {
		return err
	}

	fmt.Printf("IPFS   : %s\n", proposal.IpfsHash)

	return nil
}

func checkPermission(space string, account accounts.Account) bool {
	var validators map[string]ValidatorsItem
	switch space {
	case "staking-mainnet":
		validators = getValidators(urlGetValidatorsInMainNet)
	case "staking-testnet":
		validators = getValidators(urlGetValidatorsInTestNet)
	default:
		return true
	}

	if _, ok := validators[address.ToBech32(account.Address)]; ok {
		return true
	} else {
		return false
	}
}

type VoteMessage struct {
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
	Space     string `json:"space"`
	Type      string `json:"type"`
	Payload   struct {
		Proposal string `json:"proposal"`
		Choice   int    `json:"choice"`
		Metadata struct {
		} `json:"metadata"`
	} `json:"payload"`
}

func Vote(keyStore *keystore.KeyStore, account accounts.Account, proposalHash string, choiceText string) error {
	proposals, err := viewProposalsByProposalHash(proposalHash)
	if err != nil {
		return err
	}

	if !checkPermission(proposals.parsedMsg.Space, account) {
		return fmt.Errorf("no permission!")
	}

	chooseIndex := -1
	for i, choice := range proposals.parsedMsg.Payload.Choices {
		if choice == choiceText {
			chooseIndex = i + 1
			break
		}
	}

	if chooseIndex < 0 {
		return fmt.Errorf("error choose, please choose: %s", strings.Join(proposals.parsedMsg.Payload.Choices, " / "))
	}

	voteJson := &VoteMessage{
		Version:   "0.2.0",
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
		Space:     proposals.parsedMsg.Space,
		Type:      "vote",
		Payload: struct {
			Proposal string   `json:"proposal"`
			Choice   int      `json:"choice"`
			Metadata struct{} `json:"metadata"`
		}{
			Proposal: proposalHash,
			Choice:   chooseIndex,
			Metadata: struct{}{},
		},
	}

	voteJsonData, err := json.Marshal(voteJson)
	if err != nil {
		return err
	}

	sign, err := signMessage(keyStore, account, voteJsonData)
	if err != nil {
		return err
	}

	proposal, err := submitMessage(account.Address.String(), string(voteJsonData), fmt.Sprintf("0x%x", sign))
	if err != nil {
		return err
	}

	fmt.Printf("Vote IPFS: %s\n", proposal.IpfsHash)
	return nil
}
