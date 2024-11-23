package tui_test

import (
	"testing"

	"github.com/aimotrens/impulsar/pkg/tui"
	"github.com/stretchr/testify/assert"
)

func Test_Box(t *testing.T) {
	expected := `┌────────┐
│ Test   │
│ Foobar │
└────────┘
`

	assert.Equal(t, expected, tui.Box("Test\nFoobar"))
}

func Test_BoxRoundCorner(t *testing.T) {
	expected := `╭────────╮
│ Test   │
│ Foobar │
╰────────╯
`

	assert.Equal(t, expected, tui.BoxRoundCorner("Test\nFoobar"))
}

func Test_BoxDouble(t *testing.T) {
	expected := `╔════════╗
║ Test   ║
║ Foobar ║
╚════════╝
`

	assert.Equal(t, expected, tui.BoxDouble("Test\nFoobar"))
}
