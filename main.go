package main

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type listItems []list.Item

type lists []list.Model

// MAIN MODEL
type model struct {
	lists        lists
	focused      int
	err          error
	isLoaded     bool
	addTaskMode  bool
	editTaskMode bool
	textInput    textinput.Model
}

// InitLists initializes the lists To Do, In Progress, and Done
func (m *model) initLists(width, height int) {
	// Create new lists with calculated width and height
	defaultList := list.New(listItems{}, list.NewDefaultDelegate(), width, height-4)
	defaultList.SetShowHelp(false)

	// Add additional help keys
	AddHelpKeys(&defaultList)

	// Initialize the lists for To Do, In Progress, and Done
	m.lists = make(lists, 3)
	for i := range m.lists {
		m.lists[i] = defaultList
	}

	// Seed lists with data
	m.seedLists()
}

// SeedList seeds a list with data
func (m *model) seedList(title string, status status, items listItems) {
	m.lists[status].Title = title
	m.lists[status].SetItems(items)
}

// SeedLists seeds the lists To Do, In Progress, and Done with data
func (m *model) seedLists() {
	// Seed To Do List
	m.seedList("To do", todo, todoListItems)

	// Seed In Progress List
	m.seedList("In Progress", inProgress, inProgressListItems)

	// Seed Done List
	m.seedList("Done", done, doneListItems)
}

// InitTextInput initializes the text input
func (m *model) initTextInput() {
	m.textInput = textinput.New()
	m.textInput.Placeholder = "Buy bubbletea"
	m.textInput.Focus()
	m.textInput.CharLimit = 140
	m.textInput.Width = 50
}

// Init method to initialize the model
func initialModel() *model {
	m := &model{err: nil}
	m.initTextInput()
	return m
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.isLoaded {
			m.initLists(msg.Width, msg.Height)
			m.isLoaded = true
		}
	case tea.KeyMsg:
		// Handle key inputs differently based on the current mode
		switch {
		case m.addTaskMode:
			// In task-adding mode, forward the key events to the text input
			switch msg.String() {
			case "enter":
				// Add the task and return to list view
				m.AddTask()
				return m, nil
			case "esc":
				// Cancel adding task and return to list view
				m.addTaskMode = false
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			default:
				// Forward all other key inputs (including 'a') to the text input
				var cmd tea.Cmd
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd
			}

		case m.editTaskMode:
			// In task-editing mode, forward the key events to the text input
			switch msg.String() {
			case "enter":
				// Add the task and return to list view
				m.EditTask()
				return m, nil
			case "esc":
				// Cancel editing task and return to list view
				m.addTaskMode = false
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			default:
				// Forward all other key inputs (including 'a') to the text input
				var cmd tea.Cmd
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd
			}

		default:
			// In list mode, handle navigation and toggling between modes
			switch msg.String() {
			case "a":
				// Toggle to task-adding mode
				m.textInput.SetValue("") // Clear the input field
				m.addTaskMode = true
				m.textInput.Focus()
				return m, nil
			case "e":
				// Toggle to edit mode
				m.editTaskMode = true
				m.textInput.SetValue(m.lists[m.focused].SelectedItem().(Task).description)
				m.textInput.Focus()
				return m, nil
			case "left", "h":
				m.lists[m.focused].SetShowHelp(false)
				m.Prev()
			case "right", "l":
				m.lists[m.focused].SetShowHelp(false)
				m.Next()
			case "enter":
				msg := m.MoveTask(1)
				if msg != nil {
					return m, msg.(tea.Cmd)
				}
				return m, nil
			case "d":
				if m.lists[m.focused].Index() < 0 {
					return m, nil
				}
				m.RemoveTask()
				return m, nil
			case "backspace":
				msg := m.MoveTask(-1)
				if msg != nil {
					return m, msg.(tea.Cmd)
				}
				return m, nil
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd

	// Handle task adding mode
	if m.addTaskMode {
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	// Handle edit task mode
	if m.editTaskMode {
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	// Update the focused list
	m.lists[m.focused].SetShowHelp(true)
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m *model) View() string {
	if !m.isLoaded {
		return "Loading..."
	}

	// View for adding a new task
	if m.addTaskMode {
		return m.StyledInputTaskView("Add task")
	}

	// View for editing a task
	if m.editTaskMode {
		return m.StyledInputTaskView("Edit task")
	}

	// Get the views for each list
	todoView := m.lists[todo].View()
	inProgressView := m.lists[inProgress].View()
	doneView := m.lists[done].View()

	// Display the lists horizontally, adjusting based on the focused list
	switch m.focused {
	case int(inProgress):
		return StyledListView(1, todoView, inProgressView, doneView)
	case int(done):
		return StyledListView(2, todoView, inProgressView, doneView)
	default:
		return StyledListView(0, todoView, inProgressView, doneView)
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
