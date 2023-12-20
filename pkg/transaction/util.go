package transaction

import (
	"encoding/hex"
	"fmt"
	"strings"
)

func StringToByte(dataStr string) ([]byte, error) {
	if len(dataStr) == 0 {
		return []byte{}, nil
	}
	if !strings.HasPrefix(dataStr, "0x") {
		return nil, fmt.Errorf("invalid data literal: %q", dataStr)
	}
	return hex.DecodeString(dataStr[2:])
}
