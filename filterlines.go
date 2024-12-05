package main

import (
        "bufio"
        "flag"
        "fmt"
        "math"
        "os"
        "regexp"
        "strings"

        "github.com/agnivade/levenshtein"
)

func calculateEntropy(s string) float64 {
        if len(s) == 0 {
                return 0
        }
        frequency := make(map[rune]float64)
        for _, char := range s {
                frequency[char]++
        }
        var entropy float64
        for _, count := range frequency {
                probability := count / float64(len(s))
                entropy -= probability * math.Log2(probability)
        }
        return entropy
}

func isRandom(s string, entropyThreshold float64) bool {
        if len(s) < 5 {
                return false
        }
        randomPattern := regexp.MustCompile(`^[a-z]{5,}$`)
        highConsonantRatioPattern := regexp.MustCompile(`^[bcdfghjklmnpqrstvwxyz]{4,}[a-z]*$`)
        entropy := calculateEntropy(s)
        return entropy >= entropyThreshold || randomPattern.MatchString(s) || highConsonantRatioPattern.MatchString(s)
}

func fuzzyMatch(s, target string, threshold float64) bool {
        distance := levenshtein.ComputeDistance(s, target)
        similarity := 1 - float64(distance)/float64(max(len(s), len(target)))
        return similarity >= threshold
}

func max(a, b int) int {
        if a > b {
                return a
        }
        return b
}

func main() {
        // Command-line arguments
        listFile := flag.String("l", "", "File containing input lines")
        charCount := flag.Int("c", -1, "Exclude lines with char count greater than <int>")
        randomThreshold := flag.Float64("r", 3.5, "Exclude lines that are completely random")
        matchRegex := flag.String("mr", "", "Exclude lines matching the regex")
        matchString := flag.String("ms", "", "Exclude lines matching the string (case sensitive)")
        fuzzyString := flag.String("mf", "", "Exclude lines with fuzzy match using Levenshtein similarity")
        similarityThreshold := flag.Float64("t", 0.8, "Fuzzy match similarity threshold")
        showHelp := flag.Bool("h", false, "Show help")

        flag.Parse()

        // Show help if requested
        if *showHelp {
                fmt.Println("Usage:")
                fmt.Println("-l <file>    Input file (optional, otherwise read from stdin)")
                fmt.Println("-c <int>     Exclude lines with char count greater than <int>")
                fmt.Println("-r <float>   Exclude random lines with entropy threshold <float>")
                fmt.Println("-mr <regex>  Exclude lines matching the regex")
                fmt.Println("-ms <string> Exclude lines matching the string (case sensitive)")
                fmt.Println("-mf <string> Exclude lines with fuzzy match using Levenshtein similarity")
                fmt.Println("-t <float>   Fuzzy match similarity threshold (default 0.8)")
                fmt.Println("-h           Show this help")
                return
        }

        // Check if no arguments are provided
        noArgsProvided := *listFile == "" && *charCount == -1 && *matchRegex == "" && *matchString == "" && *fuzzyString == ""

        // Read input lines
        var lines []string
        if *listFile != "" {
                file, err := os.Open(*listFile)
                if err != nil {
                        fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
                        os.Exit(1)
                }
                defer file.Close()
                scanner := bufio.NewScanner(file)
                for scanner.Scan() {
                        lines = append(lines, scanner.Text())
                }
                if err := scanner.Err(); err != nil {
                        fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
                        os.Exit(1)
                }
        } else {
                scanner := bufio.NewScanner(os.Stdin)
                for scanner.Scan() {
                        lines = append(lines, scanner.Text())
                }
                if err := scanner.Err(); err != nil {
                        fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
                        os.Exit(1)
                }
        }

        // Process lines
        for _, line := range lines {
                line = strings.TrimSpace(line)
                if line == "" {
                        continue
                }

                // Apply default rules if no arguments are provided
                if noArgsProvided {
                        if len(line) > 10 || isRandom(line, 3.5) {
                                continue
                        }
                } else {
                        // Exclude by char count
                        if *charCount > 0 && len(line) > *charCount {
                                continue
                        }

                        // Exclude random lines
                        if isRandom(line, *randomThreshold) {
                                continue
                        }

                        // Exclude by regex match
                        if *matchRegex != "" {
                                regex, err := regexp.Compile(*matchRegex)
                                if err != nil {
                                        fmt.Fprintf(os.Stderr, "Invalid regex: %v\n", err)
                                        os.Exit(1)
                                }
                                if regex.MatchString(line) {
                                        continue
                                }
                        }

                        // Exclude by string match
                        if *matchString != "" && line == *matchString {
                                continue
                        }

                        // Exclude by fuzzy match
                        if *fuzzyString != "" && fuzzyMatch(line, *fuzzyString, *similarityThreshold) {
                                continue
                        }
                }

                // Print the line if it passes all checks
                fmt.Println(line)
        }
}
