package core

import (
	"errors"
	"hash"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

type KeccakState interface {
	hash.Hash
	Read(p []byte) (int, error)
}

var (
	errHasherHasNoRead = errors.New("hasher has no read function")
	errReadFailed      = errors.New("read failed")
	errIncorrectLength = errors.New("produced hash has an incorrect length")
)

func Keccak256Hash(data ...[]byte) (common.Hash, error) {
	var h common.Hash

	d, ok := sha3.NewLegacyKeccak256().(KeccakState)
	if !ok {
		return h, errHasherHasNoRead
	}

	for _, b := range data {
		d.Write(b)
	}

	n, err := d.Read(h[:])
	if err != nil {
		return h, errReadFailed
	}
	if n != 32 {
		return h, errIncorrectLength
	}

	return h, nil
}
