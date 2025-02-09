package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Define flags for filtering IP addresses.
	omitIPs := flag.Bool("omit-ip", false, "Omit URLs that use an IP address in the host")
	onlyIPs := flag.Bool("only-ip", false, "Only display URLs that use an IP address in the host")
	flag.Parse()

	// Ensure the two flags are not used simultaneously.
	if *omitIPs && *onlyIPs {
		fmt.Fprintln(os.Stderr, "Cannot use both -omit-ip and -only-ip flags simultaneously.")
		os.Exit(1)
	}

	// Create channels for processing input lines and results.
	lines := make(chan string)
	results := make(chan string)

	// This goroutine processes each input line.
	go func() {
		for line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if trimmedLine == "" {
				continue // skip empty lines
			}

			// Parse the URL.
			u, err := url.Parse(trimmedLine)
			if err != nil || u.Scheme == "" || u.Host == "" {
				// Skip invalid URLs (must have a scheme and host)
				continue
			}

			// Extract the host without port.
			host := u.Host
			if strings.Contains(host, ":") {
				// For IPv6 URLs the host is usually in square brackets.
				// net.SplitHostPort handles "host:port" formats.
				if h, _, err := net.SplitHostPort(host); err == nil {
					host = h
				}
			}
			// Remove any surrounding square brackets (for IPv6 literals)
			host = strings.Trim(host, "[]")

			// Determine if the host represents an IP address.
			ipFlag := isIPAddress(host)

			// Apply the filtering flags.
			if *omitIPs && ipFlag {
				continue
			}
			if *onlyIPs && !ipFlag {
				continue
			}

			// If we made it this far, print the entire original line.
			results <- trimmedLine
		}
		close(results)
	}()

	// This goroutine prints the results.
	go func() {
		for result := range results {
			fmt.Println(result)
		}
	}()

	// Read lines from standard input.
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines <- scanner.Text()
	}
	close(lines)

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading from stdin:", err)
	}
}

// isIPAddress returns true if s is recognized as an IP address in any common format.
// It checks for both dotted/colon notation and also attempts to parse s as a numerical value.
func isIPAddress(s string) bool {
	// First try the standard IP parsing.
	if net.ParseIP(s) != nil {
		return true
	}

	// Next, try if the string is a plain decimal integer (e.g. "3232235777")
	if _, err := strconv.ParseUint(s, 10, 32); err == nil {
		return true
	}

	// Lastly, try with base 0 so that hex strings (with 0x prefix) are accepted.
	if _, err := strconv.ParseUint(s, 0, 32); err == nil {
		return true
	}

	return false
}
