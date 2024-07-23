package cout

import (
	"strings"
)

// Box
const (
	TopLeft     = "┌"
	TopRight    = "┐"
	BottomLeft  = "└"
	BottomRight = "┘"
	Horizontal  = "─"
	Vertical    = "│"
)

func H1(s string) string {
	sLen := len(s)

	res := TopLeft + strings.Repeat(Horizontal, sLen+2) + TopRight + "\n"
	res += Vertical + " " + s + " " + Vertical + "\n"
	res += BottomLeft + strings.Repeat(Horizontal, sLen+2) + BottomRight + "\n"

	return res
}
