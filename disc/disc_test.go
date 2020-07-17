package disc

import (
	"testing"
)

func TestParseArgs(t *testing.T) {
	inputs := []string{"!chip echo dumpling"}
	expected := [][]string{
		[]string{
			"!chip", "echo", "dumpling",
		},
	}
	for i, in := range inputs {
		out := parseArgs(in)
		for j, o := range out {
			if o != expected[i][j] {
				t.Error("unexpected output")
			}
		}
	}
}
