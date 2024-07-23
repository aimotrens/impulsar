package cout

// Font Decorations
const (
	bold          = "\033[1m"
	underline     = "\033[4m"
	inverse       = "\033[7m"
	strikethrough = "\033[9m"
)

func Bold(s string) string {
	return bold + s + reset
}

func Underline(s string) string {
	return underline + s + reset
}

func Inverse(s string) string {
	return inverse + s + reset
}

func Strikethrough(s string) string {
	return strikethrough + s + reset
}
