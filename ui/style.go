package ui

import (
	"github.com/muesli/termenv"
)

var (
	primary   = termenv.ColorProfile().Color("205")
	secondary = termenv.ColorProfile().Color("#89F0CB")
	gray      = termenv.ColorProfile().Color("#626262")
	red       = termenv.ColorProfile().Color("#ED567A")
)

func Bold(s string) string {
	return termenv.String(s).Bold().String()
}

func Italic(s string) string {
	return termenv.String(s).Italic().String()
}

func PrimaryForeground(s string) string {
	return termenv.String(s).Foreground(primary).String()
}

func SecondaryForeground(s string) string {
	return termenv.String(s).Foreground(secondary).String()
}

func GrayForeground(s string) string {
	return termenv.String(s).Foreground(gray).String()
}

func RedForeground(s string) string {
	return termenv.String(s).Foreground(red).String()
}
