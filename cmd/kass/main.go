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

	// Check prompt
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
		response, err := llmClient.GetResponse(prompt)
		if err != nil {
			logger.Printf("Error getting response from LLM: %v", err)
			handleErrorWithAssistance(logger, llmClient, cfg, err.Error())
			return
		}
		fmt.Println(response)
	} else {
		command, err := llmClient.GetCommand(prompt)
		if err != nil {
			logger.Printf("Error getting command from LLM: %v", err)
			handleErrorWithAssistance(logger, llmClient, cfg, err.Error())
			return
		}

		// Output command for user to edit and execute
		shellHandler := shell.NewHandler(cfg.Shell, logger, llmClient, cfg, handleErrorWithAssistance)
		if err := shellHandler.OutputCommand(command); err != nil {
			logger.Printf("Error with command: %v", err)
			handleErrorWithAssistance(logger, llmClient, cfg, err.Error())
			return
		}
	}
}

func handleErrorWithAssistance(logger *log.Logger, llmClient llm.Client, cfg *config.Config, errResponse string) {
	fmt.Printf("Would you like assistance with this error? [Y/n] ")
	var willAssist string
	fmt.Scanln(&willAssist)

	if willAssist == "Y" || willAssist == "y" {
		shellHandler := shell.NewHandler(cfg.Shell, logger, llmClient, cfg, handleErrorWithAssistance)
		history, err := shellHandler.GetHistory(20)
		if err != nil {
			logger.Printf("Warning: Could not get shell history: %v", err)
			history = "No shell history available"
		}

		response, err := llmClient.HandleError(errResponse, history)
		if err != nil {
			logger.Printf("Error getting assistance: %v", err)
			return
		}
		fmt.Println(response)
	}
}
