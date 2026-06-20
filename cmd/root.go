package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/1funct0ry/gizmo/internal/agent"
	"github.com/1funct0ry/gizmo/internal/shell"
	"github.com/spf13/cobra"
)

const (
	defaultBaseURL = "http://localhost:11434/v1" // Ollama's OpenAI-compatible endpoint
	defaultModel   = "qwen3:0.6b"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gizmo",
	Short: "Minimal agent shell for AI-powered coding assistance",
	Long: `Gizmo is a minimal interactive shell that provides AI-powered coding assistance
through Ollama's local language models.

It creates an interactive REPL environment where you can chat with an AI coding
agent to get help with programming tasks, code generation, debugging, and more.

The shell connects to a local Ollama instance and uses configurable models
to provide intelligent responses. History is preserved between sessions.
`,
	Run: func(cmd *cobra.Command, args []string) {
		baseURL, _ := cmd.Flags().GetString("base-url")
		model, _ := cmd.Flags().GetString("model")

		codingAgent := agent.New(context.Background(), baseURL, model)
		sh := shell.NewShell("Gizmo> ", "gizmo.history", codingAgent)
		err := sh.Run()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "error:", err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("base-url", "u", defaultBaseURL, "Base URL for Ollama API")
	rootCmd.Flags().StringP("model", "m", defaultModel, "Model to use for the agent")
}
