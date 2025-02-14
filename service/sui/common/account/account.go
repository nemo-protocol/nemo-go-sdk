package account

import (
	"encoding/hex"
	"fmt"
	"github.com/coming-chat/go-sui/v2/account"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"strings"
)

func NewAccountPrivateKey(privateKeyHex string) *account.Account {
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")

	// decode privateKey
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		fmt.Printf("decode privateKey error: %v\n", err)
		return nil
	}

	// create sign scheme
	scheme, err := sui_types.NewSignatureScheme(0)
	if err != nil {
		fmt.Printf("create sign scheme error: %v\n", err)
		return nil
	}

	// create account
	sender := account.NewAccount(scheme, privateKeyBytes)
	return sender
}
