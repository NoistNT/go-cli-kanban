package main

import tea "github.com/charmbracelet/bubbletea"

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

// AddTask adds the task to the list
func (m *model) AddTask() tea.Msg {
	newTask := Task{
		title:       m.textInput.Value(),
		description: m.textInput.Value(),
		status:      todo,
	}

	todoList := &m.lists[todo]
	todoList.InsertItem(len(todoList.Items()), newTask)
	m.addTaskMode = false // Return to list view after adding
	m.textInput.Blur()    // Remove focus from input

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

	m.editTaskMode = false // Return to list view after editing
	m.textInput.Blur()     // Remove focus from input

	return tea.Batch(tea.EnterAltScreen, tea.ClearScreen)
}

// RemoveTask removes the selected task from the list
func (m *model) RemoveTask() tea.Msg {
	// Get the index of the selected task
	selectedIndex := m.lists[m.focused].Index()

	// Remove the task from the list
	m.lists[m.focused].RemoveItem(selectedIndex)

	return tea.Batch(tea.EnterAltScreen, tea.ClearScreen)
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

// MoveTask moves the selected task to the next list (+1) or the previous list (-1)
func (m *model) MoveTask(direction int) tea.Msg {
	// Get the current focused list and its tasks
	fl := m.lists[m.focused]

	// Check if the list is empty
	if len(fl.Items()) == 0 {
		return nil
	}

	// Get the selected task
	selectedIndex := fl.Index()
	task, ok := fl.SelectedItem().(Task)
	if !ok || (direction == 1 && task.status == done) || (direction == -1 && task.status == todo) {
		return nil
	}

	// Remove the selected task from the current list
	fl.RemoveItem(selectedIndex)

	// Move the selected task to the correct list
	if direction == 1 {
		task.Next()
	} else {
		task.Prev()
	}

	// Add the task to the next list
	nextList := m.lists[task.status]
	nextList.InsertItem(len(nextList.Items()), task)

	// Update Bubble Tea's internal list state after removal and insertion
	m.lists[m.focused] = fl
	m.lists[task.status] = nextList

	// Return a command to trigger a re-render and update the UI
	return tea.Batch(tea.EnterAltScreen, tea.ClearScreen)
}
