package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
)

func main() {
	includePorts := flag.Bool("ports", false, "include ports in the output")
	flag.Parse()

	ipv4Regex := regexp.MustCompile(`(\b\d{1,3}(\.\d{1,3}){3}\b)(:\d+)?`)
	ipv6Regex := regexp.MustCompile(`(\b[0-9a-fA-F:]+\b)(:\d+)?`)

	scanner := bufio.NewScanner(os.Stdin)
	lines := make(chan string)
	results := make(chan string)

	go func() {
		for line := range lines {
			if *includePorts {
				matches := ipv4Regex.FindAllString(line, -1)
				matches = append(matches, ipv6Regex.FindAllString(line, -1)...)
				for _, match := range matches {
					if isValidIP(match) || isValidIPWithPort(match) {
						results <- match
					}
				}
			} else {
				matches := ipv4Regex.FindAllStringSubmatch(line, -1)
				for _, match := range matches {
					if isValidIP(match[1]) {
						results <- match[1]
					}
				}
				matches = ipv6Regex.FindAllStringSubmatch(line, -1)
				for _, match := range matches {
					if isValidIP(match[1]) {
						results <- match[1]
					}
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
	host, port, err := net.SplitHostPort(s)
	if err != nil {
		return false
	}
	return isValidIP(host) && isValidPort(port)
}

func isValidIP(s string) bool {
	return net.ParseIP(s) != nil
}

func isValidPort(port string) bool {
	re := regexp.MustCompile(`^\d+$`)
	if !re.MatchString(port) {
		return false
	}
	p, err := strconv.Atoi(port)
	if err != nil || p < 1 || p > 65535 {
		return false
	}
	return true
}
