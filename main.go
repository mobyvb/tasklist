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

type TaskList struct {
	tasks    []Task
	taskFile string
}

func NewTaskList(taskFile string) *TaskList {
	return &TaskList{taskFile: taskFile}
}

func (t *TaskList) printTasks() {
	fmt.Println("Tasks:")

	t.printTaskList()

	totalTime := time.Duration(0)
	for _, task := range t.tasks {
		totalTime += task.WorkTime
	}
	fmt.Printf("\nTotal time spent on all tasks: %v\n", totalTime.Truncate(time.Second))
}

func (t *TaskList) printTaskList() {
	for i, task := range t.tasks {
		if !task.Finished {
			continue
		}
		t.printTask(task, i+1, "Finished")
	}
	fmt.Println("------------------------------------------")
	for i, task := range t.tasks {
		if task.Finished {
			continue
		}
		t.printTask(task, i+1, "Unfinished")
	}
}

func (t *TaskList) printTask(task Task, number int, status string) {
	fmt.Printf("%-5d %-10s %-10v %s\n", number, status, task.WorkTime.Truncate(time.Second), task.Description)
}

func (t *TaskList) addTask(description string) {
	t.tasks = append(t.tasks, Task{Description: description, Finished: false})
}

func (t *TaskList) removeTask(index int) {
	if index < 1 || index > len(t.tasks) {
		fmt.Println("Invalid task number")
		return
	}

	t.tasks = append(t.tasks[:index-1], t.tasks[index:]...)
}

func (t *TaskList) updateTask(index int, description string) {
	if index < 1 || index > len(t.tasks) {
		fmt.Println("Invalid task number")
		return
	}

	t.tasks[index-1].Description = description
}

func (t *TaskList) finishTask(index int) {
	if index < 1 || index > len(t.tasks) {
		fmt.Println("Invalid task number")
		return
	}

	t.tasks[index-1].Finished = true
}

func (t *TaskList) workOnTask(index int) {
	if index < 1 || index > len(t.tasks) {
		fmt.Println("Invalid task number")
		return
	}

	start := time.Now()

	fmt.Println("Press enter to stop working on the task...")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

	t.tasks[index-1].WorkTime += time.Since(start)
}

func handleCommand(input string, taskList *TaskList) {
	tokens := strings.Split(input, " ")
	command := tokens[0]

	switch command {
	case "add":
		taskList.addTask(strings.Join(tokens[1:], " "))
	case "remove":
		taskList.removeTask(atoi(tokens[1]))
	case "update":
		taskList.updateTask(atoi(tokens[1]), strings.Join(tokens[2:], " "))
	case "finish":
		taskList.finishTask(atoi(tokens[1]))
	case "work":
		taskList.workOnTask(atoi(tokens[1]))
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

func (t *TaskList) saveTasksToFile() {
	taskData, err := json.Marshal(t.tasks)
	if err != nil {
		fmt.Println("Error saving tasks to file:", err)
		return
	}

	err = ioutil.WriteFile(t.taskFile, taskData, 0644)
	if err != nil {
		fmt.Println("Error saving tasks to file:", err)
	}
}

func (t *TaskList) loadTasksFromFile() {
	taskData, err := ioutil.ReadFile(t.taskFile)
	if err != nil {
		fmt.Println("No existing task file found. Starting with an empty task list.")
		return
	}

	err = json.Unmarshal(taskData, &t.tasks)
	if err != nil {
		fmt.Println("Error loading tasks from file:", err)
	}
}

func main() {
	var taskFile string

	if len(os.Args) > 1 {
		taskFile = os.Args[1]
	} else {
		taskFile = "tasks.dat"
	}

	taskList := NewTaskList(taskFile)
	taskList.loadTasksFromFile()

	reader := bufio.NewReader(os.Stdin)

	for {
		taskList.printTasks()

		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "quit" {
			break
		}

		handleCommand(input, taskList)
	}

	taskList.saveTasksToFile()

	fmt.Println("Task list saved to file. Goodbye!")
}
