package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	Description string        `json:"description"`
	Finished    bool          `json:"finished"`
	WorkTime    time.Duration `json:"work_time"`
}

var tasks []Task

func printTasks() {
	fmt.Println("Tasks:")
	for i, task := range tasks {
		status := "Unfinished"
		if task.Finished {
			status = "Finished"
		}
		fmt.Printf("%d. %s (%s) - worked for %v\n", i+1, task.Description, status, task.WorkTime)
	}
}

func addTask(description string) {
	tasks = append(tasks, Task{Description: description, Finished: false})
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

	tasks[index-1].Description = description
}

func finishTask(index int) {
	if index < 1 || index > len(tasks) {
		fmt.Println("Invalid task number")
		return
	}

	tasks[index-1].Finished = true
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

	tasks[index-1].WorkTime += time.Since(start)
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

func saveTasksToFile() {
	taskData, err := json.Marshal(tasks)
	if err != nil {
		fmt.Println("Error saving tasks to file:", err)
		return
	}

	err = ioutil.WriteFile("tasks.dat", taskData, 0644)
	if err != nil {
		fmt.Println("Error saving tasks to file:", err)
	}
}

func loadTasksFromFile() {
	taskData, err := ioutil.ReadFile("tasks.dat")
	if err != nil {
		fmt.Println("No existing task file found. Starting with an empty task list.")
		return
	}

	err = json.Unmarshal(taskData, &tasks)
	if err != nil {
		fmt.Println("Error loading tasks from file:", err)
	}
}

func main() {
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
