package shell

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/chzyer/readline"
	"github.com/evesfect/k-assist/internal/config"
	"github.com/evesfect/k-assist/internal/llm"
)

type Handler struct {
	shellType   string
	logger      *log.Logger
	llmClient   llm.Client
	config      *config.Config
	handleError func(*log.Logger, llm.Client, *config.Config, string)
}

func NewHandler(shellType string, logger *log.Logger, llmClient llm.Client, cfg *config.Config, handleError func(*log.Logger, llm.Client, *config.Config, string)) *Handler {
	if shellType == "" {
		shellType = detectShell()
	}
	return &Handler{
		shellType:   shellType,
		logger:      logger,
		llmClient:   llmClient,
		config:      cfg,
		handleError: handleError,
	}
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
	rl, err := readline.New("")
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
		// Set the prompt with PS1-like style
		rl.SetPrompt(fmt.Sprintf("%s $ ", currentDir))

		// Pre-populate the line with the suggested command
		rl.WriteStdin([]byte(cmd))

		// Get user input (they can edit or just press enter)
		command, err := rl.Readline()
		if err != nil {
			return fmt.Errorf("error reading line: %w", err)
		}

		// Execute the command (original or modified)
		if command = strings.TrimSpace(command); command != "" {
			newDir, err := h.executeCommand(command, currentDir)
			if err != nil {
				h.logger.Printf("Error executing command: %v\n", err)
				h.handleError(h.logger, h.llmClient, h.config, err.Error())
				return nil
			}
			currentDir = newDir
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

	// Create pipes for stdout and stderr
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		// If there's stderr output, include it in the error
		if stderr.Len() > 0 {
			return workDir, fmt.Errorf("%v: %s", err, stderr.String())
		}
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

func (h *Handler) GetHistory(lines int) (string, error) {
	var historyFile string
	var cmd *exec.Cmd

	switch h.shellType {
	case "bash":
		historyFile = filepath.Join(os.Getenv("HOME"), ".bash_history")
	case "zsh":
		historyFile = filepath.Join(os.Getenv("HOME"), ".zsh_history")
	case "powershell":
		cmd = exec.Command("powershell", "-Command",
			fmt.Sprintf("Get-History -Count %d | Format-Table -Property CommandLine -HideTableHeaders", lines))
	default:
		return "", fmt.Errorf("unsupported shell type: %s", h.shellType)
	}

	if h.shellType != "powershell" {
		// For Unix shells, read from history file
		content, err := os.ReadFile(historyFile)
		if err != nil {
			return "", fmt.Errorf("error reading history file: %w", err)
		}

		// Split into lines and get last N lines
		historyLines := strings.Split(string(content), "\n")
		start := len(historyLines) - lines
		if start < 0 {
			start = 0
		}
		return strings.Join(historyLines[start:], "\n"), nil
	}

	// Execute PowerShell command
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting PowerShell history: %w", err)
	}

	return string(output), nil
}
