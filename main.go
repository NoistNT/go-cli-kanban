package main

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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

// todoListItems seed data for To Do
var todoListItems = listItems{
	Task{
		status:      todo,
		title:       "Complete Weekly Report",
		description: "Finish the weekly report by the end of the day.",
	},
	Task{
		status:      todo,
		title:       "Grocery Shopping",
		description: "Buy groceries for the week.",
	},
	Task{
		status:      todo,
		title:       "Attend Team Meeting",
		description: "Participate in the weekly team meeting.",
	},
	Task{
		status:      todo,
		title:       "Learn a New Skill",
		description: "Dedicate 30 minutes to learning a new skill.",
	},
	Task{
		status:      todo,
		title:       "Schedule Doctor's Appointment",
		description: "Book a doctor's appointment for a check-up.",
	},
}

// inProgressListItems seed data for In Progress
var inProgressListItems = listItems{
	Task{
		status:      inProgress,
		title:       "Research New Software",
		description: "Compare different software options for the project.",
	},
	Task{
		status:      inProgress,
		title:       "Write Blog Post",
		description: "Draft a blog post about [topic].",
	},
	Task{
		status:      inProgress,
		title:       "Exercise",
		description: "Do a 30-minute workout.",
	},
	Task{
		status:      inProgress,
		title:       "Prepare Presentation",
		description: "Create slides for the upcoming presentation.",
	},
}

// doneListItems seed data for Done
var doneListItems = listItems{
	Task{
		status:      done,
		title:       "Complete Project Proposal",
		description: "Submit the project proposal to the client.",
	},
	Task{
		status:      done,
		title:       "Pay Bills",
		description: "Pay all outstanding bills.",
	},
	Task{
		status:      done,
		title:       "Read Book",
		description: "Finish reading the book.",
	},
	Task{
		status:      done,
		title:       "Clean House",
		description: "Thoroughly clean the house.",
	},
}

type status int

const (
	todo status = iota
	inProgress
	done
)

// Task struct
type Task struct {
	status      status
	title       string
	description string
}

// Next method to update the task status to the next state (To Do -> In Progress -> Done)
func (t *Task) Next() {
	if t.status < done {
		t.status++
	}
}

// Prev method to update the task status to the previous state (Done -> In Progress -> To Do)
func (t *Task) Prev() {
	if t.status > todo {
		t.status--
	}
}

// FilterValue returns the title of the task
func (t Task) FilterValue() string {
	return t.title
}

// Title returns the title of the task
func (t Task) Title() string {
	return t.title
}

// Description returns the description of the task
func (t Task) Description() string {
	return t.description
}

type listItems []list.Item

type lists []list.Model

// MAIN MODEL
type model struct {
	lists      lists
	focused    int
	err        error
	isLoaded   bool
	addingTask bool
	editMode   bool
	textInput  textinput.Model
}

// ShowHelp shows the help screen when the user is in task-adding/editing mode
func (m *model) ShowHelp() string {
	// Helper message for when the user is in task-adding mode
	return helpStyle.Render("Press 'enter' to confirm, 'esc' to cancel, 'ctrl+c' to quit")
}

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

// AddTask adds the task to the list
func (m *model) AddTask() tea.Msg {
	newTask := Task{
		title:       m.textInput.Value(),
		description: m.textInput.Value(),
		status:      todo,
	}

	todoList := &m.lists[todo]
	todoList.InsertItem(len(todoList.Items()), newTask)
	m.addingTask = false // Return to list view after adding
	m.textInput.Blur()   // Remove focus from input

	return tea.Batch(tea.EnterAltScreen, tea.ClearScreen)
}

// EditTask edits the selected task in the list
func (m *model) EditTask() tea.Msg {
	// Get the index of the selected task
	selectedIndex := m.lists[m.focused].Index()

	// Update the task in the list
	selectedTask := m.lists[m.focused].Items()[selectedIndex].(Task)
	selectedTask.title = m.textInput.Value()
	selectedTask.description = m.textInput.Value()

	// Update the list with the modified items
	m.lists[m.focused].RemoveItem(selectedIndex)
	m.lists[m.focused].InsertItem(selectedIndex, selectedTask)

	m.editMode = false // Return to list view after editing
	m.textInput.Blur() // Remove focus from input

	return tea.Batch(tea.EnterAltScreen, tea.ClearScreen)
}

