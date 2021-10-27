package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
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

// A wrapper to `os.OpenFile()` to deal with errors
func openFileWithFlags(filename string, flags int) *os.File {
	file, err := os.OpenFile(filename, flags, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return file
}

// Read tasks from a file and return them in a `map[name]description`
// If the file does not exists, it will create one.
func readTasks(filename string) map[string]string {
	file := openFileWithFlags(filename, os.O_RDONLY|os.O_CREATE)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	tasks := make(map[string]string, 0)

	for scanner.Scan() {
		task := scanner.Text()
		words := strings.Fields(task)

		tasks[words[0]] = strings.Join(words[1:], " ")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return tasks
}

// Write the given tasks to a file, overriding its content or
// creating it if necessary
func writeTasks(filename string, tasks map[string]string) {
	file := openFileWithFlags(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE)
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for name, description := range tasks {
		fmt.Fprintf(writer, "%v %v\n", name, description)
	}
}

// Add a new task or update an old one
func setItem(filename string, args []string) {
	// Unmatched <task-name> <task-description>
	if (len(args)-1)%2 != 0 {
		printHelp()
		return
	}

	tasks := readTasks(filename)

	for task := 1; task < len(args)-1; task += 2 {
		name := args[task]
		if len(strings.Fields(name)) > 1 {
			log.Fatal("Task name need to be a single word")
		}

		description := args[task+1]

		tasks[name] = description
	}

	writeTasks(filename, tasks)
}

// Remove one task
func removeItem(filename string, args []string) {
	taskName := args[1]

	tasks := readTasks(filename)
	delete(tasks, taskName)

	writeTasks(filename, tasks)
}

// Print one task
func getItem(filename string, args []string) {
	if len(args) < 2 {
		printHelp()
		return
	}

	taskName := args[1]

	tasks := readTasks(filename)

	printTable(map[string]string{taskName: tasks[taskName]}, "Task", "Description")
}

// Print all the tasks in a table
func printList(filename string) {
	file := openFileWithFlags(filename, os.O_RDONLY)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	wordsTable := make(map[string]string)

	for scanner.Scan() {
		task := scanner.Text()
		words := strings.Fields(task)

		wordsTable[words[0]] = strings.Join(words[1:], " ")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	printTable(wordsTable, "Task", "Description")
}

func main() {
	var filename string

	home := os.Getenv("HOME")
	flag.StringVar(&filename, "f", home+"/.go-todo", "Path to TODO list file to use instead of the default(~/.go-todo)")

	flag.Parse()
	args := flag.Args()

	// Make sure the TODO list file exists
	file := openFileWithFlags(filename, os.O_CREATE)
	file.Close()

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
