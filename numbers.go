package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		numbersLine := extractNumbers(line)
		fmt.Println(numbersLine)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading from stdin:", err)
	}
}

func extractNumbers(s string) string {
	var result []rune
	for _, r := range s {
		if unicode.IsDigit(r) {
			result = append(result, r)
		}
	}
	return string(result)
}
