package shell

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
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

func (h *Handler) OutputCommand(command string) error {
	var cmd *exec.Cmd

	switch h.shellType {
	case "powershell":
		cmd = exec.Command("powershell", "-Command", command)
	case "cmd":
		cmd = exec.Command("cmd", "/C", command)
	default:
		cmd = exec.Command("sh", "-c", command)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
