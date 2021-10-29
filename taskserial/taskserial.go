package taskserial

import (
	"encoding/json"
	"log"
	"os"
)

type Task struct {
	Name        string
	Description string
}

// A wrapper to `os.OpenFile()` to deal with errors
func openFileWithFlags(filename string, flags int) *os.File {
	file, err := os.OpenFile(filename, flags, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return file
}

func CreateFileIfNotExists(filename string) {
	file := openFileWithFlags(filename, os.O_CREATE)
	file.close()
}

// Write the given tasks to a file, overriding its content or
// creating it if needed
func WriteTasks(filename string, tasks map[string]Task) {
	file := openFileWithFlags(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE)
	defer file.Close()

	taskList := make([]Task, 0, len(tasks))

	for _, task := range tasks {
		taskList = append(taskList, task)
	}

	b, err := json.Marshal(taskList)
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(b)
	if err != nil {
		log.Fatal(err)
	}
}

// Read tasks from a file and return them in a `map[name]description`
// If the file does not exists, it will create one.
func ReadTasks(filename string) map[string]Task {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	tasks := make(map[string]Task, 0)
	taskList := make([]Task, 0)

	// If file is not empty, parse it into json
	if len(fileContent) > 0 {
		err = json.Unmarshal(fileContent, &taskList)
		if err != nil {
			log.Fatal(err)
		}

		for _, task := range taskList {
			tasks[task.Name] = task
		}
	}

	return tasks
}
