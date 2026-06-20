package agent

import (
	"context"
	"fmt"

	"github.com/1funct0ry/gizmo/internal/toolcall"
	"github.com/1funct0ry/gizmo/internal/toolcall/tools/docker"
	"github.com/1funct0ry/gizmo/internal/toolcall/tools/fs"
	"github.com/1funct0ry/gizmo/internal/toolcall/tools/shell"
	"github.com/1funct0ry/gizmo/internal/utils"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const defaultSystemPrompt = `
You are a helpful terminal-based agent. Your name is Gizmo. Use the available tools when they help you answer the user's request.
`

type Agent struct {
	context      context.Context
	model        string
	client       openai.Client
	messages     []openai.ChatCompletionMessageParamUnion
	registry     *toolcall.Registry
	systemPrompt string
}

func New(ctx context.Context, baseURL string, model string) *Agent {

	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("ollama"), // unused by Ollama, but the SDK requires a non-empty value
	)

	registry := toolcall.NewRegistry()
	registry.AddTool(&fs.ReadTool{})
	registry.AddTool(&shell.ExecuteTool{})
	registry.AddTool(&fs.WriteTool{})
	registry.AddTool(&fs.CWDTool{})
	registry.AddTool(&docker.ImageListTool{})
	registry.AddTool(&docker.ContainerListTool{})
	registry.AddTool(&docker.ContainerLogsTool{})
	registry.AddTool(&docker.ContainerInspectTool{})
	registry.AddTool(&docker.ExecTool{})
	registry.AddTool(&docker.PullTool{})
	registry.AddTool(&docker.StopTool{})
	registry.AddTool(&docker.ComposePsTool{})

	return &Agent{
		context: ctx,
		model:   model,
		client:  client,
		messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(defaultSystemPrompt),
		},
		registry: registry,
	}
}

func (a *Agent) Turn(message string) (string, error) {
	a.messages = append(a.messages, openai.UserMessage(message))
	for {
		completion, err := a.client.Chat.Completions.New(a.context, openai.ChatCompletionNewParams{
			Model:    a.model,
			Messages: a.messages,
			Tools:    a.registry.GetTools(),
		})
		if err != nil {
			return "", err
		}
		if len(completion.Choices) == 0 {
			return "", fmt.Errorf("model returned no choices")
		}

		msg := completion.Choices[0].Message
		a.messages = append(a.messages, msg.ToParam())

		if len(msg.ToolCalls) == 0 {
			return msg.Content, nil
		}

		for _, tc := range msg.ToolCalls {
			result := a.callTool(tc.Function.Name, tc.Function.Arguments)
			fmt.Printf("  \033[2m[result] %s\033[0m\n", utils.Truncate(result, 400))
			a.messages = append(a.messages, openai.ToolMessage(result, tc.ID))
		}
		// Loop again so the model can incorporate the tool results.
	}
}

func (a *Agent) callTool(name, argsJSON string) string {
	tool, ok := a.registry.GetTool(name)

	if !ok {
		return "error: unknown tool " + name
	}
	return tool.Execute(argsJSON)
}