// MoveToNext moves the selected task to the next list
func (m *model) MoveToNext() tea.Msg {
	// Get the current focused list and its tasks
	fl := m.lists[m.focused]

	// Check if the list is empty
	if len(fl.Items()) == 0 {
		return nil
	}

	// Get the selected task
	selectedIndex := fl.Index()
	selected := fl.SelectedItem()
	task, ok := selected.(Task)
	if !ok {
		return nil
	}

	// If the task is already in "Done" state, prevent further movement
	if task.status == done {
		return nil
	}

	// Remove the selected task from the current list
	fl.RemoveItem(selectedIndex)

	// Move the selected task to the next list and update its status
	task.Next()

	// Add the task to the next list
	nextList := m.lists[task.status]
	nextList.InsertItem(len(nextList.Items()), task)

	// Update Bubble Tea's internal list state after removal and insertion
	m.lists[m.focused] = fl
	m.lists[task.status] = nextList

	// Return a command to trigger a re-render and update the UI
	return tea.Batch(tea.EnterAltScreen, tea.ClearScreen)
}

// MoveToPrev moves the selected task to the previous list
func (m *model) MoveToPrev() tea.Msg {
	// Get the current focused list and its tasks
	fl := m.lists[m.focused]

	// Check if the list is empty
	if len(fl.Items()) == 0 {
		return nil
	}

	// Get the selected task
	selectedIndex := fl.Index()
	selected := fl.SelectedItem()
	task, ok := selected.(Task)
	if !ok {
		return nil
	}

	// If the task is already in "To Do" state, prevent further movement
	if task.status == todo {
		return nil
	}

	// Remove the selected task from the current list
	fl.RemoveItem(selectedIndex)

	// Move the selected task to the previous list and update its status
	task.Prev()

	// Add the task to the previous list
	prevList := m.lists[task.status]
	prevList.InsertItem(len(prevList.Items()), task)

	// Reset index in the original list to prevent leftover selection
	if selectedIndex >= len(fl.Items()) {
		selectedIndex = len(fl.Items()) - 1
	}
	fl.Select(selectedIndex)

	// Update Bubble Tea's internal list state after removal and insertion
	m.lists[m.focused] = fl
	m.lists[task.status] = prevList

	// Return a command to trigger a re-render and update the UI
	return tea.Batch(tea.EnterAltScreen, tea.ClearScreen)
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

// Next method to update the task status to the next state (To Do -> In Progress -> Done)
func (m *model) Next() {
	m.focused = (m.focused + 1) % len(m.lists)
}

// Prev method to update the task status to the previous state (Done -> In Progress -> To Do)
func (m *model) Prev() {
	m.focused = (m.focused - 1) % len(m.lists)
	if m.focused < 0 {
		m.focused = len(m.lists) - 1
	}
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
		case m.addingTask:
			// In task-adding mode, forward the key events to the text input
			switch msg.String() {
			case "enter":
				// Add the task and return to list view
				m.AddTask()
				return m, nil
			case "esc":
				// Cancel adding task and return to list view
				m.addingTask = false
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			default:
				// Forward all other key inputs (including 'a') to the text input
				var cmd tea.Cmd
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd
			}

		case m.editMode:
			// In task-editing mode, forward the key events to the text input
			switch msg.String() {
			case "enter":
				// Add the task and return to list view
				m.EditTask()
				return m, nil
			case "esc":
				// Cancel editing task and return to list view
				m.addingTask = false
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
				m.addingTask = true
				m.textInput.Focus()
				return m, nil
			case "e":
				// Toggle to edit mode
				m.editMode = true
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
				return m, m.MoveToNext
			case "backspace":
				return m, m.MoveToPrev
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd

	// Handle task adding mode
	if m.addingTask {
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	// Handle edit task mode
	if m.editMode {
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
	if m.addingTask {
		return lipgloss.NewStyle().Padding(1).Render("Add Task: \n"+m.textInput.View()) + m.ShowHelp()
	}

	// View for editing a task
	if m.editMode {
		return lipgloss.NewStyle().Padding(1).Render("Edit Task: \n"+m.textInput.View()) + m.ShowHelp()
	}

	// Get the views for each list
	todoView := m.lists[todo].View()
	inProgressView := m.lists[inProgress].View()
	doneView := m.lists[done].View()

	// Display the lists horizontally, adjusting based on the focused list
	switch m.focused {
	case int(inProgress):
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			columnStyle.Render(todoView),
			focusedStyle.Render(inProgressView),
			columnStyle.Render(doneView),
		)
	case int(done):
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			columnStyle.Render(todoView),
			columnStyle.Render(inProgressView),
			focusedStyle.Render(doneView),
		)
	default:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			focusedStyle.Render(todoView),
			columnStyle.Render(inProgressView),
			columnStyle.Render(doneView),
		)
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
