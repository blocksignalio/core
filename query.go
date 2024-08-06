package core

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-1474.md
const (
	codeLimitExceeded = -32005
	envEthereumNode   = "ETHEREUM_NODE"
)

//nolint:gochecknoglobals
var powers = []uint64{
	0x200000,
	0x100000,
	0x080000,
	0x040000,
	0x020000,
	0x010000,
	0x008000,
	0x004000,
	0x002000,
	0x001000,
	0x000800,
	0x000400,
	0x000200,
	0x000100,
	0x000080,
	0x000040,
	0x000020,
	0x000010,
	0x000008,
	0x000004,
	0x000002,
	0x000001,
	0x000000,
}

func extractBound(err error) (uint64, bool) {
	// Check error code.
	a, ok := err.(rpc.Error) //nolint:errorlint
	if !ok {
		return 0, false
	}

	if a.ErrorCode() != codeLimitExceeded {
		return 0, false
	}

	// Get error data.
	b, ok := err.(rpc.DataError) //nolint:errorlint
	if !ok {
		return 0, false
	}

	var c any = b.ErrorData() //nolint:stylecheck
	if c == nil {
		return 0, false
	}

	d, ok := c.(map[string]any)
	if !ok {
		return 0, false
	}

	// Get data["to"] as string.
	e, ok := d["to"]
	if !ok {
		return 0, false
	}

	f, ok := e.(string)
	if !ok {
		return 0, false
	}

	// If prefixed with 0x, strip it.
	g := f
	if strings.HasPrefix(f, "0x") {
		g = f[2:]
	}

	// Parse data["to"].
	x, err := strconv.ParseUint(g, 16, 64)
	if err != nil {
		return 0, false
	}

	return x, true
}

func makeClient() (*ethclient.Client, error) {
	node := os.Getenv(envEthereumNode)
	if node == "" {
		return nil, fmt.Errorf("%w: %s", ErrUnsetEnvironmentVar, envEthereumNode)
	}

	client, err := ethclient.Dial(node)
	if err != nil {
		return client, fmt.Errorf("dial: %w", err)
	}

	return client, nil
}

// Given a transaction, return the block it's included in.
func GetTransactionBlock(ctx context.Context, tx string) (uint64, error) {
	client, err := makeClient()
	if err != nil {
		return 0, fmt.Errorf("makeClient: %w", err)
	}

	h := common.HexToHash(tx)
	receipt, err := client.TransactionReceipt(ctx, h)
	if err != nil {
		return 0, fmt.Errorf("transaction by hash: %w", err)
	}
	blockNumber := receipt.BlockNumber.Uint64()
	return blockNumber, nil
}

// QueryLogs returns (toBlock, logs, err), where `logs` is all logs in the
// range of [fromBlock, toBlock).
func QueryLogs(ctx context.Context, fromBlock uint64, contract string) (uint64, []types.Log, error) {
	if !ValidateAddress(contract) {
		return fromBlock, nil, makeErrorHex(ErrInvalidContractAddress, contract)
	}

	client, err := makeClient()
	if err != nil {
		return fromBlock, nil, fmt.Errorf("makeClient: %w", err)
	}

	head, err := client.BlockNumber(ctx)
	if err != nil {
		return fromBlock, nil, fmt.Errorf("request head: %w", err)
	}

	address := common.HexToAddress(contract)

	toBlock := fromBlock + 2*powers[0]
	for _, step := range powers {
		if toBlock > head {
			toBlock = head
		}

		if fromBlock >= toBlock {
			return 0, nil, nil
		}

		query := ethereum.FilterQuery{
			BlockHash: nil,
			FromBlock: big.NewInt(int64(fromBlock)),
			ToBlock:   big.NewInt(int64(toBlock)),
			Addresses: []common.Address{address},
			Topics: [][]common.Hash{
				{},
			},
		}

		logs, err := client.FilterLogs(ctx, query)
		if err == nil {
			// Success!
			return toBlock + 1, logs, nil
		}

		// There's an error.  See if we can extract a bounding
		// block from it.
		bound, ok := extractBound(err)
		if ok { //nolint:gocritic
			toBlock = bound
		} else if ctx.Err() == nil {
			toBlock = fromBlock + step
		} else {
			return fromBlock, nil, fmt.Errorf("context done: %w", ctx.Err())
		}
	}

	panic("unreachable")
}
