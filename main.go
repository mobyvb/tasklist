package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func printTasks() {
	fmt.Println("Tasks:")

	tasks.printTaskList()

	totalTime := time.Duration(0)
	totalTime = tasks.calculateTotalTime(totalTime)
	fmt.Printf("\nTotal time spent on all tasks: %v\n", totalTime.Truncate(time.Second))
}

func handleCommand(input string) {
	tokens := strings.Split(input, " ")
	command := tokens[0]

	switch command {
	case "add":
		tasks.addTask(strings.Join(tokens[1:], " "))
	case "remove":
		tasks.removeTask(atoi(tokens[1]))
	case "update":
		tasks.updateTask(atoi(tokens[1]), strings.Join(tokens[2:], " "))
	case "finish":
		tasks.finishTask(atoi(tokens[1]))
	case "work":
		tasks.workOnTask(atoi(tokens[1]))
	default:
		fmt.Println("Invalid command")
	}
}

func atoi(str string) int {
	result, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return result
}

func saveTasksToFile() {
	tasks.saveToFile()
}

func loadTasksFromFile() {
	tasks.loadFromFile()
}

var taskFile string

func main() {
	if len(os.Args) > 1 {
		taskFile = os.Args[1]
	} else {
		taskFile = "tasks.dat"
	}

	loadTasksFromFile()

	reader := bufio.NewReader(os.Stdin)

	for {
		printTasks()

		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "quit" {
			break
		}

		handleCommand(input)
	}

	saveTasksToFile()

	fmt.Println("Task list saved to file. Goodbye!")
}
