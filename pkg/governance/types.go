package governance

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	eip712 "github.com/ethereum/go-ethereum/signer/core"
	"github.com/pkg/errors"
)

type Vote struct {
	From         string // --key
	Space        string // --space
	Proposal     string // --proposal
	ProposalType string // --proposal-type
	Choice       string // --choice
	Privacy      string // --privacy
	App          string // --app
	Timestamp    int64  // not exposed to the end user
}

func (v *Vote) ToEIP712() (*TypedData, error) {
	// common types regardless of parameters
	// key `app` appended later because order matters
	myType := []eip712.Type{
		{
			Name: "from",
			Type: "address",
		},
		{
			Name: "space",
			Type: "string",
		},
		{
			Name: "timestamp",
			Type: "uint64",
		},
	}

	var proposal interface{}
	isHex := strings.HasPrefix(v.Proposal, "0x")
	if isHex {
		myType = append(myType, eip712.Type{
			Name: "proposal",
			Type: "bytes32",
		})
		if proposalBytes, err := hexutil.Decode(v.Proposal); err != nil {
			return nil, errors.Wrapf(
				err, "invalid proposal hash %s", v.Proposal,
			)
		} else {
			// EncodePrimitiveValue accepts only hexutil.Bytes not []byte
			proposal = hexutil.Bytes(proposalBytes)
		}
	} else {
		myType = append(myType, eip712.Type{
			Name: "proposal",
			Type: "string",
		})
		proposal = v.Proposal
	}

	// vote type, vote choice and vote privacy
	// choice needs to be converted into its native format for envelope
	var choice interface{}
	// The space between [1, 2, 3] does not matter since we parse it
	// hmy governance vote-proposal \
	// 		--space harmony-mainnet.eth \
	// 		--proposal 0xTruncated \
	// 		--proposal-type {"approval","ranked-choice"} \
	// 		--choice "[1, 2, 3]" \
	//		--app my-app \
	// 		--key <name of pk>
	if v.ProposalType == "approval" || v.ProposalType == "ranked-choice" {
		myType = append(myType, eip712.Type{
			Name: "choice",
			Type: "uint32[]",
		})
		var is []int64
		if err := json.Unmarshal([]byte(v.Choice), &is); err == nil {
			local := make([]interface{}, len(is))
			for i := range is {
				local[i] = math.NewHexOrDecimal256(is[i])
			}
			choice = local
		} else {
			return nil, errors.Wrapf(err,
				"unexpected value of choice %s (expected uint32[])", choice,
			)
		}
	// The space between [1, 2, 3] does not matter to snapshot.org
	// But for comparing with the snapshot-js library, remove it
	// hmy governance vote-proposal \
	// 		--space harmony-mainnet.eth \
	// 		--proposal 0xTruncated \
	// 		--proposal-type {"quadratic","weighted"} \
	// 		--choice "[1,2,3]" \
	//		--app my-app \
	// 		--key <name of pk>
	} else if v.ProposalType == "quadratic" || v.ProposalType == "weighted" {
		myType = append(myType, eip712.Type{
			Name: "choice",
			Type: "string",
		})
		choice = v.Choice
	// hmy governance vote-proposal \
	// 		--space harmony-mainnet.eth \
	// 		--proposal 0xTruncated \
	// 		--proposal-type ANY \
	// 		--choice "unknown-format" \
	//		--app my-app \
	// 		--key <name of pk>
	//      --privacy shutter
	} else if v.Privacy == "shutter" {
		myType = append(myType, eip712.Type{
			Name: "choice",
			Type: "string",
		})
		choice = v.Choice
	// hmy governance vote-proposal \
	// 		--space harmony-mainnet.eth \
	// 		--proposal 0xTruncated \
	// 		--proposal-type single-choice \
	// 		--choice 1 \
	//		--app my-app \
	// 		--key <name of pk>
	} else if v.ProposalType == "single-choice" {
		myType = append(myType, eip712.Type{
			Name: "choice",
			Type: "uint32",
		})
		if x, err := strconv.Atoi(v.Choice); err != nil {
			return nil, errors.Wrapf(err,
				"unexpected value of choice %s (expected uint32)", choice,
			)
		} else {
			choice = math.NewHexOrDecimal256(int64(x))
		}
	} else {
		return nil, errors.New(
			fmt.Sprintf(
				"unknown proposal type %s",
				v.ProposalType,
			),
		)
	}

	// order matters so this is added last
	myType = append(myType, eip712.Type{
		Name: "app",
		Type: "string",
	})

	if v.Timestamp == 0 {
		v.Timestamp = time.Now().Unix()
	}

	return &TypedData{
		eip712.TypedData{
			Domain: eip712.TypedDataDomain{
				Name:    name,
				Version: version,
			},
			Types: eip712.Types{
				"EIP712Domain": {
					{
						Name: "name",
						Type: "string",
					},
					{
						Name: "version",
						Type: "string",
					},
				},
				"Vote": myType,
			},
			Message: eip712.TypedDataMessage{
				"from":  v.From,
				"space": v.Space,
				// EncodePrimitiveValue accepts string, float64, or this type
				"timestamp": math.NewHexOrDecimal256(v.Timestamp),
				"proposal":  proposal,
				"choice":    choice,
				"app":       v.App,
			},
			PrimaryType: "Vote",
		},
	}, nil
}
