package main

import (
	"bufio"
	"os"
	"regexp"
)

type VBPFile struct {
	References []string
	Forms      []string
	Classes    []string
}

func ParseVBPFile(filename string) (VBPFile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return VBPFile{}, err
	}
	defer file.Close()

	vbp := VBPFile{}

	referencePattern := regexp.MustCompile(`Reference=.#.#.#(.)#`)
	formPattern := regexp.MustCompile(`Form=(..frm)`)
	classPattern := regexp.MustCompile(`Class=(.*.cls)`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if match := referencePattern.FindStringSubmatch(line); match != nil {
			vbp.References = append(vbp.References, match[1])
		} else if match := formPattern.FindStringSubmatch(line); match != nil {
			vbp.Forms = append(vbp.Forms, match[1])
		} else if match := classPattern.FindStringSubmatch(line); match != nil {
			vbp.Classes = append(vbp.Classes, match[1])
		}
	}

	if err := scanner.Err(); err != nil {
		return VBPFile{}, err
	}

	return vbp, nil
}
