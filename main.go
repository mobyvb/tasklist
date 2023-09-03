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
	Subtasks    []Task        `json:"subtasks"`
}

var tasks []Task

func printTasks() {
	fmt.Println("Tasks:")

	printTaskList(tasks, "")

	totalTime := time.Duration(0)
	for _, task := range tasks {
		totalTime += task.WorkTime
		for _, subtask := range task.Subtasks {
			totalTime += subtask.WorkTime
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
	printTaskList(task.Subtasks, number+".")
}

func addTask(description string) {
	tasks = append(tasks, Task{Description: description, Finished: false})
}

func addSubtask(index int, description string) {
	if index < 1 || index > len(tasks) {
		fmt.Println("Invalid task number")
		return
	}

	tasks[index-1].Subtasks = append(tasks[index-1].Subtasks, Task{Description: description, Finished: false})
}

func removeTask(index int) {
	if index < 1 || index > len(tasks) {
		fmt.Println("Invalid task number")
		return
	}

	tasks = append(tasks[:index-1], tasks[index:]...)
}

// ... (more code follows)
