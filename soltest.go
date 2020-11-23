// soltest provides helpful utilties to test solidity smart contracts
// with go code.
package soltest

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// TestAccount maintains an ethereum account private key used to
// authenticate transactions.
type TestAccount struct {
	Addr common.Address
	Priv *ecdsa.PrivateKey
	Auth *bind.TransactOpts
}

// TestChain is a simulated, in-memory, ethereum blockchain used for
// testing.
// Use type alias?
type TestChain struct {
	*backends.SimulatedBackend
}

// Succeed accepts the output of a typical abigen trasactor method,
// and returns whether the given transaction was successfully
// committed into the TestChain.
func (c *TestChain) Succeed(txn *types.Transaction, err error) bool {
	if err != nil {
		return false
	}
	c.Commit()
	r, err := c.TransactionReceipt(nil, txn.Hash())
	if err != nil {
		return false
	}

	return r.Status == 1
}

// New returns a TestChain, and a slice of TestAccounts which all have
// a started eth balance on the TestChain.
func New() (TestChain, []TestAccount) {
	var testAccounts []TestAccount
	genesis := make(core.GenesisAlloc)
	for _, pk := range []string{
		"1010101010101010101010101010101010101010101010101010101010101010",
		"1111111111111111111111111111111111111111111111111111111111111111",
		"2222222222222222222222222222222222222222222222222222222222222222",
		"3333333333333333333333333333333333333333333333333333333333333333",
		"4444444444444444444444444444444444444444444444444444444444444444",
		"5555555555555555555555555555555555555555555555555555555555555555",
		"6666666666666666666666666666666666666666666666666666666666666666",
		"7777777777777777777777777777777777777777777777777777777777777777",
		"8888888888888888888888888888888888888888888888888888888888888888",
		"9999999999999999999999999999999999999999999999999999999999999999",
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
		"dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd",
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
	} {
		key, _ := crypto.HexToECDSA(pk)

		genesis[crypto.PubkeyToAddress(key.PublicKey)] = core.GenesisAccount{
			Balance: big.NewInt(1e18), // change to 1 eth?
		}

		// For testing, gas price will be 0 to keep balance inquiries easy.
		t := bind.NewKeyedTransactor(key)
		t.GasPrice = big.NewInt(0)

		testAccounts = append(testAccounts, TestAccount{
			Addr: crypto.PubkeyToAddress(key.PublicKey),
			Priv: key,
			Auth: t,
		})
	}

	return TestChain{backends.NewSimulatedBackend(genesis, 0)}, testAccounts
}
