package governance

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/harmony-one/harmony/accounts"
	"github.com/harmony-one/harmony/accounts/keystore"
)

type VoteMessage struct {
	Version   string  `json:"version"`
	Timestamp string  `json:"timestamp"`
	Space     string  `json:"space"`
	Type      string  `json:"type"`
	Payload   Payload `json:"payload"`
}

type Payload struct {
	Proposal string `json:"string"`
	Choice   int    `json:"choice"`
}

func Vote(keyStore *keystore.KeyStore, account accounts.Account, space string, proposalHash string, choice int) error {
	if choice < 0 {
		return fmt.Errorf("invalid choice, please choose positive value")
	}

	voteJson := &VoteMessage{
		Version:   version,
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
		Space:     space,
		Type:      voteType,
		Payload: Payload{
			Proposal: proposalHash,
			Choice:   choice,
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

	result, err := submitMessage(account.Address.String(), string(voteJsonData), fmt.Sprintf("0x%x", sign))
	if err != nil {
		return err
	}

	fmt.Println(indent(result))
	return nil
}
