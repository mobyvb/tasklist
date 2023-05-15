package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	description string
	finished    bool
	workTime    time.Duration
}

var tasks []Task

func printTasks() {
	fmt.Println("Tasks:")
	for i, task := range tasks {
		status := "Unfinished"
		if task.finished {
			status = "Finished"
		}
		fmt.Printf("%d. %s (%s) - worked for %v\n", i+1, task.description, status, task.workTime)
	}
}

func addTask(description string) {
	tasks = append(tasks, Task{description: description, finished: false})
}

func removeTask(index int) {
	if index < 1 || index > len(tasks) {
		fmt.Println("Invalid task number")
		return
	}

	tasks = append(tasks[:index-1], tasks[index:]...)
}

func updateTask(index int, description string) {
	if index < 1 || index > len(tasks) {
		fmt.Println("Invalid task number")
		return
	}

	tasks[index-1].description = description
}

func finishTask(index int) {
	if index < 1 || index > len(tasks) {
		fmt.Println("Invalid task number")
		return
	}

	tasks[index-1].finished = true
}

func workOnTask(index int) {
	if index < 1 || index > len(tasks) {
		fmt.Println("Invalid task number")
		return
	}

	start := time.Now()

	fmt.Println("Press enter to stop working on the task...")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

	tasks[index-1].workTime += time.Since(start)
}

func handleCommand(input string) {
	tokens := strings.Split(input, " ")
	command := tokens[0]

	switch command {
	case "add":
		addTask(strings.Join(tokens[1:], " "))
	case "remove":
		removeTask(atoi(tokens[1]))
	case "update":
		updateTask(atoi(tokens[1]), strings.Join(tokens[2:], " "))
	case "finish":
		finishTask(atoi(tokens[1]))
	case "work":
		workOnTask(atoi(tokens[1]))
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

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		printTasks()

		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		handleCommand(input)
	}
}
