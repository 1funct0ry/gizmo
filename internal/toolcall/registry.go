package toolcall

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/shared"
)

type Registry struct {
	tools map[string]Tool
}

func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

func (r *Registry) AddTool(tool Tool) {
	r.tools[tool.Name()] = tool
}

func (r *Registry) GetTools() []openai.ChatCompletionToolParam {
	tools := make([]openai.ChatCompletionToolParam, 0)
	for _, tool := range r.tools {
		tools = append(tools, openai.ChatCompletionToolParam{
			Function: shared.FunctionDefinitionParam{
				Name:        tool.Name(),
				Strict:      param.Opt[bool]{},
				Description: openai.String(tool.Description()),
				Parameters:  tool.Parameters(),
			},
		})
	}
	return tools
}

func (r *Registry) GetTool(name string) (Tool, bool) {
	t, ok := r.tools[name]
	return t, ok
}
