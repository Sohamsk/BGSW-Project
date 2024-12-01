package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"bosch/converter"
	"bosch/listener"
	"bosch/parser"

	"github.com/antlr4-go/antlr/v4"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	choices   []string
	cursor    int
	selected  *int
	textInput textinput.Model
	status    string
	inputMode bool
	err       error
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter absolute path to VB6 source file"
	ti.Focus()

	return model{
		choices:   []string{"1. VB6 to IR", "2. VB6 to C#"},
		textInput: ti,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if !m.inputMode {
				m.cursor--
				if m.cursor < 0 {
					m.cursor = len(m.choices) - 1
				}
			}

		case "down", "j":
			if !m.inputMode {
				m.cursor++
				if m.cursor >= len(m.choices) {
					m.cursor = 0
				}
			}

		case "enter":
			if !m.inputMode {
				m.selected = &m.cursor
				m.inputMode = true
				m.textInput.Focus()
				return m, nil
			}
			return m, m.processFile

		case "esc":
			if m.inputMode {
				m.inputMode = false
				m.selected = nil
			}
		}

		if m.inputMode {
			m.textInput, cmd = m.textInput.Update(msg)
		}

	case errMsg:
		m.err = msg
		m.status = ""

	case successMsg:
		m.status = msg.status
		m.err = nil
		m.inputMode = false
		m.selected = nil
		m.textInput.Reset()
	}

	return m, cmd
}

func (m model) View() string {
	s := strings.Builder{}

	// Title
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#5D4CFF")).
		Bold(true).
		Padding(0, 1).
		Render("VB6 to C# Code Generation Wizard")

	s.WriteString(fmt.Sprintf("\n%s\n\n", title))

	// Error handling
	if m.err != nil {
		s.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("red")).
			Render(m.err.Error()) + "\n")
	}

	// Success message
	if m.status != "" {
		s.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("green")).
			Render(m.status) + "\n")
	}

	// Conversion type selection
	if m.selected == nil {
		s.WriteString("Select Conversion Type:\n")
		for i, choice := range m.choices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			s.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
		}
		s.WriteString("\nUse ↑/↓ to navigate, Enter to select, q to quit\n")
	}

	// File input
	if m.inputMode {
		conversionType := "Unknown"
		if m.selected != nil {
			switch *m.selected {
			case 0:
				conversionType = "VB6 to IR"
			case 1:
				conversionType = "VB6 to C#"
			}
		}
		s.WriteString(fmt.Sprintf("Selected: %s\n", conversionType))
		s.WriteString("Enter file path:\n")
		s.WriteString(m.textInput.View())
		s.WriteString("\nPress Enter to convert, Esc to go back\n")
	}

	return s.String()
}

func (m *model) processFile() tea.Msg {
	if m.selected == nil {
		return errMsg{fmt.Errorf("no conversion type selected")}
	}

	inputFileName := m.textInput.Value()
	if inputFileName == "" {
		return errMsg{fmt.Errorf("file path cannot be empty")}
	}

	// Validate file exists
	if _, err := os.Stat(inputFileName); os.IsNotExist(err) {
		return errMsg{fmt.Errorf("file does not exist: %s", inputFileName)}
	}

	// Generate output file path
	fileName, fileExtension := getFileDetails(inputFileName)
	outputFile := filepath.Join(".", fileName+"_output.json")

	// Parse VB6 file
	input, err := antlr.NewFileStream(inputFileName)
	if err != nil {
		return errMsg{fmt.Errorf("error reading input file: %v", err)}
	}

	lexer := parser.NewVisualBasic6Lexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewVisualBasic6Parser(stream)
	p.BuildParseTrees = true
	tree := p.StartRule()

	// Create output file
	f, err := os.Create(outputFile)
	if err != nil {
		return errMsg{fmt.Errorf("error creating output file: %v", err)}
	}
	defer f.Close()

	var buf bytes.Buffer
	writeToOutput(f, &buf, fileName, fileExtension, tree)

	// Perform conversion based on selected type
	var status string
	switch *m.selected {
	case 0: // VB6 to IR
		status = "Generated Intermediate Representation (IR)"
	case 1: // VB6 to C#
		converter.Convert(buf.String())
		status = "Converted to C# source code"
	}

	return successMsg{status: status}
}

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

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

type successMsg struct{ status string }

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
