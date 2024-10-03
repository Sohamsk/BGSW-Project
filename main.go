package main

import (
	"bosch/listener"
	"bosch/parser"
	"bufio"
	"log"
	"os"
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

func writeToOutput(file *os.File, fileName string, fileExtension string, tree parser.IStartRuleContext) {
	w := bufio.NewWriter(file)
	w.WriteString("{\"FileName\":\"" + fileName + "\", \"FileType\": \"" + fileExtension + "\",")
	antlr.ParseTreeWalkerDefault.Walk(listener.NewTreeShapeListener(w), tree)
	file.Seek(-1, 2)
	w.WriteString("}")
	w.Flush()
}

func main() {
	inputfileName := os.Args[1]

	input, err := antlr.NewFileStream(inputfileName)
	if err != nil {
		log.Panic("File error")
	}

	fileName, fileExtension := getFileDetails(inputfileName)

	lexer := parser.NewVisualBasic6Lexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewVisualBasic6Parser(stream)
	p.BuildParseTrees = true
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	tree := p.StartRule()

	f, err := os.Create("op.json")
	if err != nil {
		log.Panic(err)
	}
	f.Seek(0, 0)
	writeToOutput(f, fileName, fileExtension, tree)
	f.Close()
}
