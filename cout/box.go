package cout

import (
	"strings"
)

// Box
const (
	horizontal = "─"
	vertical   = "│"

	topLeftSquare     = "┌"
	topRightSquare    = "┐"
	bottomLeftSquare  = "└"
	bottomRightSquare = "┘"

	topLeftRound     = "╭"
	topRightRound    = "╮"
	bottomLeftRound  = "╰"
	bottomRightRound = "╯"

	// ---

	horizontalDouble = "═"
	verticalDouble   = "║"

	topLeftDouble     = "╔"
	topRightDouble    = "╗"
	bottomLeftDouble  = "╚"
	bottomRightDouble = "╝"
)

func Box(s string) string {
	return box(s, topLeftSquare, topRightSquare, bottomLeftSquare, bottomRightSquare, horizontal, vertical)
}

func BoxRoundCorner(s string) string {
	return box(s, topLeftRound, topRightRound, bottomLeftRound, bottomRightRound, horizontal, vertical)
}

func BoxDouble(s string) string {
	return box(s, topLeftDouble, topRightDouble, bottomLeftDouble, bottomRightDouble, horizontalDouble, verticalDouble)
}

func box(s, topLeft, topRight, bottomLeft, bottomRight, horizontal, vertical string) string {
	sLen := len(s)

	res := topLeft + strings.Repeat(horizontal, sLen+2) + topRight + "\n"
	res += vertical + " " + s + " " + vertical + "\n"
	res += bottomLeft + strings.Repeat(horizontal, sLen+2) + bottomRight + "\n"

	return res
}
