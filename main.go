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
	defer logFile.Close()

	log.SetOutput(logFile)

	fileName, fileExtension := getFileDetails(inputfileName)

	// Initialize lexer and get all tokens
	lexer := parser.NewVisualBasic6Lexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	stream.Fill() // Fill the stream with all tokens

	// Get the total number of tokens
	totalTokens := len(stream.GetAllTokens())

	// Reset the stream for parsing
	stream.Seek(0)

	// Initialize parser
	p := parser.NewVisualBasic6Parser(stream)
	p.BuildParseTrees = true
	tree := p.StartRule()

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	listen := listener.NewTreeShapeListener(writer, &buf)
	writeToOutput(listen, writer, &buf, fileName, fileExtension, tree)
	jsonContent := buf.String()

	// Calculate the percentage of tokens processed
	processedTokens := stream.Index()
	conversionPercentage := float64(processedTokens) / float64(totalTokens) * 100
	fmt.Printf("Conversion Progress: %.2f%%\n", conversionPercentage)

	convertedContent, err := converter.Convert(jsonContent, listen.SymTab)
	if err != nil {
		log.Panic(err)
	}

	err = writeOutputFiles(fileName, fileExtension, outputDir, jsonContent, convertedContent)
	if err != nil {
		log.Panic("Error writing output files")
	}
}
