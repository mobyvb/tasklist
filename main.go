```
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

func (tl *TaskList) printTasks() {
    fmt.Println("Tasks:")

    tl.printTaskList()

    totalTime := time.Duration(0)
    for _, task := range tl.Tasks {
        totalTime += task.WorkTime
    }
    fmt.Printf("\nTotal time spent on all tasks: %v\n", totalTime.Truncate(time.Second))
}

func (tl *TaskList) printTaskList() {
    for i, task := range tl.Tasks {
        if !task.Finished {
            continue
        }
        tl.printTask(task, i+1, "Finished")
    }
    fmt.Println("------------------------------------------")
    for i, task := range tl.Tasks {
        if task.Finished {
            continue
        }
        tl.printTask(task, i+1, "Unfinished")
    }
}

func (tl *TaskList) printTask(task Task, number int, status string) {
    fmt.Printf("%-5d %-10s %-10v %s\n", number, status, task.WorkTime.Truncate(time.Second), task.Description)
}

func (tl *TaskList) addTask(description string) {
    tl.Tasks = append(tl.Tasks, Task{Description: description, Finished: false})
}

func (tl *TaskList) removeTask(index int) {
    if index < 1 || index > len(tl.Tasks) {
        fmt.Println("Invalid task number")
        return
    }

    tl.Tasks = append(tl.Tasks[:index-1], tl.Tasks[index:]...)
}

func (tl *TaskList) updateTask(index int, description string) {
    if index < 1 || index > len(tl.Tasks) {
        fmt.Println("Invalid task number")
        return
    }

    tl.Tasks[index-1].Description = description
}

func (tl *TaskList) finishTask(index int) {
    if index < 1 || index > len(tl.Tasks) {
        fmt.Println("Invalid task number")
        return
    }

    tl.Tasks[index-1].Finished = true
}

func (tl *TaskList) workOnTask(index int) {
    if index < 1 || index > len(tl.Tasks) {
        fmt.Println("Invalid task number")
        return
    }

    start := time.Now()

    fmt.Println("Press enter to stop working on the task...")
    reader := bufio.NewReader(os.Stdin)
    _, _ = reader.ReadString('\n')

    tl.Tasks[index-1].WorkTime += time.Since(start)
}

func atoi(str string) int {
    result, err := strconv.Atoi(str)
    if err != nil {
        return 0
    }
    return result
}

func main() {
    taskFile := "tasks.dat"
    if len(os.Args) > 1 {
        taskFile = os.Args[1]
    }
    
    taskList := &TaskList{}
    taskList.loadTasksFromFile(taskFile)

    reader := bufio.NewReader(os.Stdin)

    for {
        taskList.printTasks()

        fmt.Print("> ")
        input, _ := reader.ReadString('\n')
        input = strings.TrimSpace(input)

        if input == "quit" {
            break
        }

        taskList.handleCommand(input)
    }

    taskList.saveTasksToFile(taskFile)

    fmt.Println("Task list saved to file. Goodbye!")
}
```
