package core_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/blocksignalio/core"
)

func TestSanitizeHex(t *testing.T) {
	t.Parallel()

	inputs := []string{
		"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		"0XC02AAA39B223FE8D0A0E5C4F27EAD9083C756CC2",
		"c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		"  0XC02AAA39B223FE8D0A0E5C4F27EAD9083C756CC2  ",
		"xx0XC02AAA39B223FE8D0A0E5C4F27EAD9083C756CC2xx",
		"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2F",
		"Whoever appeals to the law against his fellow man is ...",
		"; rm -rf ~ ; %x %p %s ; ( cd / && pwd )",
		`echo -e "\\e[31mHello!\\e[0m"`,
		"",
	}
	expected := []string{
		"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		"0xC02AAA39B223FE8D0A0E5C4F27EAD9083C756CC2",
		"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		"0xC02AAA39B223FE8D0A0E5C4F27EAD9083C756CC2",
		"0x0C02AAA39B223FE8D0A0E5C4F27EAD9083C756CC2",
		"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2F",
		"0xeeaeaeaaafee",
		"0xfcdd",
		"0xecee31ee0",
		"",
	}

	if len(inputs) != len(expected) {
		t.Fatal()
	}
	for i, want := range expected {
		have := core.SanitizeHex(inputs[i])
		if diff := cmp.Diff(want, have); diff != "" {
			t.Errorf("input='%s'\n%s", inputs[i], diff)
		}
	}
}
