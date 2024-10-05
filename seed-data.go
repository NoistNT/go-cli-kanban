package main

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
