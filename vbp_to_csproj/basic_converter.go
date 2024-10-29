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
	RootNamespace   string `xml:"RootNamespace"`
	AssemblyName    string `xml:"AssemblyName"`
}

type ItemGroup struct {
	Compile          []Compile          `xml:"Compile,omitempty"`
	PackageReference []PackageReference `xml:"PackageReference,omitempty"`
	COMReference     []COMReference     `xml:"COMReference,omitempty"`
}

type Compile struct {
	Include string `xml:"Include,attr"`
}

type PackageReference struct {
	Include string `xml:"Include,attr"`
	Version string `xml:"Version,attr"`
}

type COMReference struct {
	Include      string `xml:"Include,attr"`
	Guid         string `xml:"Guid"`
	VersionMajor int    `xml:"VersionMajor"`
	VersionMinor int    `xml:"VersionMinor"`
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
			TargetFramework: "net5.0",
			RootNamespace:   "ConvertedVBProject",
			AssemblyName:    "ConvertedVBProject",
		},
	}

	itemGroup := ItemGroup{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "Name=") {
			projectName := strings.Trim(strings.TrimPrefix(line, "Name="), "\"")
			project.PropertyGroup.RootNamespace = projectName
			project.PropertyGroup.AssemblyName = projectName
		}
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
				version := strings.Split(reference[1], ".")

				// Parse version components as integers
				majorVersion := 0
				minorVersion := 0
				if len(version) >= 1 {
					majorVersion = parseVersionPart(version[0]) // helper function to convert to integer
				}
				if len(version) >= 2 {
					minorVersion = parseVersionPart(version[1])
				}

				// Add the parsed COM reference to the ItemGroup
				comReference := COMReference{
					Include:      namespace,
					Guid:         guid,
					VersionMajor: majorVersion,
					VersionMinor: minorVersion,
				}
				itemGroup.COMReference = append(itemGroup.COMReference, comReference)
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
