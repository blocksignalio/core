package core

import "fmt"

func makeErrorHex(err error, address string) error {
	return fmt.Errorf("%w: %s", err, SanitizeHex(address))
}
