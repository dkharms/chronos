package parser

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	//go:embed testdata/gotool.txt
	gotoolOutput string
)

func TestParseGoTool(t *testing.T) {
	actual := NewGoParser(strings.NewReader(gotoolOutput)).Parse()
	// TODO(dkharms): Write better tests equality check.
	// Now I am just being too lazy...
	assert.Len(t, actual, 5)
}
