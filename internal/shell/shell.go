package shell

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/1funct0ry/gizmo/internal/agent"
	"github.com/1funct0ry/gizmo/internal/commands"
	"github.com/ergochat/readline"
)

type Shell struct {
	prompt      string
	historyFile string
	instance    *readline.Instance
	agent       *agent.Agent
	cmds        *commands.Registry
}

func NewShell(prompt string, historyFile string, agent *agent.Agent) *Shell {
	s := &Shell{
		prompt:      prompt,
		historyFile: historyFile,
		agent:       agent,
		cmds:        commands.NewRegistry(),
	}
	s.registerCommands()
	return s
}

func (s *Shell) registerCommands() {
	s.cmds.Register("exit", "exit gizmo", func(_ string) commands.Action {
		fmt.Println("Bye ... 👋")
		return commands.ActionExit
	})
	s.cmds.Register("quit", "exit gizmo (alias for /exit)", func(_ string) commands.Action {
		fmt.Println("Bye ... 👋")
		return commands.ActionExit
	})
	s.cmds.Register("help", "show available slash commands", func(_ string) commands.Action {
		fmt.Println(s.cmds.HelpText())
		return commands.ActionNone
	})
}

func (s *Shell) Run() error {
	var err error
	s.instance, err = readline.NewFromConfig(&readline.Config{
		Prompt:       s.prompt,
		HistoryFile:  s.historyFile,
		AutoComplete: s.cmds.Completer(),
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

		if action, isCmd := s.cmds.Dispatch(line); isCmd {
			if action == commands.ActionExit {
				break
			}
			continue
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
