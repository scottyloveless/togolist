package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// "os"
// "fmt"
// "strconv"

type Todo struct {
	ID   int
	Task string
	Done bool
}

var todos []Todo
var nextID = 1

func main() {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("\nToGo Task Manager")
		fmt.Print("Choose an option: ")
		scanner.Scan()
		choice := strings.ToLower(scanner.Text())

		switch choice {
		case "a":
			fmt.Print("Enter task: ")
			scanner.Scan()
			task := scanner.Text()
			if task != "" {
				todos = append(todos, Todo{ID: nextID, Task: task, Done: false})
				nextID++
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
		fmt.Printf("[%d] %s [%s]\n", todo.ID, todo.Task, status)
	}
}

func completeTodo(id int) {
	for _, todo := range todos {
		if todo.ID == id {
			todo.Done = true
			fmt.Println("Task completed!")
			return
		}
	}
	fmt.Println("Task not found!")
}

func deleteTodo(id int) {
	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			fmt.Println("Task deleted!")
			return
		}
	}
	fmt.Println("Task not found!")
}
