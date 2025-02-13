package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// countLOC counts total lines and multi-line comment lines
func countLOC(filename string) (int, int, float64) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0, 0.0
	}
	defer file.Close()

	totalLines := 0
	multiCommentLines := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		totalLines++

		// Check if the line is a VB6 comment (starts with ' or Rem)
		if strings.HasPrefix(line, "'") || strings.HasPrefix(strings.ToLower(line), "rem ") {
			multiCommentLines++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return 0, 0, 0.0
	}

	// Calculate the percentage of multi-line comments
	var percentage float64
	if totalLines > 0 {
		percentage = (float64(multiCommentLines) / float64(totalLines)) * 100
	}

	// Compute the final result (100 - percentage)
	finalResult := 100 - percentage

	return totalLines, multiCommentLines, finalResult
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <converted_vb6_file>")
		return
	}

	filename := os.Args[1]
	totalLOC, multiCommentLOC, finalPercentage := countLOC(filename)

	fmt.Printf("Total Lines of Code (LOC): %d\n", totalLOC)
	fmt.Printf("Multi-line Comment LOC: %d\n", multiCommentLOC)
	fmt.Printf("Percentage of Multi-line Comments: %.2f%%\n", 100-finalPercentage)
	fmt.Printf("Final Result (100 - Multi-line Comment %%): %.2f%%\n", finalPercentage)
}
