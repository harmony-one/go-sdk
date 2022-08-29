package governance

import (
	"encoding/json"

	"github.com/pkg/errors"
)

func submitMessage(address string, typedData *TypedData, sign string) (map[string]interface{}, error) {
	data, err := typedData.String()
	if err != nil {
		return nil, errors.Wrapf(err, "could not encode EIP712 data")
	}

	type body struct {
		Address  string          `json:"address"`
		Sig      string          `json:"sig"`
		JsonData json.RawMessage `json:"data"`
	}

	message, err := json.Marshal(body{
		Address:  address,
		Sig:      sign,
		JsonData: json.RawMessage(data),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "could not encode body")
	}
	return postAndParse(urlMessage, message)
}
