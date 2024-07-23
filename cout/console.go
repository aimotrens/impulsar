package cout

import "strings"

type FontModifier func(s string) string

func Default(s string) string {
	return s
}

func Multi(fm ...FontModifier) FontModifier {
	return func(s string) string {
		for _, f := range fm {
			s = f(s)
		}
		return s
	}
}

func FormattingArea(fm FontModifier, f func(*strings.Builder)) string {
	builder := new(strings.Builder)
	f(builder)
	return fm(builder.String())
}

func FormatByIndex(s string, start, end int, fm FontModifier) string {
	if start < 0 {
		start = 0
	}

	if end > len(s) {
		end = len(s)
	}

	return s[:start] + fm(s[start:end]) + s[end:]
}

func FormatByText(s, text string, fm FontModifier) string {
	if len(text) == 0 {
		return s
	}

	nextIndex := func(s, text string, start int) int {
		for i := start; i < len(s); i++ {
			if strings.HasPrefix(s[i:], text) {
				return i
			}
		}

		return -1
	}

	start := 0
	for {
		start = nextIndex(s, text, start)
		if start == -1 {
			break
		}

		oldLen := len(s)
		s = FormatByIndex(s, start, start+len(text), fm)

		start += len(text) + (len(s) - oldLen)
	}

	return s
}
