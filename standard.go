package core

// Definitions of standard events.

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
)

const (
	indexed = true
)

func makeEvent(name string, desc ...any) abi.Event {
	if len(desc)%3 != 0 {
		panic("each argument is described by a (name,type,indexed) tuple")
	}
	inputs := make([]abi.Argument, 0, len(desc)/3)
	for i := 0; i < len(desc); i += 3 {
		n := desc[i+0].(string)   //nolint:forcetypeassert
		t := desc[i+1].(abi.Type) //nolint:forcetypeassert
		x := desc[i+2].(bool)     //nolint:forcetypeassert
		inputs = append(inputs, abi.Argument{Name: n, Type: t, Indexed: x})
	}
	return abi.NewEvent(name, name, false, inputs)
}

//nolint:gochecknoglobals
var (
	Uint256, _    = abi.NewType("uint256", "", nil)
	Uint32, _     = abi.NewType("uint32", "", nil)
	Uint16, _     = abi.NewType("uint16", "", nil)
	String, _     = abi.NewType("string", "", nil)
	Bool, _       = abi.NewType("bool", "", nil)
	Bytes, _      = abi.NewType("bytes", "", nil)
	Bytes32, _    = abi.NewType("bytes32", "", nil)
	Address, _    = abi.NewType("address", "", nil)
	Uint64Arr, _  = abi.NewType("uint64[]", "", nil)
	AddressArr, _ = abi.NewType("address[]", "", nil)
	Int8, _       = abi.NewType("int8", "", nil)

	// // Special types
	// Uint32Arr2, _       = NewType("uint32[2]", "", nil)
	// Uint64Arr2, _       = NewType("uint64[2]", "", nil)
	// Uint256Arr, _       = NewType("uint256[]", "", nil)
	// Uint256Arr2, _      = NewType("uint256[2]", "", nil)
	// Uint256Arr3, _      = NewType("uint256[3]", "", nil)
	// Uint256ArrNested, _ = NewType("uint256[2][2]", "", nil)
	// Uint8ArrNested, _   = NewType("uint8[][2]", "", nil)
	// Uint8SliceNested, _ = NewType("uint8[][]", "", nil).

	Approval = makeEvent("Approval",
		"owner", Address, indexed,
		"spender", Address, indexed,
		"value", Uint256, false)

	Transfer = makeEvent("Transfer",
		"from", Address, indexed,
		"to", Address, indexed,
		"amount", Uint256, false)
)

/*
https://etherscan.io/address/0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2#code
    event  Approval(address indexed src, address indexed guy, uint wad);
    event  Transfer(address indexed src, address indexed dst, uint wad);
    event  Deposit(address indexed dst, uint wad);
    event  Withdrawal(address indexed src, uint wad);

TODO: https://www.4byte.directory/event-signatures/?page=3
*/
