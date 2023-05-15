# tasklist
a tool to keep track of tasks and how long you've been working on them

## Running
```
go run main.go
```

## Usage

Add a new task:
```
> add <task description>
```

Update an existing task:
```
> update <tasknumber> <new task description>
```

Remove a task:
```
> remove <tasknumber>
```

Work on a task. When you press enter again, the duration will be added to the task:
```
> work <tasknumber>
```

Finish task:
```
> finish <tasknumber>
```
