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

	// create a logs file
	logfileName := filepath.Join(outputDir, "logs.log")
	logFile, err := os.OpenFile(logfileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Could not create log file but program execution will continue")
	}
	defer logFile.Close()
	if err == nil {
		log.SetOutput(logFile)
	}
	fileName, fileExtension := getFileDetails(inputfileName)

	lexer := parser.NewVisualBasic6Lexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewVisualBasic6Parser(stream)
	p.BuildParseTrees = true
	tree := p.StartRule()

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	listen := listener.NewTreeShapeListener(writer, &buf)
	writeToOutput(listen, writer, &buf, fileName, fileExtension, tree)
	jsonContent := buf.String()
	// fmt.Println(jsonContent) //uncomment this while debugging json
	// start debug
	//for key, val := range listen.SymTab {
	//	fmt.Printf("%s: %s\n", key, val)
	//}
	// stop debug

	convertedContent, err := converter.Convert(jsonContent, listen.SymTab)
	if err != nil {
		log.Panic(err)
	}

	err = writeOutputFiles(fileName, fileExtension, outputDir, jsonContent, convertedContent)
	if err != nil {
		log.Panic("Error writing output files")
	}
	absPath, err := filepath.Abs(outputDir)
	if err != nil {
		panic("couldn't find output directory absolute path")
	}
	fmt.Println("The tool ran successfully, find output files in ", absPath)
}
