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

func writeOutputFiles(fileName, fileExtension, outputDir string, jsonContent string, convertedContent string) error {
	// Write JSON output file
	jsonFilePath := filepath.Join(outputDir, "op.json")
	err := os.WriteFile(jsonFilePath, []byte(jsonContent), 0644)
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

func writeToOutput(listen *listener.TreeShapeListener, writer *bufio.Writer, buf *bytes.Buffer, fileName string, fileExtension string, tree parser.IStartRuleContext) {
	buf.WriteString("{\"FileName\":\"" + fileName + "\", \"FileType\": \"" + fileExtension + "\",")
	antlr.ParseTreeWalkerDefault.Walk(listen, tree)
	writer.Flush()
	buf.WriteString("}")
}

func parseFile(input *antlr.FileStream, inputfileName string, outputDir string) {
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

}
