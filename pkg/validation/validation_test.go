package validation

import (
	"testing"

	"github.com/harmony-one/go-sdk/pkg/sharding"
)

func TestIsValidAddress(t *testing.T) {
	tests := []struct {
		str string
		exp bool
	}{
		{"one1ay37rp2pc3kjarg7a322vu3sa8j9puahg679z3", true},
		{"0x7c41E0668B551f4f902cFaec05B5Bdca68b124CE", true},
		{"onefoofoo", false},
		{"0xbarbar", false},
		{"dsasdadsasaadsas", false},
		{"32312123213213212321", false},
	}

	for _, test := range tests {
		err := ValidateAddress(test.str)
		valid := false

		if err == nil {
			valid = true
		}

		if valid != test.exp {
			t.Errorf(`ValidateAddress("%s") returned %v, expected %v`, test.str, valid, test.exp)
		}
	}
}

func TestIsValidShard(t *testing.T) {
	if err := ValidateNodeConnection("http://localhost:9500"); err != nil {
		t.Skip()
	}
	s, _ := sharding.Structure("http://localhost:9500")

	tests := []struct {
		shardID uint32
		exp     bool
	}{
		{0, true},
		{1, true},
		{98, false},
		{99, false},
	}

	for _, test := range tests {
		valid := ValidShardID(test.shardID, uint32(len(s)))

		if valid != test.exp {
			t.Errorf("ValidShardID(%d) returned %v, expected %v", test.shardID, valid, test.exp)
		}
	}
}
