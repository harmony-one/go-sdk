package keys

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"strings"

	ffiBls "github.com/harmony-one/bls/ffi/go/bls"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/sharding"
	"github.com/harmony-one/go-sdk/pkg/validation"
	"github.com/harmony-one/harmony/crypto/bls"
	"github.com/harmony-one/harmony/crypto/hash"
	"github.com/harmony-one/harmony/shard"
	"github.com/harmony-one/harmony/staking/types"
	"golang.org/x/crypto/ssh/terminal"
)

//BlsKey - struct to represent bls key data
type BlsKey struct {
	PrivateKey     *ffiBls.SecretKey
	PublicKey      *ffiBls.PublicKey
	PublicKeyHex   string
	PrivateKeyHex  string
	Passphrase     string
	FilePath       string
	ShardPublicKey *shard.BlsPublicKey
}

//Initialize - initialize a bls key and assign a random private bls key if not already done
func (blsKey *BlsKey) Initialize() {
	if blsKey.PrivateKey == nil {
		blsKey.PrivateKey = bls.RandPrivateKey()
		blsKey.PrivateKeyHex = blsKey.PrivateKey.SerializeToHexStr()
		blsKey.PublicKey = blsKey.PrivateKey.GetPublicKey()
		blsKey.PublicKeyHex = blsKey.PublicKey.SerializeToHexStr()
	}
}

//Reset - resets the currently assigned private and public key fields
func (blsKey *BlsKey) Reset() {
	blsKey.PrivateKey = nil
	blsKey.PrivateKeyHex = ""
	blsKey.PublicKey = nil
	blsKey.PublicKeyHex = ""
}

// GenBlsKeys - generate a random bls key using the supplied passphrase, write it to disk at the given filePath
func GenBlsKeys(blsKey *BlsKey) error {
	blsKey.Initialize()
	out, err := writeBlsKeyToFile(blsKey)
	if err != nil {
		return err
	}
	fmt.Println(common.JSONPrettyFormat(out))
	return nil
}

// GenMultiBlsKeys - generate multiple BLS keys for a given shard and node/network
func GenMultiBlsKeys(blsKeys []*BlsKey, node string, count uint32, shardID uint32) error {
	blsKeys, _, err := generateMultipleBlsKeys(blsKeys, node, count, shardID)
	if err != nil {
		return err
	}

	outputs := []string{}
	for _, blsKey := range blsKeys {
		out, err := writeBlsKeyToFile(blsKey)
		if err != nil {
			return err
		}

		outputs = append(outputs, out)
	}

	if len(outputs) > 0 {
		fmt.Println(common.JSONPrettyFormat(fmt.Sprintf("[%s]", strings.Join(outputs[:], ","))))
	}

	return nil
}

func RecoverBlsKeyFromFile(passphrase, filePath string) error {
	if !path.IsAbs(filePath) {
		return common.ErrNotAbsPath
	}
	encryptedPrivateKeyBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	decryptedPrivateKeyBytes, err := decrypt(encryptedPrivateKeyBytes, passphrase)
	if err != nil {
		return err
	}
	privateKey, err := getBlsKey(string(decryptedPrivateKeyBytes))
	if err != nil {
		return err
	}
	publicKey := privateKey.GetPublicKey()
	publicKeyHex := publicKey.SerializeToHexStr()
	privateKeyHex := privateKey.SerializeToHexStr()
	out := fmt.Sprintf(`
{"public-key" : "0x%s", "private-key" : "0x%s"}`,
		publicKeyHex, privateKeyHex)
	fmt.Println(common.JSONPrettyFormat(out))
	return nil

}

func SaveBlsKey(passphrase, filePath, privateKeyHex string) error {
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
	privateKey, err := getBlsKey(privateKeyHex)
	if err != nil {
		return err
	}
	if filePath == "" {
		cwd, _ := os.Getwd()
		filePath = fmt.Sprintf("%s/%s.key", cwd, privateKey.GetPublicKey().SerializeToHexStr())
	}
	if !path.IsAbs(filePath) {
		return common.ErrNotAbsPath
	}
	encryptedPrivateKeyStr, err := encrypt([]byte(privateKeyHex), passphrase)
	if err != nil {
		return err
	}
	err = writeToFile(filePath, encryptedPrivateKeyStr)
	if err != nil {
		return err
	}
	fmt.Printf("Encrypted and saved bls key to: %s\n", filePath)
	return nil
}

func GetPublicBlsKey(privateKeyHex string) error {
	privateKey, err := getBlsKey(privateKeyHex)
	if err != nil {
		return err
	}
	publicKeyHex := privateKey.GetPublicKey().SerializeToHexStr()
	out := fmt.Sprintf(`
{"public-key" : "0x%s", "private-key" : "0x%s"}`,
		publicKeyHex, privateKeyHex)
	fmt.Println(common.JSONPrettyFormat(out))
	return nil

}

