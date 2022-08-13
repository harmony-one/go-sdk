package governance

import (
	"fmt"
	"github.com/harmony-one/harmony/accounts"
	"github.com/harmony-one/harmony/accounts/keystore"
	"github.com/harmony-one/harmony/crypto/hash"
	"strconv"
	"strings"
	"time"
)

func timestampToDateString(timestamp float64) string {
	return time.Unix(int64(timestamp), 0).Format(time.RFC822)
}

func linePaddingPrint(content string, trim bool) {
	for _, line := range strings.Split(content, "\n") {
		trimLine := line
		if trim {
			trimLine = strings.TrimSpace(line)
		}
		if trimLine != "" {
			fmt.Printf("        %s\n", trimLine)
		}
	}
}

func signMessage(keyStore *keystore.KeyStore, account accounts.Account, data []byte) (sign []byte, err error) {
	signData := append([]byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(data))))
	msgHash := hash.Keccak256(append(signData, data...))

	sign, err = keyStore.SignHash(account, msgHash)
	if err != nil {
		return nil, err
	}

	if len(sign) != 65 {
		return nil, fmt.Errorf("sign error")
	}

	sign[64] += 0x1b
	return sign, nil
}
