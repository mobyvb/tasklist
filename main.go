package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Task struct {
	Description string
	Duration    time.Duration
}

func main() {
	taskListFile := "taskList.json"
	tasks := loadTasks(taskListFile)

	defer saveTasks(taskListFile, tasks)

	for {
		fmt.Println("Please enter a command (list, add, work): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('
')
		input = strings.TrimSuffix(input, "
")

		switch input {
		case "list":
			printTasks(tasks)
		case "add":
			tasks = addTask(tasks)
		case "work":
			timer(tasks)
		default:
			fmt.Println("Unknown command.")
		}
	}
}

func printTasks(tasks []Task) {
	for i, task := range tasks {
		fmt.Printf("%d. Description: %s, Duration: %v (sec)
", i+1, task.Description, task.Duration.Truncate(time.Second))
	}
}

func addTask(tasks []Task) []Task {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter a task description: ")
	input, _ := reader.ReadString('
')
	desc := strings.TrimSuffix(input, "
")

	tasks = append(tasks, Task{desc, 0})
	return tasks
}

func timer(tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println("No tasks found. Please add one first.")
		return
	}

	printTasks(tasks)

	fmt.Println("Enter the number of the task to work on: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('
')

	var taskIndex int
	_, err := fmt.Sscanf(input, "%d", &taskIndex)

	if err != nil || taskIndex < 1 || taskIndex > len(tasks) {
		fmt.Println("Invalid task number.")
		return
	}

	workTask(&tasks[taskIndex-1])

	// Print updated list of tasks.
	fmt.Println("Updated task list:")
	printTasks(tasks)
}

func workTask(task *Task) {
	fmt.Println("Starting task: ", task.Description)

	startTime := time.Now()

	fmt.Println("Do you want to set a predefined countdown time? (yes/no): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('
')
	input = strings.TrimSuffix(input, "
")

	if input == "yes" {
		fmt.Println("Set countdown time in the format (mm): ")
		reader := bufio.NewReader(os.Stdin)
		counter, _ := reader.ReadString('
')
		counter = strings.TrimSuffix(counter, "
")

		var countdown int
		_, err := fmt.Sscanf(counter, "%d", &countdown)
		if err != nil {
			fmt.Println("Invalid countdown time.")
			return
		}

		fmt.Printf("Countdown timer set for %d minutes.
", countdown)
		time.Sleep(time.Duration(countdown) * time.Minute)
		fmt.Println("Countdown complete. Moving to next task.")

	} else {
		fmt.Println("Press enter to stop.")
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('
')

		timeElapsed := time.Since(startTime)
		task.Duration += timeElapsed
	}
}

func loadTasks(taskListFile string) []Task {
	tasks := []Task{}

	raw, err := ioutil.ReadFile(taskListFile)
	if err != nil {
		return tasks
	}

	json.Unmarshal(raw, &tasks)
	return tasks
}

func saveTasks(taskListFile string, tasks []Task) {
	raw, _ := json.Marshal(tasks)
	ioutil.WriteFile(taskListFile, raw, 0644)
}