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
	Priority    string        `json:"priority"`
	Deadline    string        `json:"deadline"`
}

var tasks []Task

func printTasks() {
	fmt.Println("Tasks:")

	printTaskList(tasks)

	totalTime := time.Duration(0)
	for _, task := range tasks {
		totalTime += task.WorkTime
	}
	fmt.Printf("\nTotal time spent on all tasks: %v\n", totalTime.Truncate(time.Second))
}

func printTaskList(taskList []Task) {
	fmt.Printf("%-5s %-10s %-5s %-10s %-10v %s\n", "Num", "Status", "Prio", "Deadline", "Worktime", "Desc")
	for i, task := range taskList {
		if !task.Finished {
			continue
		}
		printTask(task, i+1, "Finished")
	}
	fmt.Println("------------------------------------------")
	for i, task := range taskList {
		if task.Finished {
			continue
		}
		printTask(task, i+1, "Unfinished")
	}
}

func printTask(task Task, number int, status string) {
	fmt.Printf("%-5d %-10s %-5s %-10s %-10v %s\n", number, status, task.Priority, task.Deadline, task.WorkTime.Truncate(time.Second), task.Description)
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

func prioritizeTask(index int, priority string) {
	if index < 1 || index > len(tasks) {
		fmt.Println("Invalid task number")
		return
	}

	tasks[index-1].Priority = priority
}

func setDeadline(index int, deadline string) {
	if index < 1 || index > len(tasks) {
		fmt.Println("Invalid task number")
		return
	}

	tasks[index-1].Deadline = deadline
}

func swap(index1, index2 int) {
	if index1 < 1 || index1 > len(tasks) || index2 < 1 || index2 > len(tasks) {
		fmt.Println("Invalid task number")
		return
	}
	if index1 == index2 {
		fmt.Println("identical task numbers - no swap")
	}

	t1 := tasks[index1-1]
	t2 := tasks[index2-1]
	tasks[index1-1] = t2
	tasks[index2-1] = t1
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
	case "priority":
		prio := ""
		if len(tokens) > 2 {
			prio = tokens[2]
		}
		prioritizeTask(atoi(tokens[1]), prio)
	case "deadline":
		dl := ""
		if len(tokens) > 2 {
			dl = tokens[2]
		}
		setDeadline(atoi(tokens[1]), dl)
	case "swap":
		if len(tokens) < 3 {
			fmt.Println("need two args")
			break
		}
		swap(atoi(tokens[1]), atoi(tokens[2]))
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
		saveTasksToFile()
	}

	fmt.Println("Task list saved to file. Goodbye!")
}
