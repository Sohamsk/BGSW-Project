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
	"time"

	"github.com/antlr4-go/antlr/v4"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\033[31m%s\033[0m\n", r)
		}
	}()

	if len(os.Args) != 2 {
		log.Panic("Usage: go run main.go <VB6File>")
	}

	vb6File := os.Args[1]
	startTotal := time.Now()

	startParsing := time.Now()
	csharpFile := generateCSharpFile(vb6File)
	parsingDuration := time.Since(startParsing)

	// Count LOC and comments for VB6
	vb6Total, vb6Comments := countVB6LinesAndComments(vb6File)

	csharpTotal, csharpSingleComments, csharpMultiComments := countCSharpLinesAndComments(csharpFile)

	csharpCodeLOC := csharpTotal - csharpMultiComments

	csharpFinalResult := (float64(csharpCodeLOC) / float64(csharpTotal)) * 100

	totalDuration := time.Since(startTotal)

	fmt.Printf("VB6 - File: %s\n", filepath.Base(vb6File))
	fmt.Printf("  Total LOC: %d, Comment LOC: %d, Code LOC: %d\n", vb6Total, vb6Comments, vb6Total-vb6Comments)

	fmt.Printf("C# - File: %s\n", filepath.Base(csharpFile))
	fmt.Printf("  Total LOC: %d, Single-line Comments: %d, Multi-line Comments: %d, Code LOC: %d, Final Result: %.2f%%\n",
		csharpTotal, csharpSingleComments, csharpMultiComments, csharpCodeLOC, csharpFinalResult)

	// Time measurements
	fmt.Printf("\nTime Measurements:\n")
	fmt.Printf("  Parsing and Conversion Duration: %v\n", parsingDuration)
	fmt.Printf("  Total Conversion Duration: %v\n", totalDuration)
}

func generateCSharpFile(vb6File string) string {
	input, err := antlr.NewFileStream(vb6File)
	if err != nil {
		log.Panic("File error")
	}

	outputDir := "output"
	os.MkdirAll(outputDir, 0755)

	fileName := strings.TrimSuffix(filepath.Base(vb6File), filepath.Ext(vb6File))
	csharpFile := filepath.Join(outputDir, fileName+".cs")

	lexer := parser.NewVisualBasic6Lexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewVisualBasic6Parser(stream)
	p.BuildParseTrees = true
	tree := p.StartRule()

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	listen := listener.NewTreeShapeListener(writer, &buf)
	writeToOutput(listen, writer, &buf, fileName, ".cs", tree)

	convertedContent, err := converter.Convert(buf.String(), listen.SymTab)
	if err != nil {
		log.Panic(err)
	}

	os.WriteFile(csharpFile, []byte(convertedContent), 0644)
	return csharpFile
}

func countVB6LinesAndComments(filename string) (int, int) {
	file, err := os.Open(filename)
	if err != nil {
		log.Panic("Cannot open VB6 file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	totalLines, commentLines := 0, 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		totalLines++
		if strings.HasPrefix(line, "'") {
			commentLines++
		}
	}

	return totalLines, commentLines
}

func countCSharpLinesAndComments(filename string) (int, int, int) {
	file, err := os.Open(filename)
	if err != nil {
		log.Panic("Cannot open C# file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	totalLines, singleLineComments, multiLineComments := 0, 0, 0
	inMultiLineComment := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		totalLines++

		if inMultiLineComment {
			multiLineComments++
			if strings.Contains(line, "*/") {
				inMultiLineComment = false
			}
			continue
		}

		if strings.HasPrefix(line, "/*") {
			multiLineComments++
			inMultiLineComment = true
			continue
		}

		if strings.HasPrefix(line, "//") || strings.Contains(line, "// ") {
			singleLineComments++
		}
	}

	return totalLines, singleLineComments, multiLineComments
}
