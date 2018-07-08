package bit

import (
	"strings"
	"testing"
	"time"
)

func TestFormatElapsed(t *testing.T) {
	duration := func(duration string) time.Duration {
		d, err := time.ParseDuration(duration)
		if err != nil {
			panic(err)
		}
		return d
	}

	var tests = []struct {
		name  string
		value time.Duration
		want  string
	}{
		{name: "zero", value: duration("0s"), want: "0:00:00"},
		{name: "seconds", value: duration("32s"), want: "0:00:32"},
		{name: "minutes", value: duration("32m"), want: "0:32:00"},
		{name: "hours", value: duration("32h"), want: "32:00:00"},
		{name: "full", value: duration("1h23m45s"), want: "1:23:45"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			have := strings.TrimSpace(formatElapsed(test.value))
			if test.want != have {
				t.Errorf("\n want: %v \n have: %v", test.want, have)
			}
		})
	}

}
