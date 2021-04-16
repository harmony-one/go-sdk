package governance

import (
	"encoding/json"
	"fmt"
)

type Space struct {
	Name    string `json:"name"`
	Key     string `json:"key"`
	Network string `json:"network"`
	Symbol  string `json:"symbol"`
}

func listSpaces() (spaces map[string]*Space, err error) {
	var result map[string]*Space

	err = getAndParse(urlListSpace, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type ProposalMsgPayload struct {
	End      float64  `json:"end"`
	Body     string   `json:"body"`
	Name     string   `json:"name"`
	Start    float64  `json:"start"`
	Choices  []string `json:"choices"`
	Snapshot int      `json:"snapshot"`
}

type ProposalMsg struct {
	Version   string             `json:"version"`
	Timestamp string             `json:"timestamp"`
	Space     string             `json:"space"`
	Type      string             `json:"type"`
	Payload   ProposalMsgPayload `json:"payload"`
}

type Proposal struct {
	Address         string      `json:"address"`
	Msg             ProposalMsg `json:"msg"`
	Sig             string      `json:"sig"`
	AuthorIpfsHash  string      `json:"authorIpfsHash"`
	RelayerIpfsHash string      `json:"relayerIpfsHash"`
}

func listProposalsBySpace(spaceName string) (spaces map[string]*Proposal, err error) {
	var result map[string]*Proposal

	err = getAndParse(governanceApi(fmt.Sprintf(urlListProposalsBySpace, spaceName)), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type ProposalIPFSMsg struct {
	Version   string             `json:"version"`
	Timestamp string             `json:"timestamp"`
	Space     string             `json:"space"`
	Type      string             `json:"type"`
	Payload   ProposalMsgPayload `json:"payload"`
}

type ProposalVoteMsgPayload struct {
	Choice   int    `json:"choice"`
	Proposal string `json:"proposal"`
}

type ProposalVoteMsg struct {
	Version   string                 `json:"version"`
	Timestamp string                 `json:"timestamp"`
	Space     string                 `json:"space"`
	Type      string                 `json:"type"`
	Payload   ProposalVoteMsgPayload `json:"payload"`
}

type ProposalVote struct {
	Address         string          `json:"address"`
	Msg             ProposalVoteMsg `json:"msg"`
	Sig             string          `json:"sig"`
	AuthorIpfsHash  string          `json:"authorIpfsHash"`
	RelayerIpfsHash string          `json:"relayerIpfsHash"`
}

type ProposalIPFS struct {
	Address   string `json:"address"`
	Msg       string `json:"msg"`
	Sig       string `json:"sig"`
	Version   string `json:"version"`
	parsedMsg *ProposalIPFSMsg
	votes     map[string]*ProposalVote
}

func viewProposalsByProposalHash(proposalHash string) (proposal *ProposalIPFS, err error) {
	var result *ProposalIPFS = &ProposalIPFS{
		parsedMsg: &ProposalIPFSMsg{},
		votes:     make(map[string]*ProposalVote),
	}

	err = getAndParse(governanceApi(fmt.Sprintf(urlGetProposalInfo, proposalHash)), result)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(result.Msg), result.parsedMsg)
	if err != nil {
		return nil, err
	}

	err = getAndParse(governanceApi(fmt.Sprintf(urlListProposalsVoteBySpaceAndProposal, result.parsedMsg.Space, proposalHash)), &result.votes)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type NewProposalJson struct {
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
	Space     string `json:"space"`
	Type      string `json:"type"`
	Payload   struct {
		Name     string   `json:"name"`
		Body     string   `json:"body"`
		Choices  []string `json:"choices"`
		Start    float64  `json:"start"`
		End      float64  `json:"end"`
		Snapshot int      `json:"snapshot"`
		Metadata struct {
			Strategies []struct {
				Name   string `json:"name"`
				Params struct {
					Address  string `json:"address"`
					Symbol   string `json:"symbol"`
					Decimals int    `json:"decimals"`
				} `json:"params"`
			} `json:"strategies"`
		} `json:"metadata"`
	} `json:"payload"`
}

type NewProposalResponse struct {
	IpfsHash string `json:"ipfsHash"`
}

func submitMessage(address string, content string, sign string) (resp *NewProposalResponse, err error) {
	message, err := json.Marshal(map[string]string{
		"address": address,
		"msg":     content,
		"sig":     sign,
	})
	if err != nil {
		return nil, err
	}

	resp = &NewProposalResponse{}
	err = postAndParse(urlMessage, message, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type ValidatorsItem struct {
	Active           bool     `json:"active"`
	Apr              float64  `json:"apr,omitempty"`
	Address          string   `json:"address"`
	Name             string   `json:"name"`
	Rate             string   `json:"rate"`
	TotalStake       string   `json:"total_stake"`
	UptimePercentage *float64 `json:"uptime_percentage"`
	Identity         string   `json:"identity"`
	HasLogo          bool     `json:"hasLogo"`
}

type ValidatorsInfo struct {
	Validators  []ValidatorsItem `json:"validators"`
	TotalFound  int              `json:"totalFound"`
	Total       int              `json:"total"`
	TotalActive int              `json:"total_active"`
}

func getValidators(url governanceApi) map[string]ValidatorsItem {
	info := &ValidatorsInfo{}
	err := getAndParse(url, info)
	if err != nil {
		return nil
	}

	res := map[string]ValidatorsItem{}
	for _, validator := range info.Validators {
		res[validator.Address] = validator
	}

	return res
}
