package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Show lines that are present on file 1 but are not present in file 2 (line orde does not matter)")
		fmt.Println("Usage: uniquelines file1 file2")
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

	linesFile1 := make(map[string]bool)
	linesFile2 := make(map[string]bool)

	scanner := bufio.NewScanner(file1)
	for scanner.Scan() {
		line := scanner.Text()
		linesFile1[line] = true
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading file1: %v\n", err)
	}

	scanner = bufio.NewScanner(file2)
	for scanner.Scan() {
		line := scanner.Text()
		linesFile2[line] = true
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading file2: %v\n", err)
	}

	//fmt.Println("Lines unique to file1:")
	for line := range linesFile1 {
		if !linesFile2[line] {
			fmt.Println(line)
		}
	}

	//fmt.Println("\nLines unique to file2:")
	//for line := range linesFile2 {
	//	if !linesFile1[line] {
	//		fmt.Println(line)
	//	}
	//}
}
