package main

import (
	"bosch/converter"
	"bosch/listener"
	"bosch/parser"
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\033[31m%s\033[0m\n", r)
		}
	}()

	if len(os.Args) != 2 {
		log.Panic("File Not specified.")
	}

	inputfileName := os.Args[1]

	// Count total lines and multi-line comments
	totalLines, multiLineCommentLines := countLinesAndComments(inputfileName)
	percentageMultiLineComments := 0.0
	if totalLines > 0 {
		percentageMultiLineComments = (float64(multiLineCommentLines) / float64(totalLines)) * 100
	}
	finalResult := 100 - percentageMultiLineComments

	fmt.Printf("Total Lines of Code (LOC): %d\n", totalLines)
	fmt.Printf("Multi-line Comment LOC: %d\n", multiLineCommentLines)
	fmt.Printf("Percentage of Multi-line Comments: %.2f%%\n", percentageMultiLineComments)
	fmt.Printf("Final Result (100 - Multi-line Comment %%): %.2f%%\n", finalResult)

	input, err := antlr.NewFileStream(inputfileName)
	if err != nil {
		log.Panic("File error")
	}

	// Create output directory if it doesn't exist
	outputDir := "output"
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		panic(fmt.Errorf("failed to create output directory: %v", err))
	}

	// Create a logs file
	logfileName := filepath.Join(outputDir, "logs.log")
	logFile, err := os.OpenFile(logfileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Panic("Failed to create log file")
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	fileName, fileExtension := getFileDetails(inputfileName)

	// Initialize lexer and get all tokens
	lexer := parser.NewVisualBasic6Lexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	stream.Fill() // Fill the stream with all tokens

	// Initialize parser
	p := parser.NewVisualBasic6Parser(stream)
	p.BuildParseTrees = true
	tree := p.StartRule()

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	listen := listener.NewTreeShapeListener(writer, &buf)
	writeToOutput(listen, writer, &buf, fileName, fileExtension, tree)
	jsonContent := buf.String()

	convertedContent, err := converter.Convert(jsonContent, listen.SymTab)
	if err != nil {
		log.Panic(err)
	}

	err = writeOutputFiles(fileName, fileExtension, outputDir, jsonContent, convertedContent)
	if err != nil {
		log.Panic("Error writing output files")
	}
}

// countLinesAndComments counts the total lines and comment lines in a VB6 file
func countLinesAndComments(filename string) (int, int) {
	file, err := os.Open(filename)
	if err != nil {
		log.Panic("Cannot open file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	totalLines := 0
	multiLineCommentLines := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		totalLines++

		// Single-line comments in VB6 start with a single quote '
		if strings.HasPrefix(line, "'") {
			multiLineCommentLines++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Panic("Error reading file")
	}

	return totalLines, multiLineCommentLines
}
