package runner

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/viper"
)

// RunShellCommand executes a command string in the appropriate shell for the OS.
// It connects stdin, stdout, and stderr to match the current terminal session.
func Run(commandStr string) error {
	var cmd *exec.Cmd

	shell := viper.GetString("shell")
	if shell != "" {
		fmt.Println("Using shell configured in fli config", shell)
		cmd = exec.Command(shell, "-c", commandStr)
	} else {
		// Prioritize user-defined shell regardless of OS
		if shell = os.Getenv("SHELL"); shell != "" {
			fmt.Println("Using user-defined SHELL to run command:", shell)
			cmd = exec.Command(shell, "-c", commandStr)
		} else if runtime.GOOS == "windows" {
			// Fall back to PowerShell or cmd on Windows
			if _, err := exec.LookPath("powershell"); err == nil {
				fmt.Println("Using PowerShell to run command")
				cmd = exec.Command("powershell", "-Command", commandStr)
			} else {
				fmt.Println("Using cmd to run command")
				cmd = exec.Command("cmd", "/C", commandStr)
			}
		} else {
			// Fall back to sh on Unix-like systems
			fmt.Println("Using default shell 'sh' to run command")
			cmd = exec.Command("sh", "-c", commandStr)
		}
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
