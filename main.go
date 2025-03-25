package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

	var vbp VBPFile
	var err error

	files := []string{}
	_, ext := getFileDetails(inputfileName)
	if strings.ToLower(ext) == ".vbp" {
		vbp, err = ParseVBPFile(inputfileName)
		if err != nil {
			log.Panic(err)
		}
	} else {
		files = append(files, inputfileName)
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
	files = append(files, vbp.Classes...)
	files = append(files, vbp.Forms...)

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			parseFile(file, outputDir)
		}(file)
	}
	wg.Wait()

	absPath, err := filepath.Abs(outputDir)
	if err != nil {
		log.Panic("couldn't find output directory absolute path")
	}
	fmt.Println("The tool ran successfully, find output files in ", absPath)
}
