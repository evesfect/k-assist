package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/evesfect/k-assist/internal/config"
	"github.com/evesfect/k-assist/internal/dirutil"
	"github.com/evesfect/k-assist/internal/llm"
	"github.com/evesfect/k-assist/internal/shell"
)

func main() {
	// Define flags
	codeFlag := flag.Bool("c", false, "Get code-related information")
	allFlag := flag.Bool("a", false, "Include all subdirectories and files")
	flag.Parse()

	// Check if a prompt is provided
	if flag.NArg() < 1 {
		log.Fatal("Usage: kass [-flag] \"<prompt>\"")
	}

	// Initialize logger
	logger := log.New(os.Stderr, "[kass] ", log.LstdFlags)

	// Get current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		logger.Fatalf("Error getting current directory: %v", err)
	}

	// List directory contents
	var dirInfo string
	if *allFlag {
		dirInfo, err = dirutil.GetAllDirectoryContents(currentDir)
	} else {
		dirInfo, err = dirutil.GetCurrentDirectoryContents(currentDir)
	}
	if err != nil {
		logger.Fatalf("Error reading directory contents: %v", err)
	}

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

	// Process prompt
	prompt := flag.Arg(0)

	// Add current directory information to the prompt
	prompt = dirInfo + "\n" + prompt

	if *codeFlag {
		// Chat (-c flag) single prompt
		response, err := llmClient.GetResponse(prompt)
		if err != nil {
			logger.Fatalf("Error getting response from LLM: %v", err)
		}
		fmt.Println(response)
	} else {
		// Normal operation
		command, err := llmClient.GetCommand(prompt)
		if err != nil {
			logger.Fatalf("Error getting command from LLM: %v", err)
		}

		// Output command for user to edit and potentially execute
		shellHandler := shell.NewHandler(cfg.Shell)
		if err := shellHandler.OutputCommand(command); err != nil {
			logger.Printf("Error with command: %v", err)
		}
	}
}
