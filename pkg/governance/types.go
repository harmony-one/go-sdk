package governance

import (
	"encoding/json"
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
	if v.ProposalType == "approval" || v.ProposalType == "ranked-choice" {
		myType = append(myType, eip712.Type{
			Name: "choice",
			Type: "uint32[]",
		})
		var is []int
		if err := json.Unmarshal([]byte(v.Choice), &is); err == nil {
			choice = is
		} else {
			return nil, errors.Wrapf(err,
				"unexpected value of choice %s (expected uint32[])", choice,
			)
		}
	} else if v.ProposalType == "quadratic" || v.ProposalType == "weighted" {
		myType = append(myType, eip712.Type{
			Name: "choice",
			Type: "string",
		})
		choice = v.Choice
	} else if v.Privacy == "shutter" {
		myType = append(myType, eip712.Type{
			Name: "choice",
			Type: "string",
		})
		choice = v.Choice
	} else {
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
	}

	// order matters so this is added last
	myType = append(myType, eip712.Type{
		Name: "app",
		Type: "string",
	})

	if v.Timestamp == 0 {
		v.Timestamp = time.Now().Unix()
	}

	typedData := TypedData{
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
	}

	return &typedData, nil
}
