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
	Tasks []Task
}

func newTaskList() *TaskList {
	return &TaskList{}
}

func (tl *TaskList) addTask(description string) {
	tl.Tasks = append(tl.Tasks, Task{Description: description, Finished: false})
}

func (tl *TaskList) printTasks() {
	fmt.Println("Tasks:")

	tl.printTaskList()

	totalTime := time.Duration(0)
	for _, task := range tl.Tasks {
		totalTime += task.WorkTime
	}

	fmt.Printf("\nTotal time spent on all tasks: %v\n", totalTime.Truncate(time.Second))

}

// {...Rest of your code}

func main() {
	tl := newTaskList()

	if len(os.Args) > 1 {
		taskFile = os.Args[1]
	} else {
		taskFile = "tasks.dat"
	}

	loadTasksFromFile(tl)

	reader := bufio.NewReader(os.Stdin)

	for {
		tl.printTasks()

		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "quit" {
			break
		}

		handleCommand(input, tl)
	}

	saveTasksToFile(tl)

	fmt.Println("Task list saved to file. Goodbye!")
}
