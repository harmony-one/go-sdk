package account

import (
	"testing"

	"github.com/harmony-one/go-sdk/pkg/store"
)

func TestAccountGetsRemoved(t *testing.T) {
	tests := []struct {
		Name     string
		Expected bool
	}{
		{"foobar", false},
	}

	for _, test := range tests {
		acc := Creation{
			Name:            test.Name,
			Passphrase:      "",
			Mnemonic:        "",
			HdAccountNumber: nil,
			HdIndexNumber:   nil,
		}

		err := CreateNewLocalAccount(&acc)

		if err != nil {
			t.Errorf(`RemoveAccount("%s") failed to create keystore account %v`, test.Name, err)
		}

		exists := store.DoesNamedAccountExist(test.Name)

		if !exists {
			t.Errorf(`RemoveAccount("%s") account should exist - but it can't be found`, test.Name)
		}

		RemoveAccount(test.Name)
		exists = store.DoesNamedAccountExist(test.Name)

		if exists != test.Expected {
			t.Errorf(`RemoveAccount("%s") returned %v, expected %v`, test.Name, exists, test.Expected)
		}
	}
}
