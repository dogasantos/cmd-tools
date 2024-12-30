/*
filter out masscan ips with more than 30 open ports
*/
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func main() {
    scanner := bufio.NewScanner(os.Stdin)

    // We'll store each line plus a reference to the ip:port
    var lines []string
    var ips  []string

    // Keep a map of IP -> set of ports (we’ll simulate with map[string]map[string]bool)
    portsByIP := make(map[string]map[string]bool)

    for scanner.Scan() {
        line := scanner.Text()
        if line == "" {
            continue
        }
        parts := strings.Fields(line)
        if len(parts) < 4 {
            // Not enough fields to parse "open tcp port ip ...", skip or handle differently
            continue
        }

        port := parts[2]
        ip   := parts[3]

        // Save the line and its IP for later
        lines = append(lines, line)
        ips   = append(ips, ip)

        // Initialize a sub-map if needed
        if portsByIP[ip] == nil {
            portsByIP[ip] = make(map[string]bool)
        }
        // Mark that we’ve seen this port for this IP
        portsByIP[ip][port] = true
    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "Error reading input:", err)
        return
    }

    // Now print lines corresponding to IPs having <= 20 unique ports
    for i, line := range lines {
        ip := ips[i]
        if len(portsByIP[ip]) <= 30 {
            fmt.Println(line)
        }
    }
}
