package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/antlr4-go/antlr/v4"
)

// TODO: Handle inbuilt functions

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\033[31m%s\033[0m\n", r)
		}
	}()

	if len(os.Args) != 2 {
		log.Panic("File Not specified.")
	}

	// TODO: check if .vbp file is given (if it is then parse it seperately and get a list of all modules, classes and forms.
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

	parseFile(input, inputfileName, outputDir)

	absPath, err := filepath.Abs(outputDir)
	if err != nil {
		log.Panic("couldn't find output directory absolute path")
	}
	fmt.Println("The tool ran successfully, find output files in ", absPath)
}
