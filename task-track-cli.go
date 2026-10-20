package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updateAt"`
}

type Query struct {
	NextID int
	Tasks  []Task
}

func LoadTask() []Task {
	data, err := os.ReadFile("task.json")
	if err != nil {
		return []Task{}
	}
	var tasks []Task
	json.Unmarshal(data, &tasks)
	return tasks
}

func SaveTask(tasks []Task) {
	data, _ := json.MarshalIndent(tasks, "", "  ")
	os.WriteFile("task.json", data, 0644)
}

func Add(description string) {
	tasks := LoadTask()
	var newID int = 1

	for i := 0; i < len(tasks); i++ {
		newID = max(newID, tasks[i].ID+1)
	}

	newTask := Task{
		ID:          newID,
		Description: description,
		Status:      "todo",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tasks = append(tasks, newTask)
	SaveTask(tasks)
	fmt.Printf("Task added successfully (ID: %d)\n", newID)
}

func List() {
	tasks := LoadTask()
	if len(tasks) == 0 {
		fmt.Print("No tasks found")
		return
	}
	for i := 0; i < len(tasks); i++ {
		fmt.Printf("ID: %d\n", tasks[i].ID)
		fmt.Printf("Description: %s\n", tasks[i].Description)
		fmt.Printf("Status: %s\n", tasks[i].Status)
		fmt.Printf("CreatedAt: %v\n", tasks[i].CreatedAt)
		fmt.Printf("UpdateAt: %v\n", tasks[i].UpdatedAt)
	}
}

func ReadQuery() (int, string, error) {
	if len(os.Args) < 4 {
		return 0, "", fmt.Errorf("usage : update <id> <description>")
	}

	id, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return 0, "", fmt.Errorf("id must be a number")
	}

	description := strings.Join(os.Args[3:], " ")
	if description == "" {
		return 0, "", fmt.Errorf("description cannot be empty")
	}
	return id, description, nil
}

func Update(id int, description string) {

	tasks := LoadTask()
	if len(tasks) == 0 {
		fmt.Print("Not found 404")
		return
	}

	for i := 0; i < len(tasks); i++ {
		if tasks[i].ID == id {
			tasks[i].Description = description
			tasks[i].UpdatedAt = time.Now()
			SaveTask(tasks)
			fmt.Println("Task updated successfully")
			return
		}
	}

	fmt.Println("Task not found :V")
}

// hàm update , sẽ viết 1 hàm trong đó gồm các chức năng là đọc truy vấn
// và lấy ra id , chuỗi xâu cần update lại
// sau đó truy cập vô tasks bằng cách lấy về thông tin từ json(hàm LoadTask)
// Update lại description , udpateat và sau đó save task

func Delete(id int) {

	tasks := LoadTask()
	if len(tasks) == 0 {
		fmt.Print("Not Found")
		return
	}

	for i := 0; i < len(tasks); i++ {
		if tasks[i].ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			SaveTask(tasks)
			fmt.Println("Task delete successfully")
			return
		}
	}
	fmt.Print("Not Found")
	return
}

func MarkInProgress(id int) {

	tasks := LoadTask()
	if len(tasks) == 0 {
		fmt.Print("NOT FOUND")
		return
	}

	for i := 0; i < len(tasks); i++ {
		if tasks[i].ID == id {
			tasks[i].Status = "in-progress"
			tasks[i].UpdatedAt = time.Now()
			SaveTask(tasks)
			fmt.Println("Task updated successfully")
			return
		}
	}
	fmt.Print("NOT FOUND")
	return
}

func MarkInDone(id int) {

	tasks := LoadTask()
	if len(tasks) == 0 {
		fmt.Print("NOT FOUND")
		return
	}

	for i := 0; i < len(tasks); i++ {
		if tasks[i].ID == id {
			tasks[i].Status = "done"
			tasks[i].UpdatedAt = time.Now()
			SaveTask(tasks)
			fmt.Println("Task updated successfully")
			return
		}
	}
	fmt.Print("NOT FOUND")
	return
}

func ListStatus(query string) {

	tasks := LoadTask()
	if len(tasks) == 0 {
		fmt.Print("NOT FOUND")
		return
	}

	for i := 0; i < len(tasks); i++ {
		if tasks[i].Status == query {
			fmt.Printf("ID: %d\n", tasks[i].ID)
			fmt.Printf("Description: %s\n", tasks[i].Description)
			fmt.Printf("Status: %s\n", tasks[i].Status)
			fmt.Printf("CreatedAt: %v\n", tasks[i].CreatedAt)
			fmt.Printf("UpdateAt: %v\n", tasks[i].UpdatedAt)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: task-cli <command>")
		return
	}

	command := os.Args[1]
	defer fmt.Println()

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli add <description>")
			return
		}

		description := strings.Join(os.Args[2:], " ")
		Add(description)

	case "list":
		if len(os.Args) == 2 {
			List()
			return
		}

		if len(os.Args) != 3 {
			fmt.Println("Usage: task-cli list [todo|in-progress|done]")
			return
		}

		status := os.Args[2]
		if status != "todo" && status != "in-progress" && status != "done" {
			fmt.Println("Invalid status. Use todo, in-progress, or done.")
			return
		}

		ListStatus(status)

	case "update":
		id, description, err := ReadQuery()
		if err != nil {
			fmt.Println(err)
			return
		}

		Update(id, description)

	case "delete", "mark-in-progress", "mark-done":
		if len(os.Args) != 3 {
			fmt.Printf("Usage: task-cli %s <id>\n", command)
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil || id <= 0 {
			fmt.Println("ID must be a positive number")
			return
		}

		switch command {
		case "delete":
			Delete(id)
		case "mark-in-progress":
			MarkInProgress(id)
		case "mark-done":
			MarkInDone(id)
		}

	default:
		fmt.Println("Unknown command:", command)
		fmt.Println("Commands: add, list, update, delete, mark-in-progress, mark-done")
	}
}
