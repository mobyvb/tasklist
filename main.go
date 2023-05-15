package main

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "os"
    "strconv"
    "strings"
    "time"
)

const (
    taskListFile = "task_list.txt"
)

func main() {
    // Load saved task list, if it exists
    savedTaskList, err := ioutil.ReadFile(taskListFile)
    if err == nil {
        savedTasks := strings.Split(string(savedTaskList), "
")
        fmt.Println("Loaded saved task list:")
        for _, task := range savedTasks {
            fmt.Println(task)
        }
    }

    reader := bufio.NewReader(os.Stdin)
    startTime := time.Now()
    var totalTime time.Duration

    for {
        // Prompt user to enter command
        fmt.Print("> ")
        command, _ := reader.ReadString('
')
        command = strings.TrimSpace(command)

        // Process commands
        switch strings.ToLower(command) {

        case "exit":
            // Save task list to file
            ioutil.WriteFile(taskListFile, []byte(strings.Join(taskList, "
")), 0644)
            os.Exit(0)

        case "work":
            // Switch to work mode
            workTime, err := work(reader, 0)
            if err != nil {
                fmt.Println("Invalid countdown time")
            } else {
                totalTime += workTime
            }

        default:
            if strings.HasPrefix(command, "work ") {
                countdownStr := strings.TrimSpace(strings.TrimPrefix(command, "work "))
                countdown, err := strconv.Atoi(countdownStr)
                if err != nil || countdown <= 0 {
                    fmt.Println("Invalid countdown time")
                } else {
                    // If countdown time provided, start countdown
                    workTime, err := work(reader, countdown)
                    if err != nil {
                        fmt.Println("Invalid countdown time")
                    } else {
                        totalTime += workTime
                    }
                }
            } else {
                fmt.Println("invalid command")
            }
        }
    }
}

// work function starts a timer and counts up or down depending on countdown parameter.
// If countdown is 0, timer counts up. If countdown is positive, counts down.
func work(reader *bufio.Reader, countdown int) (time.Duration, error) {
    startTime := time.Now()
    timer := time.NewTimer(1 * time.Second)

    if countdown == 0 {
        // Count up mode
        go func() {
            for range timer.C {
                fmt.Printf("\r%d seconds", time.Now().Sub(startTime)/time.Second)
                timer.Reset(1 * time.Second)
            }
        }()
    } else {
        // Countdown mode
        go func() {
            for range timer.C {
                remaining := time.Duration(countdown)*time.Second - time.Now().Sub(startTime)
                if remaining <= 0 {
                    break
                }
                fmt.Printf("\r%d seconds", remaining/time.Second)
                timer.Reset(1 * time.Second)
            }
        }()
    }

    fmt.Print("Type 'stop' to stop: ")
    for {
        line, _ := reader.ReadString('
')
        line = strings.TrimSpace(strings.ToLower(line))
        if line == "stop" {
            timer.Stop()
            break
        }
    }
    workTime := time.Now().Sub(startTime)
    return workTime.Truncate(time.Second), nil
}