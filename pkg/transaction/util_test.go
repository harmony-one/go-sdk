package transaction

import (
	"math/big"
	"testing"
)

const DefaultPassphrase = "harmony-one"

var TESTNET = big.NewInt(2)

func TestSignTx(t *testing.T) {
	expected := "0xf86d80808252088080940a7d4bbd75eecaf11f8c891ed47269006bf91dc389056bc75e2d631000008361626328a0e002a9a03e3dd69e7eca8f811b656d793b1670df89820439df7bc7a12187eb27a059c4fd25fc08f0f01f5a8e6a995b980949f929c56b6e7bcbe391d47b33dbae14"
	nonce := big.NewInt(0).Uint64()
	gasLimit := big.NewInt(0).Uint64()
	to := "one1pf75h0t4am90z8uv3y0dgunfqp4lj8wr3t5rsp"
	shardID := uint32(0)
	toShardID := uint32(0)
	amount := big.NewInt(100)
	gasPrice := big.NewInt(0)
	data := []byte("abc")

	prv := "fd416cb87dcf8ed187e85545d7734a192fc8e976f5b540e9e21e896ec2bc25c3"
	accountName := "ac2" // use a different account name next time

	tx, err := NewTx(nonce, gasLimit, to, shardID, toShardID, amount, gasPrice, data)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	raw, err := SignTx(tx, prv, DefaultPassphrase, accountName, TESTNET)
	if err != nil {
		t.Errorf("failed with error: %s", err)
	}
	if string(raw) != expected {
		t.Errorf("test failed generated raw bytes does not match %s and %s", raw, expected)
	}
}
