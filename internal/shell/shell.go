package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/chzyer/readline"
)

type Handler struct {
	shellType string
}

func NewHandler(shellType string) *Handler {
	if shellType == "" {
		shellType = detectShell()
	}
	return &Handler{shellType: shellType}
}

func detectShell() string {
	if runtime.GOOS == "windows" {
		return "powershell"
	}
	shell := os.Getenv("SHELL")
	if strings.Contains(shell, "zsh") {
		return "zsh"
	}
	return "bash"
}

func (h *Handler) FormatCommand(command string) string {
	switch h.shellType {
	case "powershell":
		return fmt.Sprintf("echo %s | Out-String", command)
	case "cmd":
		return fmt.Sprintf("echo %s", command)
	default: // bash, zsh, and other Unix shells
		return fmt.Sprintf("echo '%s'", strings.Replace(command, "'", "'\\''", -1))
	}
}

func (h *Handler) OutputCommand(commands string) error {
	rl, err := readline.New("> ")
	if err != nil {
		return fmt.Errorf("error creating readline instance: %w", err)
	}
	defer rl.Close()

	// Split commands by newline
	cmdList := strings.Split(strings.TrimSpace(commands), "\n")

	// Keep track of the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	for _, cmd := range cmdList {
		fmt.Printf("Command: %s\n", cmd)
		rl.SetPrompt("Execute? (y/n/q): ")

		for {
			response, err := rl.Readline()
			if err != nil {
				return fmt.Errorf("error reading line: %w", err)
			}

			response = strings.ToLower(strings.TrimSpace(response))
			if response == "y" {
				newDir, err := h.executeCommand(cmd, currentDir)
				if err != nil {
					fmt.Printf("Error executing command: %v\n", err)
				} else {
					currentDir = newDir
				}
				break
			} else if response == "n" {
				break
			} else if response == "q" {
				return nil
			} else {
				fmt.Println("Please enter 'y' to execute, 'n' to skip, or 'q' to quit.")
			}
		}
	}

	return nil
}

func (h *Handler) executeCommand(command, workDir string) (string, error) {
	var cmd *exec.Cmd

	switch h.shellType {
	case "powershell":
		cmd = exec.Command("powershell", "-Command", command)
	case "cmd":
		cmd = exec.Command("cmd", "/C", command)
	default:
		cmd = exec.Command("sh", "-c", command)
	}

	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return workDir, err
	}

	// Check if the command was a cd command
	if strings.HasPrefix(strings.TrimSpace(command), "cd ") {
		// Extract the new directory from the cd command
		newDir := strings.TrimSpace(strings.TrimPrefix(command, "cd "))
		// Resolve relative paths
		if !filepath.IsAbs(newDir) {
			newDir = filepath.Join(workDir, newDir)
		}
		// Verify the new directory exists
		if _, err := os.Stat(newDir); err == nil {
			return newDir, nil
		}
	}

	return workDir, nil
}
