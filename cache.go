package core

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
)

// (address string) (map[string]abi.Event, error) {
func GetContractEventsCached(address string) (map[string]abi.Event, error) {
	if !ValidateAddress(address) {
		return nil, makeErrorHex(ErrInvalidContractAddress, address)
	}
	return GetContractEvents(address)
}
