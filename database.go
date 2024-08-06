package core

import (
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const envDatabaseURL = "DATABASE_URL"

func Open() (*gorm.DB, error) {
	url := os.Getenv(envDatabaseURL)
	if url == "" {
		return nil, fmt.Errorf("%w: %s", ErrUnsetEnvironmentVar, envDatabaseURL)
	}
	db, err := gorm.Open(
		postgres.Open(url),
		&gorm.Config{ //nolint:exhaustruct
			CreateBatchSize: 256,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	err = db.AutoMigrate(&Log{}) //nolint:exhaustruct
	if err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}

func Paginate(query *gorm.DB, page int, pageSize int) (*gorm.DB, error) {
	// Validate input.
	if page < 0 {
		return nil, ErrNegativePage
	}

	// A zero or negative `pageSize` is a special case where the
	// query is not paginated.
	if pageSize <= 0 {
		return query, nil
	}

	// Calculate the offset.
	offset := page * pageSize

	// Paginate.
	return query.Limit(pageSize).Offset(offset), nil
}

func SelectLogs(db *gorm.DB, contract string, topic string, page, pageSize int) ([]Log, error) {
	// Validate input.
	if !ValidateAddress(contract) {
		return nil, makeErrorHex(ErrInvalidContractAddress, contract)
	}
	if topic != "" && !ValidateTopic(topic) {
		return nil, makeErrorHex(ErrInvalidTopic, topic)
	}

	// Prepare query.
	query := db.Select("*").Where("address = ?", strings.ToLower(contract))
	if topic != "" {
		query = query.Where("topic0 = ?", strings.ToLower(topic))
	}
	query = query.Order("block_number desc, index desc")
	query, err := Paginate(query, page, pageSize)
	if err != nil {
		return nil, err
	}

	// Execute query.
	var xs []Log
	result := query.Find(&xs)
	if result.Error != nil {
		return xs, fmt.Errorf("find: %w", result.Error)
	}

	return xs, nil
}
