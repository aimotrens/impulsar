package cout

// Farben
const (
	yellow      = "\033[33m"
	orange      = "\033[38;5;208m"
	red         = "\033[31m"
	green       = "\033[32m"
	blue        = "\033[34m"
	lightBlue   = "\033[94m"
	cyan        = "\033[36m"
	gray        = "\033[37m"
	darkGray    = "\033[90m"
	lightYellow = "\033[93m"

	reset = "\033[0m"
)

func colorize(s, color string) string {
	return color + s + reset
}

func Yellow(s string) string {
	return colorize(s, yellow)
}

func LightYellow(s string) string {
	return colorize(s, lightYellow)
}

func Orange(s string) string {
	return colorize(s, orange)
}

func Red(s string) string {
	return colorize(s, red)
}

func Green(s string) string {
	return colorize(s, green)
}

func Blue(s string) string {
	return colorize(s, blue)
}

func Cyan(s string) string {
	return colorize(s, cyan)
}

func DarkGray(s string) string {
	return colorize(s, darkGray)
}

func Gray(s string) string {
	return colorize(s, gray)
}
