package main

import "github.com/charmbracelet/lipgloss"

// StyledListView returns the styles for the list view
func StyledListView(focused int, views ...string) string {
	var result []string
	for i, view := range views {
		if i == focused {
			result = append(result, focusedStyle.Render(view))
		} else {
			result = append(result, columnStyle.Render(view))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, result...)
}

// StyledInputTaskView returns the styles for the text input in task adding/editing view
func (m *model) StyledInputTaskView(msg string) string {
	return lipgloss.NewStyle().Padding(1).Render(""+msg+": \n\n"+m.textInput.View()) + m.ShowHelp()
}

// ShowHelp shows the help screen when the user is in task-adding/editing mode
func (m *model) ShowHelp() string {
	return helpStyle.Render("Press 'enter' to confirm, 'esc' to cancel, 'ctrl+c' to quit")
}