func VerifyBLSKeys(blsPubKeys []string) ([]shard.BLSSignature, error) {
	blsSigs := make([]shard.BLSSignature, len(blsPubKeys))
	for i := 0; i < len(blsPubKeys); i++ {
		sig, err := VerifyBLS(strings.TrimPrefix(blsPubKeys[i], "0x"))
		if err != nil {
			return nil, err
		}
		blsSigs[i] = sig
	}
	return blsSigs, nil
}

func VerifyBLS(blsPubKey string) (shard.BLSSignature, error) {
	var sig shard.BLSSignature
	// look for key file in the current directory
	// if not ask for the absolute path
	cwd, _ := os.Getwd()
	filePath := fmt.Sprintf("%s/%s.key", cwd, blsPubKey)
	encryptedPrivateKeyBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("For bls public key: %s\n", blsPubKey)
		fmt.Println("Enter the absolute path to the encrypted bls private key file:")
		filePath, _ := reader.ReadString('\n')
		if !path.IsAbs(filePath) {
			return sig, common.ErrNotAbsPath
		}
		filePath = strings.TrimSpace(filePath)
		encryptedPrivateKeyBytes, err = ioutil.ReadFile(filePath)
		if err != nil {
			return sig, err
		}
	}
	// ask passphrase for bls key twice
	fmt.Println("Enter the bls passphrase:")
	pass, _ := terminal.ReadPassword(int(os.Stdin.Fd()))

	decryptedPrivateKeyBytes, err := decrypt(encryptedPrivateKeyBytes, string(pass))
	if err != nil {
		return sig, err
	}
	privateKey, err := getBlsKey(string(decryptedPrivateKeyBytes))
	if err != nil {
		return sig, err
	}
	publicKey := privateKey.GetPublicKey()
	publicKeyHex := publicKey.SerializeToHexStr()

	if publicKeyHex != blsPubKey {
		return sig, errors.New("bls key could not be verified")
	}

	messageBytes := []byte(types.BlsVerificationStr)
	msgHash := hash.Keccak256(messageBytes)
	signature := privateKey.SignHash(msgHash[:])

	bytes := signature.Serialize()
	if len(bytes) != shard.BLSSignatureSizeInBytes {
		return sig, errors.New("bls key length is not 96 bytes")
	}
	copy(sig[:], bytes)
	return sig, nil
}

func getBlsKey(privateKeyHex string) (*ffiBls.SecretKey, error) {
	privateKey := &ffiBls.SecretKey{}
	err := privateKey.DeserializeHexStr(string(privateKeyHex))
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func writeBlsKeyToFile(blsKey *BlsKey) (string, error) {
	if blsKey.FilePath == "" {
		cwd, _ := os.Getwd()
		blsKey.FilePath = fmt.Sprintf("%s/%s.key", cwd, blsKey.PublicKeyHex)
	}
	if !path.IsAbs(blsKey.FilePath) {
		return "", common.ErrNotAbsPath
	}
	encryptedPrivateKeyStr, err := encrypt([]byte(blsKey.PrivateKeyHex), blsKey.Passphrase)
	if err != nil {
		return "", err
	}
	err = writeToFile(blsKey.FilePath, encryptedPrivateKeyStr)
	if err != nil {
		return "", err
	}
	out := fmt.Sprintf(`
{"public-key" : "%s", "private-key" : "%s", "encrypted-private-key-path" : "%s"}`,
		blsKey.PublicKeyHex, blsKey.PrivateKeyHex, blsKey.FilePath)

	return out, nil
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

func generateMultipleBlsKeys(blsKeys []*BlsKey, node string, count uint32, shardID uint32) ([]*BlsKey, int, error) {
	shardingStructure, err := sharding.Structure(node)
	if err != nil {
		return blsKeys, -1, err
	}
	shardCount := len(shardingStructure)

	if !validation.ValidShardID(shardID, uint32(shardCount)) {
		return blsKeys, shardCount, fmt.Errorf("node %s only supports a total of %d shards - supplied shard id %d isn't valid", node, shardCount, shardID)
	}

	for _, blsKey := range blsKeys {
		for {
			blsKey.Initialize()
			shardPubKey := new(shard.BlsPublicKey)
			if err = shardPubKey.FromLibBLSPublicKey(blsKey.PublicKey); err != nil {
				return blsKeys, shardCount, err
			}

			if blsKeyMatchesShardID(shardPubKey, shardID, shardCount) {
				blsKey.ShardPublicKey = shardPubKey
				break
			} else {
				blsKey.Reset()
			}
		}
	}

	return blsKeys, shardCount, nil
}

func blsKeyMatchesShardID(pubKey *shard.BlsPublicKey, shardID uint32, shardCount int) bool {
	bigShardCount := big.NewInt(int64(shardCount))
	resolvedShardID := int(new(big.Int).Mod(pubKey.Big(), bigShardCount).Int64())
	return int(shardID) == resolvedShardID
}
