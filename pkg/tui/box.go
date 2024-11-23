package tui

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
	sLen := lenOfLargestLine(s)

	res := topLeft + strings.Repeat(horizontal, sLen+2) + topRight + "\n"
	for _, line := range strings.Split(s, "\n") {
		res += vertical + " " + line + strings.Repeat(" ", sLen-len(line)) + " " + vertical + "\n"
	}
	res += bottomLeft + strings.Repeat(horizontal, sLen+2) + bottomRight + "\n"

	return res
}

func lenOfLargestLine(s string) int {
	lines := strings.Split(s, "\n")

	max := 0
	for _, line := range lines {
		if len(line) > max {
			max = len(line)
		}
	}

	return max
}
