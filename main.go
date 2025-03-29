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
	csharpFile, jsonFile := generateCSharpAndJSONFile(vb6File)
	parsingDuration := time.Since(startParsing)

	// Count LOC and comments for VB6
	vb6Total, vb6Comments := countVB6LinesAndComments(vb6File)

	csharpTotal, csharpSingleComments, csharpMultiComments := countCSharpLinesAndComments(csharpFile)

	csharpCodeLOC := csharpTotal - csharpMultiComments
	csharpFinalResult := (float64(csharpCodeLOC) / float64(csharpTotal)) * 100

	totalDuration := time.Since(startTotal)

	fmt.Printf("\n==================== Conversion Summary ====================\n")

// VB6 File Details
fmt.Printf("\n📂 VB6 File: %s\n", filepath.Base(vb6File))
fmt.Printf("   📌 Total LOC: %-6d | 📝 Comment LOC: %-6d | 💻 Code LOC: %-6d\n", vb6Total, vb6Comments, vb6Total-vb6Comments)

// C# File Details
fmt.Printf("\n🚀 C# File: %s\n", filepath.Base(csharpFile))
fmt.Printf("   📌 Total LOC: %-6d | 📝 Single-line Comments: %-6d | 📝 Multi-line Comments: %-6d\n", 
          csharpTotal, csharpSingleComments, csharpMultiComments)
fmt.Printf("   💻 Code LOC: %-6d | ✅ Final Result: %.2f%%\n", csharpCodeLOC, csharpFinalResult)

// JSON File Details
fmt.Printf("\n📂 JSON File: %s\n", filepath.Base(jsonFile))

// Time Measurements
fmt.Printf("\n⏳ Time Measurements:\n")
fmt.Printf("   ⏱️ Parsing and Conversion Duration: %v\n", parsingDuration)
fmt.Printf("   ⏱️ Total Conversion Duration      : %v\n", totalDuration)
fmt.Printf("\n===========================================================\n")
}

func generateCSharpAndJSONFile(vb6File string) (string, string) {
	input, err := antlr.NewFileStream(vb6File)
	if err != nil {
		log.Panic("File error")
	}

	outputDir := "output"
	os.MkdirAll(outputDir, 0755)

	fileName := strings.TrimSuffix(filepath.Base(vb6File), filepath.Ext(vb6File))
	csharpFile := filepath.Join(outputDir, fileName+".cs")
	jsonFile := filepath.Join(outputDir, fileName+".json")

	lexer := parser.NewVisualBasic6Lexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewVisualBasic6Parser(stream)
	p.BuildParseTrees = true
	tree := p.StartRule()

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	listen := listener.NewTreeShapeListener(writer, &buf)
	writeToOutput(listen, writer, &buf, fileName, ".cs", tree)

	jsonContent := buf.String()
	generateJSONFile(jsonFile, jsonContent)

	convertedContent, err := converter.Convert(jsonContent, listen.SymTab)
	if err != nil {
		log.Panic(err)
	}

	err = os.WriteFile(csharpFile, []byte(convertedContent), 0644)
	if err != nil {
		log.Panic("Error writing C# file")
	}

	return csharpFile, jsonFile
}

func generateJSONFile(jsonFile, jsonContent string) {
	err := os.WriteFile(jsonFile, []byte(jsonContent), 0644)
	if err != nil {
		log.Panic("Error writing JSON file")
	}
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
		}
		if strings.HasPrefix(line, "//") {
			singleLineComments++
		}
	}

	return totalLines, singleLineComments, multiLineComments
}
