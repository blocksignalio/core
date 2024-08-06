package core_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/blocksignalio/core"
	"github.com/google/go-cmp/cmp"
)

func TestGetContractABI(t *testing.T) {
	t.Parallel()

	const (
		weth = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
		want = ("" +
			"event Approval(address indexed src, address indexed guy, uint256 wad) " +
			"event Deposit(address indexed dst, uint256 wad) " +
			"event Transfer(address indexed src, address indexed dst, uint256 wad) " +
			"event Withdrawal(address indexed src, uint256 wad)")
	)

	abi, err := core.GetContractABI(weth)
	if err != nil {
		t.Fatal(err)
	}
	events := make([]string, 0, 4)
	for _, e := range abi.Events {
		events = append(events, e.String())
	}
	sort.Strings(events)
	have := strings.Join(events, " ")
	if diff := cmp.Diff(have, want); diff != "" {
		t.Error(diff)
	}
}
