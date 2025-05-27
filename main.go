package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Todo struct {
	ID       int    `json:"id"`
	Task     string `json:"task"`
	Done     bool   `json:"done"`
	Priority int    `json:"priority"`
}

var todos []Todo
var nextID = 1
var sortType = false

const storageFile = "storage.json"

func main() {
	loadTodos()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		clearTerminal()
		topMenu()
		fmt.Print("Choose an option: ")
		scanner.Scan()
		choice := strings.ToLower(scanner.Text())

		switch choice {
		case "a":
			clearTerminal()
			fmt.Print("Enter task: ")
			scanner.Scan()
			task := scanner.Text()
			fmt.Print("Enter priority 1-5: ")
			scanner.Scan()
			var prio int
			if _, err := fmt.Sscanf(scanner.Text(), "%d", &prio); err != nil {
				fmt.Print("Invalid entry: please enter 1-5 for priority")
			}
			if prio < 1 {
				prio = 1
				fmt.Print("Priority rounded up to 1\n")
			}
			if prio > 5 {
				prio = 5
				fmt.Print("Priority rounded down to 5\n")
			}

			if task != "" {
				todos = append(todos, Todo{ID: nextID, Task: task, Done: false, Priority: prio})
				nextID++
				saveTodos()
				fmt.Println("Task added!")
			}
		// case "l":
		// 	clearTerminal()
		// 	listTodos()
		case "d":
			clearTerminal()
			topMenu()
			fmt.Print("Enter task ID to delete: ")
			scanner.Scan()
			var id int
			fmt.Sscanf(scanner.Text(), "%d", &id)
			deleteTodo(id)
			saveTodos()
		case "c":
			clearTerminal()
			topMenu()
			fmt.Print("Enter task ID to complete: ")
			scanner.Scan()
			var id int
			fmt.Sscanf(scanner.Text(), "%d", &id)
			completeTodo(id)
			saveTodos()
		case "q":
			fmt.Println("Goodbye!")
			return
		case "s":
			clearTerminal()
			topMenu()
			fmt.Println("Sort by [i]d | [p]riority: ")
			scanner.Scan()
			sortChoice := strings.ToLower(scanner.Text())
			switch sortChoice {
			case "i":
				sortTodosById()
				saveTodos()
			case "p":
				sortTodosByPriority()
				saveTodos()
			default:
				fmt.Println("Invalid option")

			}

		case "x":
			clearTerminal()
			clearCompleted()
			saveTodos()
			fmt.Println("Completed todos deleted")
		default:
			fmt.Println("Invalid option")
		}
	}
}

func listTodos() {
	if len(todos) == 0 {
		fmt.Println("No tasks!")
		return
	}

	for _, todo := range todos {
		status := " "
		if todo.Done {
			status = "✓"
		}
		fmt.Printf("[%d]  %s    !%d   [%s]\n", todo.ID, todo.Task, todo.Priority, status)
	}

}

func sortTodosByPriority() {
	n := len(todos)
	for i := 0; i < n+1; i++ {
		for j := 0; j < n-i-1; j++ {
			if todos[j].Priority < todos[j+1].Priority {
				todos[j], todos[j+1] = todos[j+1], todos[j]
			}
		}
	}
}

func sortTodosById() {
	n := len(todos)
	for i := 0; i < n+1; i++ {
		for j := 0; j < n-i-1; j++ {
			if todos[j].ID > todos[j+1].ID {
				todos[j], todos[j+1] = todos[j+1], todos[j]
			}
		}
	}
}

func completeTodo(id int) {
	for i, todo := range todos {
		if todo.ID == id {
			if todo.Done == true {
				todos[i].Done = false
				return
			}
			todos[i].Done = true
			return
		}
	}
	fmt.Println("Task not found!")
}

func deleteTodo(id int) {
	for i, todo := range todos {
		if todo.ID == id {
			todos = slices.Delete(todos, i, i+1)
			fmt.Println("Task deleted!")
			return
		}
	}
	fmt.Println("Task not found!")
}

func loadTodos() {
	file, err := os.ReadFile(storageFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Printf("Error loading todos: %v\n", err)
		return
	}
	if err := json.Unmarshal(file, &todos); err != nil {
		fmt.Printf("Error parsing todos: %v\n", err)
		return
	}
	for _, todo := range todos {
		if todo.ID >= nextID {
			nextID = todo.ID + 1
		}
	}
}

func saveTodos() {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		fmt.Printf("Error saving todos: %v\n", err)
		return
	}
	if err := os.WriteFile(storageFile, data, 0644); err != nil {
		fmt.Printf("Error writing todos: %v\n", err)
	}
}

func clearCompleted() {
	var completed []int
	for _, todo := range todos {
		if todo.Done == true {
			completed = append(completed, todo.ID)
		}
	}
	slices.Reverse(completed)

	if len(completed) == 0 {
		fmt.Println("No completed todos to clear!")
		return
	}

	for _, v := range completed {
		deleteTodo(v)
	}
}

func topMenu() {
	fmt.Println("Todos:")
	fmt.Println(" ")
	listTodos()
	fmt.Println(" ")
	fmt.Println("[a]dd | [d]elete | [c]omplete | [s]ort | [x]cleanup | [q]uit")
}
