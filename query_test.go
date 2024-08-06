package core_test

import (
	"context"
	"testing"

	"github.com/blocksignalio/core"
)

func TestQuery(t *testing.T) {
	t.Parallel()

	const (
		fromBlock = 17899693
		weth      = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"

		wantBlocks = 130
		wantLogs   = 9976
	)

	toBlock, logs, err := core.QueryLogs(context.Background(), fromBlock, weth)
	if err != nil {
		t.Error(err)
	}

	if haveBlocks := toBlock - fromBlock; haveBlocks != wantBlocks {
		t.Errorf("blocks: have=%d want=%d", haveBlocks, wantBlocks)
	}

	if haveLogs := len(logs); haveLogs != wantLogs {
		t.Errorf("logs: have=%d want=%d", haveLogs, wantLogs)
	}
}
