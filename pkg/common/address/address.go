package address

// TODO Stopgap measure for now to get functionality up and running
// because this code was locked away under harmony/internal/..., later
// Make main harmony repo be a consumer of this, as this is indeed the sdk package, right.
import (
	"github.com/btcsuite/btcutil/bech32"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	// "github.com/tendermint/tendermint/libs/bech32"
)

const (
	Bech32AddressHRP = "one"
)

func decodeAndConvert(bech string) (string, []byte, error) {
	hrp, data, err := bech32.Decode(bech)
	if err != nil {
		return "", nil, errors.Wrap(err, "decoding bech32 failed")
	}
	converted, err := bech32.ConvertBits(data, 5, 8, false)
	if err != nil {
		return "", nil, errors.Wrap(err, "decoding bech32 failed")
	}
	return hrp, converted, nil
}

// TODO ek – the following functions use Ethereum addresses until we have a
//  proper abstraction set in place.

// ParseBech32Addr decodes the given bech32 address and populates the given
// human-readable-part string and address with the decoded result.
func ParseBech32Addr(b32 string, hrp *string, addr *ethCommon.Address) error {
	h, b, err := decodeAndConvert(b32)
	if err != nil {
		return errors.Wrapf(err, "cannot decode %#v as bech32 address", b32)
	}
	if len(b) != ethCommon.AddressLength {
		return errors.Errorf("decoded bech32 %#v has invalid length %d",
			b32, len(b))
	}
	*hrp = h
	addr.SetBytes(b)
	return nil
}

func Bech32ToAddress(b32 string) (addr ethCommon.Address, err error) {
	var hrp string
	err = ParseBech32Addr(b32, &hrp, &addr)
	if err == nil && hrp != Bech32AddressHRP {
		err = errors.Errorf("%#v is not a %#v address", b32, Bech32AddressHRP)
	}
	return
}

// ParseAddr parses the given address, either as bech32 or as hex.
// The result can be 0x00..00 if the passing param is not a correct address.
func Parse(s string) ethCommon.Address {
	if addr, err := Bech32ToAddress(s); err == nil {
		return addr
	}
	// The result can be 0x00...00 if the passing param is not a correct address.
	return ethCommon.HexToAddress(s)
}

func ToBech32(addr ethCommon.Address) string {
	b32, err := BuildBech32Addr(Bech32AddressHRP, addr)
	if err != nil {
		panic(err)
	}
	return b32
}

func BuildBech32Addr(hrp string, addr ethCommon.Address) (string, error) {
	return ConvertAndEncode(hrp, addr.Bytes())
}

func ConvertAndEncode(hrp string, data []byte) (string, error) {
	converted, err := bech32.ConvertBits(data, 8, 5, true)
	if err != nil {
		return "", errors.Wrap(err, "encoding bech32 failed")
	}
	return bech32.Encode(hrp, converted)

}
