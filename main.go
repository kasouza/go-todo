package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func print_help() {
	fmt.Println("CU")
}

func set_item(writer *bufio.Writer, args []string) {
	// Unmatched <task-name> <task-description>
	if (len(args)-1)%2 != 0 {
		log.Fatal("ETA RAPAZ CÊ ESCREVEU MAIS NOME DO Q DESCRILÇAO")
	}

	for task := 1; task < len(args); task += 2 {
		_, err := fmt.Fprintf(writer, "%v %v\n", args[task], args[task+1])
		if err != nil {
			log.Fatal(err)
		}
	}
}

func remove_item(filename string, args []string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)

	taskKey := args[1]
	resultingWords := make([]string, 0)

	for scanner.Scan() {
		task := scanner.Text()
		words := strings.Fields(task)

		if words[0] != taskKey {
			resultingWords = append(resultingWords, strings.Join(words, " "))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	file.Close()

	file, err = os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(resultingWords))
	file.WriteString(strings.Join(resultingWords, "\n"))
}

func get_item(scanner *bufio.Scanner, args []string) {
	taskKey := args[1]

	for scanner.Scan() {
		task := scanner.Text()
		words := strings.Fields(task)

		if words[0] == taskKey {
			fmt.Println(strings.Join(words[1:], " "))
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func print_list(scanner *bufio.Scanner) {
	for scanner.Scan() {
		task := scanner.Text()
		words := strings.Fields(task)

		fmt.Printf("%v: %v\n", words[0], strings.Join(words[1:], " "))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	filename := "./temp.txt"
	file, err := os.OpenFile("./temp.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	writer := bufio.NewWriter(file)
	scanner := bufio.NewScanner(file)

	flag.Parse()
	args := flag.Args()

	// If no command was given, print some help
	if len(args) < 1 {
		print_help()

	} else {
		switch args[0] {
		case "set":
			set_item(writer, args)

		case "get":
			get_item(scanner, args)

		case "remove":
			remove_item(filename, args)

		case "print":
			print_list(scanner)

		default:
			fmt.Println("PORA")
			print_help()
		}
	}

	writer.Flush()
}
