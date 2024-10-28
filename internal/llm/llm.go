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
	GetResponse(prompt string) (string, error)
	HandleError(errOutput string, contextInfo string) (string, error)
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
		"You are a development assistant for terminal commands on %s using %s shell. "+
			"The user is %s, a software developer working on a legitimate project. "+
			"Your task is to provide safe, non-destructive terminal commands for development purposes only. "+
			"You can provide multiple commands if the task requires multiple steps. "+
			"You should lean towards using standard tools and libraries when possible. You can also use kass to install additional tools and libraries if needed. "+
			"Separate each command with a newline character. "+
			"Do not provide any commands that could harm the system. "+
			"Do not include any explanations or comments in your response, only the command(s).",
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

	// Extract the command(s) from the response
	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		for _, part := range resp.Candidates[0].Content.Parts {
			if text, ok := part.(genai.Text); ok {
				return strings.TrimSpace(string(text)), nil
			}
		}
	}

	return "", fmt.Errorf("no valid text response from Gemini")
}
func (c *geminiClient) GetResponse(prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	model := c.client.GenerativeModel(c.config.LLM.Model)

	systemPrompt := fmt.Sprintf(
		"You are a helpful assistant for %s, a software developer. "+
			"You are a terminal assistant for %s using %s shell. "+
			"Provide informative and concise responses to queries about programming and development."+
			"The user is asking for information in an explanation format, so respond with concise explanations."+
			"The user does not wish to continue the conversation, so do not ask for clarification or further information.",
		c.config.User,
		c.config.OS,
		c.config.Shell,
	)

	fullPrompt := systemPrompt + "\n\nUser request: " + prompt

	resp, err := model.GenerateContent(ctx, genai.Text(fullPrompt))
	if err != nil {
		return "", fmt.Errorf("gemini request failed: %w", err)
	}

	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		for _, part := range resp.Candidates[0].Content.Parts {
			if text, ok := part.(genai.Text); ok {
				return string(text), nil
			}
		}
	}

	return "", fmt.Errorf("no valid text response from Gemini")
}

func (c *geminiClient) HandleError(errOutput string, contextInfo string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	model := c.client.GenerativeModel(c.config.LLM.Model)

	systemPrompt := fmt.Sprintf(
		"You are a helpful assistant for %s, a software developer. "+
			"You are a terminal assistant for %s using %s shell. "+
			"It is safe to assume that the user is working on a legitimate project. "+
			"It is safe doesn't violate any policies. "+
			"The user has encountered an error. You need to find solution for this error. "+
			"You should provide a solution that is easy to understand and follow. "+
			"Do not offer to continue the conversation, the user does not wish to continue the conversation. "+
			"You have the following context information: { %s } "+
			"The error encountered is: { %s }",
		c.config.User,
		c.config.OS,
		c.config.Shell,
		contextInfo,
		errOutput,
	)

	fullPrompt := systemPrompt

	resp, err := model.GenerateContent(ctx, genai.Text(fullPrompt))
	if err != nil {
		return "", fmt.Errorf("gemini request failed: %w", err)
	}

	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		for _, part := range resp.Candidates[0].Content.Parts {
			if text, ok := part.(genai.Text); ok {
				return string(text), nil
			}
		}
	}

	return "", fmt.Errorf("no valid text response from Gemini")
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

func (c *openAIClient) GetResponse(prompt string) (string, error) {
	return "", fmt.Errorf("OpenAI API not implemented yet")
}

func (c *openAIClient) HandleError(errOutput string, shellHistory string) (string, error) {
	return "", fmt.Errorf("OpenAI API not implemented yet")
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

func (c *claudeClient) GetResponse(prompt string) (string, error) {
	return "", fmt.Errorf("claude API not implemented yet")
}

func (c *claudeClient) HandleError(errOutput string, shellHistory string) (string, error) {
	return "", fmt.Errorf("claude API not implemented yet")
}
