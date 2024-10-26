package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type LLMConfig struct {
	Provider string `json:"provider"` // "openai", "gemini", or "claude"
	APIKey   string `json:"api_key"`
	Model    string `json:"model"`
}

type Config struct {
	OS        string    `json:"os"`
	User      string    `json:"user"`
	LLM       LLMConfig `json:"llm"`
	MaxTokens int       `json:"max_tokens"`
	Shell     string    `json:"shell"`
}

// Default configuration values
const (
	DefaultMaxTokens = 50
	ConfigFileName   = "config.json"

	// Default models for each provider
	DefaultOpenAIModel = "gpt-3.5-turbo"
	DefaultGeminiModel = "gemini-pro"
	DefaultClaudeModel = "claude-3-sonnet-20240229"
)

// Load reads and parses the configuration file
func Load() (*Config, error) {
	configPath, err := ensureConfigFile()
	if err != nil {
		return nil, fmt.Errorf("config file error: %w", err)
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := validateAndSetDefaults(&config); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	return &config, nil
}

// validateAndSetDefaults validates the config and sets default values
func validateAndSetDefaults(config *Config) error {
	// Set MaxTokens default
	if config.MaxTokens == 0 {
		config.MaxTokens = DefaultMaxTokens
	}

	// Validate LLM configuration
	if config.LLM.Provider == "" {
		return fmt.Errorf("LLM provider must be specified")
	}

	// Set default model based on provider if not specified
	if config.LLM.Model == "" {
		switch config.LLM.Provider {
		case "openai":
			config.LLM.Model = DefaultOpenAIModel
		case "gemini":
			config.LLM.Model = DefaultGeminiModel
		case "claude":
			config.LLM.Model = DefaultClaudeModel
		default:
			return fmt.Errorf("unsupported LLM provider: %s", config.LLM.Provider)
		}
	}

	// Check for API key in environment variables if not in config
	if config.LLM.APIKey == "" {
		envVar := fmt.Sprintf("KASS_%s_API_KEY", config.LLM.Provider)
		config.LLM.APIKey = os.Getenv(envVar)
		if config.LLM.APIKey == "" {
			return fmt.Errorf("API key not found in config or environment variable %s", envVar)
		}
	}

	return nil
}

// ensureConfigFile ensures the config file exists and returns its path
func ensureConfigFile() (string, error) {
	// First, check current directory
	if _, err := os.Stat(ConfigFileName); err == nil {
		return ConfigFileName, nil
	}

	// Then check user config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	kassConfigDir := filepath.Join(configDir, "kass")
	configPath := filepath.Join(kassConfigDir, ConfigFileName)

	// If config doesn't exist in user config directory, create it
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create directory if it doesn't exist
		if err := os.MkdirAll(kassConfigDir, 0755); err != nil {
			return "", err
		}

		// Create default config file
		defaultConfig := Config{
			OS:        "linux", // This should be detected
			User:      os.Getenv("USER"),
			MaxTokens: DefaultMaxTokens,
			LLM: LLMConfig{
				Provider: "gemini", // Default to Gemini
				Model:    DefaultGeminiModel,
			},
		}

		configJSON, err := json.MarshalIndent(defaultConfig, "", "    ")
		if err != nil {
			return "", err
		}

		if err := os.WriteFile(configPath, configJSON, 0644); err != nil {
			return "", err
		}
	}

	return configPath, nil
}
