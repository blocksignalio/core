package core

import (
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
)

var (
	ErrInvalidContractAddress = errors.New("invalid contract address")
	ErrUnsetEnvironmentVar    = errors.New("environment variable not set")
	ErrNegativePage           = errors.New("page cannot be negative")
	ErrInvalidTopic           = errors.New("invalid topic")

	ErrInvalidResponse     = errors.New("invalid response")
	ErrInvalidResponseBody = errors.New("invalid response body")
)

func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

func SanitizeHex(address string) string {
	address = strings.TrimSpace(address)
	prefix := strings.HasPrefix(address, "0x") || strings.HasPrefix(address, "0X")
	if prefix {
		address = address[2:]
	}
	re := regexp.MustCompile(`[0-9a-fA-F]+`)
	xs := re.FindAllString(address, -1)
	ss := strings.Join(xs, "")
	if ss == "" {
		return ""
	}
	return "0x" + ss
}

func ValidateAddress(address string) bool {
	if has0xPrefix(address) {
		address = address[2:]
	}
	if len(address) != 40 {
		return false
	}
	_, err := hex.DecodeString(address)
	return err == nil
}

func ValidateTopic(topic string) bool {
	if has0xPrefix(topic) {
		topic = topic[2:]
	}
	if len(topic) != 64 {
		return false
	}
	_, err := hex.DecodeString(topic)
	return err == nil
}
