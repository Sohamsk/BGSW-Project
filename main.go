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

func getFileDetails(inputFileName string) (string, string) {
	filePath := strings.Split(inputFileName, "/")
	fileName := filePath[len(filePath)-1]

	fileNameSlice := strings.Split(fileName, ".")
	fileName, fileExtension := fileNameSlice[0], fileNameSlice[1]
	return fileName, fileExtension
}

func writeOutputFiles(fileName, fileExtension string, jsonContent string, convertedContent string) error {
	// Create output directory if it doesn't exist
	outputDir := "output"
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Write JSON output file
	jsonFilePath := filepath.Join(outputDir, "op.json")
	err = os.WriteFile(jsonFilePath, []byte(jsonContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON file: %v", err)
	}

	// Write converted CS file
	csFilePath := filepath.Join(outputDir, fileName+".cs")
	err = os.WriteFile(csFilePath, []byte(convertedContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write CS file: %v", err)
	}

	return nil
}

func writeToOutput(file *os.File, buf *bytes.Buffer, fileName string, fileExtension string, tree parser.IStartRuleContext) {
	buf.WriteString("{\"FileName\":\"" + fileName + "\", \"FileType\": \"" + fileExtension + "\",")
	writer := bufio.NewWriter(buf)
	antlr.ParseTreeWalkerDefault.Walk(listener.NewTreeShapeListener(writer, buf), tree)
	writer.Flush()
	buf.WriteString("}")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\033[31m%s\033[0m\n", r)
		}
	}()

	if len(os.Args) != 2 {
		panic("File Not specified.")
	}

	inputfileName := os.Args[1]

	input, err := antlr.NewFileStream(inputfileName)
	if err != nil {
		panic("File Not Found:" + err.Error())
	}

	fileName, fileExtension := getFileDetails(inputfileName)

	lexer := parser.NewVisualBasic6Lexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewVisualBasic6Parser(stream)
	p.BuildParseTrees = true
	tree := p.StartRule()

	var buf bytes.Buffer
	writeToOutput(nil, &buf, fileName, fileExtension, tree)
	jsonContent := buf.String()
	//println(jsonContent)

	convertedContent := converter.Convert(jsonContent)

	err = writeOutputFiles(fileName, fileExtension, jsonContent, convertedContent)
	if err != nil {
		log.Panic("Error writing output files")
	}

}
