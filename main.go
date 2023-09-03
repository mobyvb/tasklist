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
	SubTasks    []Task        `json:"subtasks"`
}

var tasks []Task

func printTasks() {
	fmt.Println("Tasks:")

	printTaskList(tasks, "")

	totalTime := time.Duration(0)
	for _, task := range tasks {
		totalTime += task.WorkTime
		for _, subTask := range task.SubTasks {
			totalTime += subTask.WorkTime
		}
	}
	fmt.Printf("\nTotal time spent on all tasks: %v\n", totalTime.Truncate(time.Second))
}

func printTaskList(taskList []Task, prefix string) {
	for i, task := range taskList {
		if !task.Finished {
			continue
		}
		printTask(task, prefix+strconv.Itoa(i+1), "Finished")
	}
	fmt.Println("------------------------------------------")
	for i, task := range taskList {
		if task.Finished {
			continue
		}
		printTask(task, prefix+strconv.Itoa(i+1), "Unfinished")
	}
}

func printTask(task Task, number string, status string) {
	fmt.Printf("%-5s %-10s %-10v %s\n", number, status, task.WorkTime.Truncate(time.Second), task.Description)
	printTaskList(task.SubTasks, number+".")
}

func addTask(description string) {
	tasks = append(tasks, Task{Description: description, Finished: false})
}

func addSubTask(parentIndex int, subDescription string) {
	if parentIndex < 1 || parentIndex > len(tasks) {
		fmt.Println("Invalid task number")
		return
	}

	tasks[parentIndex-1].SubTasks = append(tasks[parentIndex-1].SubTasks, Task{Description: subDescription, Finished: false})
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

func workOnTask(index int, subIndex int) {
	if index < 1 || index > len(tasks) {
		fmt.Println("Invalid task number")
		return
	}

	var taskToWork *Task
	if subIndex > 0 {
		if subIndex > len(tasks[index-1].SubTasks) {
			fmt.Println("Invalid subtask number")
			return
		}

		taskToWork = &tasks[index-1].SubTasks[subIndex-1]
	} else {
		taskToWork = &tasks[index-1]
	}

	start := time.Now()

	fmt.Println("Press enter to stop working on the task...")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

	taskToWork.WorkTime += time.Since(start)
}

func handleCommand(input string) {
	tokens := strings.Split(input, " ")
	command := tokens[0]

	switch command {
	case "add":
		addTask(strings.Join(tokens[1:], " "))
	case "add-subtask":
		parentTaskNum, _ := strconv.Atoi(tokens[1])
		addSubTask(parentTaskNum, strings.Join(tokens[2:], " "))
	case "remove":
		removeTask(atoi(tokens[1]))
	case "update":
		updateTask(atoi(tokens[1]), strings.Join(tokens[2:], " "))
	case "finish":
		finishTask(atoi(tokens[1]))
	case "work":
		parts := strings.Split(tokens[1], ".")
		index := atoi(parts[0])
		subIndex := 0
		if len(parts) > 1 {
			subIndex = atoi(parts[1])
		}
		workOnTask(index, subIndex)
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

	err = ioutil.WriteFile(taskFile, taskData, 0644)
	if err != nil {
		fmt.Println("Error saving tasks to file:", err)
	}
}

func loadTasksFromFile() {
	taskData, err := ioutil.ReadFile(taskFile)
	if err != nil {
		fmt.Println("No existing task file found. Starting with an empty task list.")
		return
	}

	err = json.Unmarshal(taskData, &tasks)
	if err != nil {
		fmt.Println("Error loading tasks from file:", err)
	}
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
