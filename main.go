package main

import (
	"bosch/converter"
	"bosch/listener"
	"bosch/parser"
	//	vbptocsproj "bosch/vbp_to_csproj"
	"bufio"
	"bytes"
	"fmt"
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

func writeToOutput(file *os.File, buf *bytes.Buffer, fileName string, fileExtension string, tree parser.IStartRuleContext) {
	buf.WriteString("{\"FileName\":\"" + fileName + "\", \"FileType\": \"" + fileExtension + "\",")
	writer := bufio.NewWriter(buf)
	antlr.ParseTreeWalkerDefault.Walk(listener.NewTreeShapeListener(writer, buf), tree)
	writer.Flush()
	buf.WriteString("}")
	file.WriteString(buf.String())
}

func main() {
	//	vbptocsproj.ConvertVBpFiletoCSprojFile("./vbp_to_csproj/Complex.vbp") // TODO : just an example file change later
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
	//p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	tree := p.StartRule()

	f, err := os.Create("op.json")
	if err != nil {
		log.Panic(err)
	}
	f.Seek(0, 0)
	var buf bytes.Buffer
	writeToOutput(f, &buf, fileName, fileExtension, tree)
	f.Close()

	fmt.Println(buf.String())
	converter.Convert(buf.String())
}
