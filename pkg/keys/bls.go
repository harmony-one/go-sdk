package keys

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/harmony/crypto/bls"
)

func GenBlsKeys(passphrase, filePath string) error {
	privateKey := bls.RandPrivateKey()
	publicKey := privateKey.GetPublicKey()
	publicKeyHex := publicKey.SerializeToHexStr()

	if filePath == "" {
		cwd, _ := os.Getwd()
		filePath = fmt.Sprintf("%s/%s.key", cwd, publicKeyHex)
	}
	if !path.IsAbs(filePath) {
		return common.ErrNotAbsPath
	}
	privateKeyHex := privateKey.SerializeToHexStr()
	encryptedPrivateKeyStr, err := encrypt([]byte(privateKeyHex), passphrase)
	if err != nil {
		return err
	}
	err = writeToFile(filePath, encryptedPrivateKeyStr)
	if err != nil {
		return err
	}
	out := fmt.Sprintf(`
{"public-key" : "0x%s", "private-key" : "0x%s", "encrypted-private-key-path" : "%s"}`,
		publicKeyHex, privateKeyHex, filePath)
	fmt.Println(common.JSONPrettyFormat(out))
	return nil
}

func writeToFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) (string, error) {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return hex.EncodeToString(ciphertext), nil
}

func decrypt(encrypted []byte, passphrase string) (decrypted []byte, err error) {
	unhexed := make([]byte, hex.DecodedLen(len(encrypted)))
	if _, err = hex.Decode(unhexed, encrypted); err == nil {
		if decrypted, err = decryptRaw(unhexed, passphrase); err == nil {
			return decrypted, nil
		}
	}
	// At this point err != nil, either from hex decode or from decryptRaw.
	decrypted, binErr := decryptRaw(encrypted, passphrase)
	if binErr != nil {
		// Disregard binary decryption error and return the original error,
		// because our canonical form is hex and not binary.
		return nil, err
	}
	return decrypted, nil
}

func decryptRaw(data []byte, passphrase string) ([]byte, error) {
	var err error
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	return plaintext, err
}
