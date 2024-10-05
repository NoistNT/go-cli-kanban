package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

// AddHelpKeys adds the help keys
func AddHelpKeys(defaultList *list.Model) {
	defaultList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp("a", "add task"),
			),
			key.NewBinding(
				key.WithKeys("e"),
				key.WithHelp("e", "edit task"),
			),
		}
	}
	defaultList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp("a", "add a task"),
			),
			key.NewBinding(
				key.WithKeys("e"),
				key.WithHelp("e", "edit a task"),
			),
		}
	}
}
