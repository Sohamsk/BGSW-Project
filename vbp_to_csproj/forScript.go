package vbptocsproj

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func handleCOMReference(dllPath string) (success bool) {
	// func main() {
	cmd := exec.Command("powershell", "-Command", "reg query \"HKCR\\WOW6432Node\\CLSID\" /s | findstr \""+dllPath+"\"")
	// cmd := exec.Command("powershell", "-Command", "reg query \"HKCR\\WOW6432Node\\CLSID\" /s | findstr \"Project1.dll\"")

	// Run the command and capture the output
	output, err := cmd.CombinedOutput()
	strOP := string(output)
	if err != nil {
		fmt.Println("Error in finding Registry entry:", err)
		fmt.Println(strOP)
		return
	}

	AbsolutedllPath := strings.Split(strOP, "   ")[3]

	dllName := filepath.Base(AbsolutedllPath)
	fmt.Println("dllName:", dllName)
	tlbimpCmd := "tlbimp" + strings.TrimSuffix(AbsolutedllPath, "\r\n") + " /out:InteropDlls\\Interop." + dllName
	tlb := exec.Command("powershell", "-Command", strings.TrimSuffix(tlbimpCmd, "\r\n"))
	output, err = tlb.CombinedOutput()
	if err != nil {
		fmt.Println("Error converting using tlbimp:", string(err.Error()))
		fmt.Println(string(output))
		return
	}

	// Print the output
	// fmt.Println(string(output))
	return true
}
