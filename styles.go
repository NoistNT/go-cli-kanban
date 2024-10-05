package main

import "github.com/charmbracelet/lipgloss"

// Styles for the application
var (
	columnStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Width(100 * 2 / 3).
			Border(lipgloss.HiddenBorder())
	focusedStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Width(100 * 2 / 3).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)
