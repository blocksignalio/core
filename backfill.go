package core

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

func retrieveFromBlock(ctx context.Context, db *gorm.DB, contract string) (uint64, error) {
	var last []Log
	result := db.Distinct("address", "block_number", "index").
		Where("address = ?", prepareHex(contract)).
		Order("block_number desc, index desc").
		Limit(1).
		Find(&last)
	if result.Error != nil {
		return 0, fmt.Errorf("find: %w", result.Error)
	}
	if len(last) > 0 {
		fromBlock := last[0].BlockNumber + 1
		return fromBlock, nil
	}
	creation, err := GetContractCreation1(contract)
	if err != nil {
		return 0, err
	}
	return GetTransactionBlock(ctx, creation.TxHash)
}

// TODO: Query up to head-64 (?) to retrieve only finalized logs?  And
// then check types.Log.Removed to confirm everything was OK!
func BackfillLogs(ctx context.Context, db *gorm.DB, contract string) error {
	if !ValidateAddress(contract) {
		return makeErrorHex(ErrInvalidContractAddress, contract)
	}

	fromBlock, err := retrieveFromBlock(ctx, db, contract)
	if err != nil {
		return err
	}

	for {
		toBlock, xs, err := QueryLogs(ctx, fromBlock, contract)
		if err != nil {
			return err
		}

		ys := adaptLogs(xs)
		db.Create(&ys)

		// fmt.Printf(
		// 	"Backfill: querying: fromBlock=%d toBlock=%d result=%d\n",
		// 	fromBlock,
		// 	toBlock,
		// 	len(ys),
		// )

		if len(xs) > 0 {
			fromBlock = toBlock
		} else {
			break
		}
	}
	return nil
}
