package main

import (
	"flag"
	"fmt"
	"github.com/kasouza/go-todo/taskserial"
	"os"
	"strings"
	"unicode/utf8"
)

// Print a table using Unicode box building characters
func printTable(table map[string]string, keyTitle string, valueTitle string) {
	keyWidth := utf8.RuneCountInString(keyTitle)
	valueWidth := utf8.RuneCountInString(valueTitle)

	for k, v := range table {
		if k_width := utf8.RuneCountInString(k); k_width > keyWidth {
			keyWidth = k_width
		}

		if v_width := utf8.RuneCountInString(v); v_width > valueWidth {
			valueWidth = v_width
		}
	}

	totalWidth := keyWidth + valueWidth + 1

	top := []rune(fmt.Sprintf("┌%v┐\n", strings.Repeat("─", totalWidth)))
	top[keyWidth+1] = '┬'

	middle := []rune(fmt.Sprintf("├%v┤\n", strings.Repeat("─", totalWidth)))
	middle[keyWidth+1] = '┼'

	bottom := []rune(fmt.Sprintf("└%v┘\n", strings.Repeat("─", totalWidth)))
	bottom[keyWidth+1] = '┴'

	fmt.Printf(string(top))

	alignedKeyTitle := keyTitle + strings.Repeat(" ", keyWidth-utf8.RuneCountInString(keyTitle))
	alignedValueTitle := valueTitle + strings.Repeat(" ", valueWidth-utf8.RuneCountInString(valueTitle))

	fmt.Printf("│%v│%v│\n", alignedKeyTitle, alignedValueTitle)
	fmt.Printf(string(middle))

	for key, value := range table {
		// These string.Repeat add spaces to align the key and values
		alignedKey := key + strings.Repeat(" ", keyWidth-utf8.RuneCountInString(key))
		alignedValue := value + strings.Repeat(" ", valueWidth-utf8.RuneCountInString(value))

		fmt.Printf("│%v│%v│\n", alignedKey, alignedValue)
	}

	fmt.Printf(string(bottom))
}

// Print a default help message
func printHelp() {
	fmt.Println("CU")
}

// Add a new task or update an old one
func setItem(filename string, args []string) {
	// Unmatched <task-name> <task-description>
	if (len(args)-1)%2 != 0 {
		printHelp()
		return
	}

	tasks := taskserial.ReadTasks(filename)

	for task := 1; task < len(args)-1; task += 2 {
		name := args[task]
		description := args[task+1]

		tasks[name] = taskserial.Task{name, description}
	}

	if len(tasks) > 0 {
		taskserial.WriteTasks(filename, tasks)
	}
}

// Remove one task
func removeItem(filename string, args []string) {
	taskName := args[1]

	tasks := taskserial.ReadTasks(filename)
	delete(tasks, taskName)

	taskserial.WriteTasks(filename, tasks)
}

// Print one task
func getItem(filename string, args []string) {
	if len(args) < 2 {
		printHelp()
		return
	}

	taskName := args[1]

	tasks := taskserial.ReadTasks(filename)
	task, isInMap := tasks[taskName]

	tasksTable := make(map[string]string)

	if isInMap {
		tasksTable[taskName] = task.Description
	}

	printTable(tasksTable, "Task", "Description")
}

// Print all the tasks in a table
func printList(filename string) {
	tasks := taskserial.ReadTasks(filename)
	tasksTable := make(map[string]string)

	for name, task := range tasks {
		tasksTable[name] = task.Description
	}

	printTable(tasksTable, "Task", "Description")
}

func main() {
	var filename string

	home := os.Getenv("HOME")
	flag.StringVar(&filename, "f", home+"/.go-todo", "Path to TODO list file to use instead of the default(~/.go-todo)")

	flag.Parse()
	args := flag.Args()

	CreateFileIfNotExists(filename)

	// If no command was given, print some help
	if len(args) < 1 {
		printHelp()

	} else {
		switch args[0] {
		case "set":
			setItem(filename, args)

		case "get":
			getItem(filename, args)

		case "remove":
			removeItem(filename, args)

		case "print":
			printList(filename)

		default:
			fmt.Println("PORA")
			printHelp()
		}
	}
}
