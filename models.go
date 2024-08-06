package core

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func prepareHex(hex string) string {
	return strings.ToLower(SanitizeHex(hex))
}

// +--------+
// | Events |
// +--------+

type Events struct {
	Contract string `gorm:"uniqueContract:idx_addr;not null"`
	Events  string `gorm:"not null"`
}

func Serialize(address map[string]abi.Event) Events {
	return Events{"", ""}
}

// +-----+
// | Log |
// +-----+

// Unique constraings:
//   - idx_logs_abi: (address,block_number,index)
//   - idx_logs_hi: (tx_hash,index)
type Log struct {
	ID          uint64 `gorm:"primaryKey"`
	Address     string `gorm:"uniqueIndex:idx_logs_abi;not null"`
	Topic0      string ``
	Topic1      string ``
	Topic2      string ``
	Topic3      string ``
	Data        string ``
	BlockNumber uint64 `gorm:"uniqueIndex:idx_logs_abi;not null"`
	TxHash      string `gorm:"uniqueIndex:idx_logs_hi;not null"`
	// Index of the transaction in the block.
	TxIndex uint `gorm:"not null"`
	// Index of the log in the block.
	Index uint `gorm:"uniqueIndex:idx_logs_abi;uniqueIndex:idx_logs_hi;not null"`
}

func FromGethLog(log types.Log) Log {
	topics := make([]string, 4)
	for i, t := range log.Topics {
		topics[i] = t.Hex()
	}
	return Log{
		ID:          0,
		Address:     prepareHex(log.Address.Hex()),
		Topic0:      prepareHex(topics[0]),
		Topic1:      prepareHex(topics[1]),
		Topic2:      prepareHex(topics[2]),
		Topic3:      prepareHex(topics[3]),
		Data:        prepareHex(common.Bytes2Hex(log.Data)),
		BlockNumber: log.BlockNumber,
		TxHash:      prepareHex(log.TxHash.Hex()),
		TxIndex:     log.TxIndex,
		Index:       log.Index,
	}
}

func adaptLogs(xs []types.Log) []Log {
	ys := make([]Log, len(xs))
	for i, x := range xs {
		ys[i] = FromGethLog(x)
	}
	return ys
}

func (o Log) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Log:\n")
	fmt.Fprintf(&b, "\tID          : %d\n", o.ID)
	fmt.Fprintf(&b, "\tAddress     : %s\n", o.Address)
	fmt.Fprintf(&b, "\tTopic0      : %s\n", o.Topic0)
	fmt.Fprintf(&b, "\tTopic1      : %s\n", o.Topic1)
	fmt.Fprintf(&b, "\tTopic2      : %s\n", o.Topic2)
	fmt.Fprintf(&b, "\tTopic3      : %s\n", o.Topic3)
	fmt.Fprintf(&b, "\tData        : %s\n", o.Data)
	fmt.Fprintf(&b, "\tBlockNumber : %d\n", o.BlockNumber)
	fmt.Fprintf(&b, "\tTxHash      : %s\n", o.TxHash)
	fmt.Fprintf(&b, "\tTxIndex     : %d\n", o.TxIndex)
	fmt.Fprintf(&b, "\tIndex       : %d\n", o.Index)
	return b.String()
}
