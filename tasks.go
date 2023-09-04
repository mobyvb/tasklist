package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Task struct {
	Description string        `json:"description"`
	Finished    bool          `json:"finished"`
	WorkTime    time.Duration `json:"work_time"`
}

type TaskList struct {
	Tasks []Task
}

var tasks TaskList

func (t *TaskList) printTaskList() {
	for i, task := range t.Tasks {
		if !task.Finished {
			continue
		}
		t.printTask(task, i+1, "Finished")
	}
	fmt.Println("------------------------------------------")
	for i, task := range t.Tasks {
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
	t.Tasks = append(t.Tasks, Task{Description: description, Finished: false})
}

func (t *TaskList) removeTask(index int) {
	if index < 1 || index > len(t.Tasks) {
		fmt.Println("Invalid task number")
		return
	}

	t.Tasks = append(t.Tasks[:index-1], t.Tasks[index:]...)
}

func (t *TaskList) updateTask(index int, description string) {
	if index < 1 || index > len(t.Tasks) {
		fmt.Println("Invalid task number")
		return
	}

	t.Tasks[index-1].Description = description
}

func (t *TaskList) finishTask(index int) {
	if index < 1 || index > len(t.Tasks) {
		fmt.Println("Invalid task number")
		return
	}

	t.Tasks[index-1].Finished = true
}

func (t *TaskList) workOnTask(index int) {
	if index < 1 || index > len(t.Tasks) {
		fmt.Println("Invalid task number")
		return
	}

	start := time.Now()

	fmt.Println("Press enter to stop working on the task...")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

	t.Tasks[index-1].WorkTime += time.Since(start)
}

func (t *TaskList) saveToFile() {
	taskData, err := json.Marshal(t.Tasks)
	if err != nil {
		fmt.Println("Error saving tasks to file:", err)
		return
	}

	err = ioutil.WriteFile(taskFile, taskData, 0644)
	if err != nil {
		fmt.Println("Error saving tasks to file:", err)
	}
}

func (t *TaskList) loadFromFile() {
	taskData, err := ioutil.ReadFile(taskFile)
	if err != nil {
		fmt.Println("No existing task file found. Starting with an empty task list.")
		return
	}

	err = json.Unmarshal(taskData, &t.Tasks)
	if err != nil {
		fmt.Println("Error loading tasks from file:", err)
	}
}

func (t *TaskList) calculateTotalTime(totalTime time.Duration) time.Duration {
	for _, task := range t.Tasks {
		totalTime += task.WorkTime
	}
	return totalTime
}
