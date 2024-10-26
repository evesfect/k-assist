package llm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/evesfect/k-assist/internal/config"
	"github.com/google/generative-ai-go/genai"
	openai "github.com/sashabaranov/go-openai"
	"google.golang.org/api/option"
)

type Client interface {
	GetCommand(prompt string) (string, error)
}

// Factory function to create the appropriate LLM client
func NewClient(cfg *config.Config) (Client, error) {
	switch cfg.LLM.Provider {
	case "openai":
		return newOpenAIClient(cfg), nil
	case "gemini":
		return newGeminiClient(cfg)
	case "claude":
		return newClaudeClient(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", cfg.LLM.Provider)
	}
}

// OpenAI implementation
type openAIClient struct {
	client *openai.Client
	config *config.Config
}

func newOpenAIClient(cfg *config.Config) *openAIClient {
	return &openAIClient{
		client: openai.NewClient(cfg.LLM.APIKey),
		config: cfg,
	}
}

func (c *openAIClient) GetCommand(prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: c.config.LLM.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: fmt.Sprintf(
						"You are a terminal assistant for %s using %s shell. "+
							"The user is %s. Respond with only the terminal command, "+
							"no explanations.",
						c.config.OS,
						c.config.Shell,
						c.config.User,
					),
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: c.config.MaxTokens,
		},
	)

	if err != nil {
		return "", fmt.Errorf("OpenAI request failed: %w", err)
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}

// Gemini implementation
type geminiClient struct {
	client *genai.Client
	config *config.Config
}

func newGeminiClient(cfg *config.Config) (*geminiClient, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.LLM.APIKey))
	if err != nil {
		return nil, fmt.Errorf("creating Gemini client: %w", err)
	}

	return &geminiClient{
		client: client,
		config: cfg,
	}, nil
}

func (c *geminiClient) GetCommand(prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	model := c.client.GenerativeModel(c.config.LLM.Model)

	systemPrompt := fmt.Sprintf(
		"You are a terminal assistant for %s using %s shell. "+
			"The user is %s. Respond with only the terminal command, "+
			"no explanations.",
		c.config.OS,
		c.config.Shell,
		c.config.User,
	)

	// Combine system prompt and user prompt
	fullPrompt := systemPrompt + "\n\nUser request: " + prompt

	resp, err := model.GenerateContent(ctx, genai.Text(fullPrompt))
	if err != nil {
		return "", fmt.Errorf("gemini request failed: %w", err)
	}

	// Extract the command from the response
	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		for _, part := range resp.Candidates[0].Content.Parts {
			if text, ok := part.(genai.Text); ok {
				return strings.TrimSpace(string(text)), nil
			}
		}
	}

	return "", fmt.Errorf("no valid text response from Gemini")
}

// Claude implementation (placeholder - implement if needed)
type claudeClient struct {
	config *config.Config
}

func newClaudeClient(cfg *config.Config) *claudeClient {
	return &claudeClient{config: cfg}
}

func (c *claudeClient) GetCommand(prompt string) (string, error) {
	return "", fmt.Errorf("claude API not implemented yet")
}
