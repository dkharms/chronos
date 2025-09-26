package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGoTool(t *testing.T) {
	actual := NewGoParser(strings.NewReader(GotoolOutput)).Parse()
	// TODO(dkharms): Write better tests equality check.
	// Now I am just being too lazy...
	assert.Len(t, actual, 5)
}
