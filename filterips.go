package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

func main() {
	includePorts := flag.Bool("ports", false, "include ports in the output")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	lines := make(chan string)
	results := make(chan string)

	go func() {
		for line := range lines {
			if *includePorts {
				if isValidIPWithPort(line) {
					results <- line
				}
			} else {
				if isValidIP(line) {
					results <- line
				}
			}
		}
		close(results)
	}()

	go func() {
		for result := range results {
			fmt.Println(result)
		}
	}()

	for scanner.Scan() {
		lines <- scanner.Text()
	}
	close(lines)

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading from stdin:", err)
	}
}

func isValidIPWithPort(s string) bool {
	if strings.Contains(s, ":") {
		if strings.Count(s, ":") > 1 {
			host, port, err := net.SplitHostPort(s)
			if err != nil {
				return false
			}
			if isValidIP(host) && isValidPort(port) {
				return true
			}
		}
		re := regexp.MustCompile(`^(\d+\.\d+\.\d+\.\d+):(\d+)$`)
		if matches := re.FindStringSubmatch(s); matches != nil {
			if isValidIP(matches[1]) && isValidPort(matches[2]) {
				return true
			}
		}
	}
	return isValidIP(s)
}

func isValidIP(s string) bool {
	return net.ParseIP(s) != nil
}

func isValidPort(port string) bool {
	re := regexp.MustCompile(`^\d+$`)
	return re.MatchString(port)
}
