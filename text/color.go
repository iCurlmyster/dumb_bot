package text

import "fmt"

// Red - Red text for terminal
func Red(s string) string {
	return fmt.Sprintf("\033[0;31m%s\033[0m", s)
}

// Green - Green text for terminal
func Green(s string) string {
	return fmt.Sprintf("\033[0;32m%s\033[0m", s)
}

// Brown - Brown text for terminal
func Brown(s string) string {
	return fmt.Sprintf("\033[0;33m%s\033[0m", s)
}

// Blue - Blue text for terminal
func Blue(s string) string {
	return fmt.Sprintf("\033[0;34m%s\033[0m", s)
}

// Magenta - Magenta text for terminal
func Magenta(s string) string {
	return fmt.Sprintf("\033[0;35m%s\033[0m", s)
}

// Cyan - Cyan text for terminal
func Cyan(s string) string {
	return fmt.Sprintf("\033[0;36m%s\033[0m", s)
}

// White - White text for terminal
func White(s string) string {
	return fmt.Sprintf("\033[0;37m%s\033[0m", s)
}

// Yellow - Yellow text for terminal
func Yellow(s string) string {
	return fmt.Sprintf("\033[1;33m%s\033[0m", s)
}

// LightBlue - LightBlue text for terminal
func LightBlue(s string) string {
	return fmt.Sprintf("\033[1;34m%s\033[0m", s)
}
