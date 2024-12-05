/*
The program removes all lines from file1 that also appear in file2 and outputs the remaining lines. 
It’s useful for filtering out unwanted lines from a file based on a list of lines provided in another file.
*/
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: filterlines file1 file2")
		return
	}

	file1, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file1: %v\n", err)
		return
	}
	defer file1.Close()

	file2, err := os.Open(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file2: %v\n", err)
		return
	}
	defer file2.Close()

	linesToRemove := make(map[string]bool)

	scanner := bufio.NewScanner(file2)
	for scanner.Scan() {
		line := scanner.Text()
		linesToRemove[line] = true
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading file2: %v\n", err)
	}

	scanner = bufio.NewScanner(file1)
	for scanner.Scan() {
		line := scanner.Text()
		if !linesToRemove[line] {
			fmt.Println(line)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading file1: %v\n", err)
	}
}
