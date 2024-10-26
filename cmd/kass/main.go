package main

import (
	"log"
	"os"

	"github.com/evesfect/k-assist/internal/config"
	"github.com/evesfect/k-assist/internal/llm"
	"github.com/evesfect/k-assist/internal/shell"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: kass \"<prompt>\"")
	}

	// Initialize logger
	logger := log.New(os.Stderr, "[kass] ", log.LstdFlags)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}

	// Create LLM client
	llmClient, err := llm.NewClient(cfg)
	if err != nil {
		logger.Fatalf("Error creating LLM client: %v", err)
	}

	// Get shell handler
	shellHandler := shell.NewHandler(cfg.Shell)

	// Process prompt and get command
	prompt := os.Args[1]
	command, err := llmClient.GetCommand(prompt)
	if err != nil {
		logger.Fatalf("Error getting command from LLM: %v", err)
	}

	// Format and output command
	formattedCommand := shellHandler.FormatCommand(command)
	if err := shellHandler.OutputCommand(formattedCommand); err != nil {
		logger.Fatalf("Error outputting command: %v", err)
	}
}
