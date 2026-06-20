package shell

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/1funct0ry/gizmo/internal/agent"
	"github.com/ergochat/readline"
)

type Shell struct {
	prompt      string
	historyFile string
	instance    *readline.Instance
	agent       *agent.Agent
}

func NewShell(prompt string, historyFile string, agent *agent.Agent) *Shell {
	return &Shell{prompt: prompt, historyFile: historyFile, instance: nil, agent: agent}
}

func (s *Shell) Run() error {
	var err error
	s.instance, err = readline.NewFromConfig(&readline.Config{
		Prompt:      s.prompt,
		HistoryFile: s.historyFile,
	})
	if err != nil {
		return err
	}
	defer func(rl *readline.Instance) {
		_ = rl.Close()
	}(s.instance)
	s.instance.CaptureExitSignal()

	for {
		line, err := s.instance.Readline()
		if errors.Is(err, readline.ErrInterrupt) {
			continue
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "input error:", err)
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "/exit" || line == "/quit" {
			break
		}

		reply, err := s.agent.Turn(line)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "error:", err)
			continue
		}
		fmt.Println(reply)
	}

	return nil
}
