package cablemodemutil

import (
	"fmt"
	"math"
	"testing"
)

var validParseUint32Tests = []struct {
	name      string
	str       string
	hasSuffix bool
	suffix    string
	desc      string
	want      uint32
}{
	{
		name:      "Zero without suffix",
		str:       "0",
		hasSuffix: false,
		suffix:    "",
		desc:      "desc",
		want:      0,
	},
	{
		name:      "MaxUint32 without suffix",
		str:       fmt.Sprintf("%d", math.MaxUint32),
		hasSuffix: false,
		suffix:    "",
		desc:      "desc",
		want:      math.MaxUint32,
	},
	{
		name:      "Valid uint32 without suffix",
		str:       "67245",
		hasSuffix: false,
		suffix:    "",
		desc:      "desc",
		want:      67245,
	},
	{
		name:      "Zero with suffix",
		str:       "0Foo",
		hasSuffix: true,
		suffix:    "Foo",
		desc:      "desc",
		want:      0,
	},
	{
		name:      "MaxUint32 with suffix",
		str:       fmt.Sprintf("%d Bar", math.MaxUint32),
		hasSuffix: true,
		suffix:    " Bar",
		desc:      "desc",
		want:      math.MaxUint32,
	},
	{
		name:      "Valid uint32 with suffix",
		str:       "67245x",
		hasSuffix: true,
		suffix:    "x",
		desc:      "desc",
		want:      67245,
	},
}

func TestParseUint32(t *testing.T) {
	for _, tc := range validParseUint32Tests {
		if result, gotErr := parseUint32(tc.str, tc.hasSuffix, tc.suffix, tc.desc); nil != gotErr || result != tc.want {
			t.Errorf(
				"%q: parseUint32(%q, %t, %q, %q) = (%d, %s) want: (%d, nil)",
				tc.name,
				tc.str,
				tc.hasSuffix,
				tc.suffix,
				tc.desc,
				result,
				gotErr,
				tc.want,
			)
		}
	}
}
