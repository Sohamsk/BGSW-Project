package vbptocsproj

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

// Define XML structure for .csproj file
type Project struct {
	XMLName       xml.Name      `xml:"Project"`
	Sdk           string        `xml:"Sdk,attr"`
	PropertyGroup PropertyGroup `xml:"PropertyGroup"`
	ItemGroups    []ItemGroup   `xml:"ItemGroup"`
}

type PropertyGroup struct {
	OutputType      string `xml:"OutputType"`
	TargetFramework string `xml:"TargetFramework"`
	ImplicitUsings  string `xml:"ImplicitUsings"`
	Nullable        string `xml:"Nullable"`
}

type ItemGroup struct {
	Compile          []Compile          `xml:"Compile,omitempty"`
	PackageReference []PackageReference `xml:"PackageReference,omitempty"`
	Reference        []Reference        `xml:"Reference,omitempty"`
}

type Compile struct {
	Include string `xml:"Include,attr"`
}

type PackageReference struct {
	Include string `xml:"Include,attr"`
	Version string `xml:"Version,attr"`
}

type Reference struct {
	Include  string `xml:"Include,attr"`
	HintPath string `xml:"HintPath"`
}

// Function to parse .vbp file and create Project struct for .csproj
func parseVbpFile(filename string) (Project, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Project{}, err
	}
	defer file.Close()

	project := Project{
		Sdk: "Microsoft.NET.Sdk",
		PropertyGroup: PropertyGroup{
			TargetFramework: "net7.0",
			ImplicitUsings:  "enable",
			Nullable:        "enable",
		},
	}

	itemGroup := ItemGroup{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Determine OutputType from Type field in vbp file
		if strings.HasPrefix(line, "Type=") {
			projectType := strings.TrimPrefix(line, "Type=")
			project.PropertyGroup.OutputType = projectType
			//	if projectType == "Exe" {
			//		project.PropertyGroup.OutputType = "Exe"
			//	} else if projectType == "Library" {
			//		project.PropertyGroup.OutputType = "Library"
			//	}

		} else if strings.HasPrefix(line, "Reference=") {
			reference := strings.Split(strings.TrimPrefix(line, "Reference="), "#")
			// Extract the GUID (text inside the `{}` braces)
			guidStart := strings.Index(reference[0], "{")
			guidEnd := strings.Index(reference[0], "}")
			if guidStart != -1 && guidEnd != -1 {
				guid := reference[0][guidStart : guidEnd+1]

				// Extract the namespace name and version
				namespace := reference[3]
				// Check Extension of file Depending on that:
				// 1. .tlb --> do tlb logic
				// 2. .dll --> run tlbimp then get the output path and put it inside HintPath Property

				// Add the parsed COM reference to the ItemGroup
				comReference := Reference{
					Include: namespace,
				}
				itemGroup.Reference = append(itemGroup.Reference, comReference)
			}
		} else if strings.HasPrefix(line, "Object=") {
			objectParts := strings.Split(strings.TrimPrefix(line, "Object="), "#")

			// Extract the GUID from inside the `{}` braces
			guidStart := strings.Index(objectParts[0], "{")
			guidEnd := strings.Index(objectParts[0], "}")
			if guidStart != -1 && guidEnd != -1 {
				guid := objectParts[0][guidStart : guidEnd+1]

				// Extract the major and minor version
				version := strings.Split(objectParts[1], ".")
				majorVersion := parseVersionPart(version[0]) // Helper function to convert to integer
				minorVersion := parseVersionPart(version[1])

				// Extract the control name, which comes after the last `;` separator
				controlName := strings.TrimSpace(strings.Split(objectParts[2], ";")[1])

				// Add the parsed COM reference to the ItemGroup
				comReference := COMReference{
					Include:      controlName,
					Guid:         guid,
					VersionMajor: majorVersion,
					VersionMinor: minorVersion,
				}
				itemGroup.COMReference = append(itemGroup.COMReference, comReference)
			}

		}
	}

	// Add collected items to project
	project.ItemGroups = append(project.ItemGroups, itemGroup)

	if err := scanner.Err(); err != nil {
		return Project{}, err
	}
	return project, nil
}

// Helper function to parse version part to integer, handling errors
func parseVersionPart(part string) int {
	if val, err := strconv.Atoi(part); err == nil {
		return val
	}
	return 0 // Fallback to 0 if parsing fails
}

// Function to write Project struct as .csproj XML file
func writeCsprojFile(project Project, outputFilename string) error {
	file, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	return encoder.Encode(project)
}

func ConvertVBpFiletoCSprojFile(vbpFilePath string) {
	// Specify input and output files
	dir, vbpFileName := path.Split(vbpFilePath)
	vbpFileNameWithoutExtension := strings.Split(vbpFileName, ".")[0]
	outputCsproj := dir + vbpFileNameWithoutExtension + ".csproj" // Desired output csproj filename

	project, err := parseVbpFile(vbpFilePath)
	if err != nil {
		fmt.Println("Error parsing .vbp file:", err)
		return
	}

	err = writeCsprojFile(project, outputCsproj)
	if err != nil {
		fmt.Println("Error writing .csproj file:", err)
		return
	}

	fmt.Printf("Successfully converted %s to %s\n", vbpFilePath, outputCsproj)
}
