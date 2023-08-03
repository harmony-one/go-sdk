package governance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core"
	"github.com/pkg/errors"
)

// This embedded type was created to override the EncodeData function
// and remove the validation for a mandatory chain id
type TypedData struct {
	core.TypedData
}

// dataMismatchError generates an error for a mismatch between
// the provided type and data
func dataMismatchError(encType string, encValue interface{}) error {
	return fmt.Errorf("provided data '%v' doesn't match type '%s'", encValue, encType)
}

// EncodeData generates the following encoding:
// `enc(value₁) ‖ enc(value₂) ‖ … ‖ enc(valueₙ)`
//
// each encoded member is 32-byte long
// This method overridden here to remove the validation for mandatory chain id
func (typedData *TypedData) EncodeData(primaryType string, data map[string]interface{}, depth int) (hexutil.Bytes, error) {
	// if err := typedData.validate(); err != nil {
	// 	return nil, err
	// }

	buffer := bytes.Buffer{}

	// Verify extra data
	if len(typedData.Types[primaryType]) < len(data) {
		return nil, errors.New("there is extra data provided in the message")
	}

	// Add typehash
	buffer.Write(typedData.TypeHash(primaryType))

	// Add field contents. Structs and arrays have special handlers.
	for _, field := range typedData.Types[primaryType] {
		encType := field.Type
		encValue := data[field.Name]
		if encType[len(encType)-1:] == "]" {
			arrayValue, ok := encValue.([]interface{})
			if !ok {
				return nil, dataMismatchError(encType, encValue)
			}

			arrayBuffer := bytes.Buffer{}
			parsedType := strings.Split(encType, "[")[0]
			for _, item := range arrayValue {
				if typedData.Types[parsedType] != nil {
					mapValue, ok := item.(map[string]interface{})
					if !ok {
						return nil, dataMismatchError(parsedType, item)
					}
					encodedData, err := typedData.EncodeData(parsedType, mapValue, depth+1)
					if err != nil {
						return nil, err
					}
					arrayBuffer.Write(encodedData)
				} else {
					bytesValue, err := typedData.EncodePrimitiveValue(parsedType, item, depth)
					if err != nil {
						return nil, err
					}
					arrayBuffer.Write(bytesValue)
				}
			}

			buffer.Write(crypto.Keccak256(arrayBuffer.Bytes()))
		} else if typedData.Types[field.Type] != nil {
			mapValue, ok := encValue.(map[string]interface{})
			if !ok {
				return nil, dataMismatchError(encType, encValue)
			}
			encodedData, err := typedData.EncodeData(field.Type, mapValue, depth+1)
			if err != nil {
				return nil, err
			}
			buffer.Write(crypto.Keccak256(encodedData))
		} else {
			byteValue, err := typedData.EncodePrimitiveValue(encType, encValue, depth)
			if err != nil {
				return nil, err
			}
			buffer.Write(byteValue)
		}
	}
	return buffer.Bytes(), nil
}

type TypedDataMessage = map[string]interface{}

// HashStruct generates a keccak256 hash of the encoding of the provided data
// This method overridden here to allow calling the overriden EncodeData above
func (typedData *TypedData) HashStruct(primaryType string, data TypedDataMessage) (hexutil.Bytes, error) {
	encodedData, err := typedData.EncodeData(primaryType, data, 1)
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(encodedData), nil
}

func (typedData *TypedData) String() (string, error) {
	type domain struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	// this data structure created to remove unused fields
	// for example, domain type is not sent in post request
	// and neither are the blank fields in the domain
	type data struct {
		Domain  domain                `json:"domain"`
		Types   core.Types            `json:"types"`
		Message core.TypedDataMessage `json:"message"`
	}
	var ts uint64
	var err error
	if ts, err = toUint64(typedData.Message["timestamp"]); err != nil {
		// should not happen
		return "", errors.Wrapf(err, "timestamp")
	}
	formatted := data{
		Domain: domain{
			Name:    typedData.Domain.Name,
			Version: typedData.Domain.Version,
		},
		Types: core.Types{
			typedData.PrimaryType: typedData.Types[typedData.PrimaryType],
		},
		Message: core.TypedDataMessage{
			"space":    typedData.Message["space"],
			"proposal": typedData.Message["proposal"],
			"choice":   typedData.Message["choice"],
			"app":      typedData.Message["app"],
			"reason":   typedData.Message["reason"],
			// this conversion is required to stop snapshot
			// from complaining about `wrong envelope format`
			"timestamp": ts,
			"from":      typedData.Message["from"],
		},
	}
	// same comment as above
	if typedData.Types["Vote"][4].Type == "uint32" {
		if choice, err := toUint64(typedData.Message["choice"]); err != nil {
			return "", errors.Wrapf(err, "choice")
		} else {
			formatted.Message["choice"] = choice
		}
	// prevent hex choice interpretation
	} else if typedData.Types["Vote"][4].Type == "uint32[]" {
		arr := typedData.Message["choice"].([]interface{})
		res := make([]uint64, len(arr))
		for i, a := range arr {
			if c, err := toUint64(a); err != nil {
				return "", errors.Wrapf(err, "choice member %d", i)
			} else {
				res[i] = c
			}
		}
		formatted.Message["choice"] = res
	}
	message, err := json.Marshal(formatted)
	if err != nil {
		return "", err
	} else {
		return string(message), nil
	}
}

func toUint64(x interface{}) (uint64, error) {
	y, ok := x.(*math.HexOrDecimal256)
	if !ok {
		return 0, errors.New(
			fmt.Sprintf("%+v is not a *math.HexOrDecimal256", x),
		)
	}
	z := (*big.Int)(y)
	return z.Uint64(), nil
}