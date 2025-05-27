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

const storageFile = "storage.json"

func main() {
	loadTodos()
	scanner := bufio.NewScanner(os.Stdin)
	for {

		fmt.Println("\nToGo Task Manager")
		fmt.Println("[a] add | [l] list | [d] delete | [c] complete | [x] clear completed")
		fmt.Print("Choose an option: ")
		scanner.Scan()
		choice := strings.ToLower(scanner.Text())

		switch choice {
		case "a":
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
		case "l":
			listTodos()
		case "d":
			fmt.Print("Enter task ID to delete: ")
			scanner.Scan()
			var id int
			fmt.Sscanf(scanner.Text(), "%d", &id)
			deleteTodo(id)
			saveTodos()
		case "c":
			fmt.Print("Enter task ID to complete: ")
			scanner.Scan()
			var id int
			fmt.Sscanf(scanner.Text(), "%d", &id)
			completeTodo(id)
			saveTodos()
		case "q":
			fmt.Println("Goodbye!")
			return
		case "x":
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
		fmt.Printf("[%d] %s [%s] !%d\n", todo.ID, todo.Task, status, todo.Priority)
	}
}

func completeTodo(id int) {
	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Done = true
			fmt.Println("Task completed!")
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
