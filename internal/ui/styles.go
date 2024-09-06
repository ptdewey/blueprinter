package ui

import "github.com/charmbracelet/lipgloss"

// TODO: figure out how styling and style application works

var DocStyle = lipgloss.NewStyle().Margin(1, 2).
	Foreground(lipgloss.Color("default"))

var HighlightStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("red")).
	Foreground(lipgloss.Color("white"))

var CursorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("green"))
