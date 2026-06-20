package commands

import (
	"fmt"
	"sort"
	"strings"
)

// HandlerFunc is the signature every slash-command handler must satisfy.
// It receives the arguments that follow the command name (trimmed) and
// returns an action that the REPL should take after the handler runs.
type HandlerFunc func(args string) Action

// Action tells the REPL what to do after a command executes.
type Action int

const (
	ActionNone Action = iota // continue the REPL loop
	ActionExit               // break out of the REPL loop
)

// Command pairs a handler with a short help description.
type Command struct {
	Handler     HandlerFunc
	Description string
}

// Registry maps slash-command names (without the leading "/") to Commands.
type Registry struct {
	cmds map[string]Command
}

func NewRegistry() *Registry {
	return &Registry{cmds: make(map[string]Command)}
}

// Register adds or replaces a command. name must not include the leading "/".
func (r *Registry) Register(name, description string, fn HandlerFunc) {
	r.cmds[name] = Command{Handler: fn, Description: description}
}

// Dispatch checks whether line is a slash command. If it is, the handler
// runs and (result, true) is returned. Otherwise (ActionNone, false) is
// returned and the caller should treat the line as a normal chat message.
func (r *Registry) Dispatch(line string) (Action, bool) {
	if !strings.HasPrefix(line, "/") {
		return ActionNone, false
	}
	parts := strings.SplitN(line[1:], " ", 2)
	name := parts[0]
	args := ""
	if len(parts) == 2 {
		args = strings.TrimSpace(parts[1])
	}
	cmd, ok := r.cmds[name]
	if !ok {
		fmt.Printf("unknown command: /%s  (type /help for a list)\n", name)
		return ActionNone, true
	}
	return cmd.Handler(args), true
}

// Completer implements readline.AutoCompleter for slash commands.
// It completes the command name after a leading "/" and leaves everything
// else (normal chat input, mid-word positions) untouched.
func (r *Registry) Completer() *Completer {
	return &Completer{registry: r}
}

// Completer satisfies the readline.AutoCompleter interface.
type Completer struct {
	registry *Registry
}

// Do is called by readline on every TAB press.
//
// readline contract:
//   - line is the full input buffer up to the cursor
//   - pos is the cursor position (== len(line) for end-of-line TAB)
//   - return (candidates, prefixLen) where prefixLen is the number of runes
//     already typed that each candidate extends
func (c *Completer) Do(line []rune, pos int) ([][]rune, int) {
	// Only work on the portion up to the cursor.
	buf := string(line[:pos])

	// Trigger only when the buffer starts with "/" and contains no space
	// (i.e. the user is still typing the command name, not its arguments).
	if !strings.HasPrefix(buf, "/") || strings.ContainsRune(buf, ' ') {
		return nil, 0
	}

	typed := buf[1:] // strip the leading "/"

	var candidates [][]rune
	for name := range c.registry.cmds {
		if strings.HasPrefix(name, typed) {
			// Offer the suffix the user hasn't typed yet.
			candidates = append(candidates, []rune(name[len(typed):]))
		}
	}
	// prefixLen = len(typed): readline will replace that many runes before
	// the cursor with the chosen candidate (so the "/" is preserved).
	return candidates, len([]rune(typed))
}

// HelpText returns a formatted list of all registered commands.
func (r *Registry) HelpText() string {
	names := make([]string, 0, len(r.cmds))
	for n := range r.cmds {
		names = append(names, n)
	}
	sort.Strings(names)

	var b strings.Builder
	b.WriteString("Available commands:\n")
	for _, n := range names {
		fmt.Fprintf(&b, "  /%-12s %s\n", n, r.cmds[n].Description)
	}
	return strings.TrimRight(b.String(), "\n")
}
