package core_test

import (
	"testing"

	"github.com/blocksignalio/core"
)

func TestSignatures(t *testing.T) {
	t.Parallel()

	have := []string{
		core.Approval.ID.Hex(),
		core.Transfer.ID.Hex(),
	}

	want := []string{
		"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925",
		"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
	}

	for i := range have {
		if have[i] != want[i] {
			t.Errorf("signature: have=%s want=%s", have[i], want[i])
		}
	}
}
