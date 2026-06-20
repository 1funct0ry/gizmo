# Gizmo

> A minimal terminal-based AI coding agent powered by your local Ollama models.

[![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE.md)
[![Status](https://img.shields.io/badge/status-alpha-orange.svg)](#roadmap)
[![Ollama](https://img.shields.io/badge/powered%20by-Ollama-black.svg)](https://ollama.com/)
[![Build with Make](https://img.shields.io/badge/build-make-blue.svg)](#building--development)

Gizmo runs an interactive REPL that talks to a local [Ollama](https://ollama.com/) instance through its OpenAI-compatible API, and lets the model call local tools — filesystem, shell, and Docker — to act on your machine. Everything runs locally; no cloud API keys required.

---

## Features

- 🧠 **Local-first** — connects to a local Ollama server; your code and prompts never leave your machine.
- 🔁 **Agentic tool loop** — the model can chain multiple tool calls in a single turn until it has what it needs to answer.
- 🛠️ **Built-in tools** — filesystem, shell, and Docker operations out of the box.
- 💬 **Persistent REPL** — readline-based prompt with command history saved between turns.
- 🔌 **OpenAI-compatible** — uses the standard OpenAI Go SDK pointed at Ollama, so any Ollama-served model works.
- ⚙️ **Configurable** — swap the model or Ollama endpoint with a single flag.

## How it works

```
main.go → cmd.Execute() (Cobra) → agent.New() + shell.NewShell() → shell.Run() REPL
                                                                         │
                                                          each line → agent.Turn()
                                                                         │
                                    ┌────────────────────────────────────┘
                                    ▼
              chat completion (history + tool schemas) ──► tool calls? ──► execute tools
                                    ▲                                          │
                                    └──────────── append results, loop ────────┘
                                                  (until no tool calls)
```

Conversation history is held in memory for the lifetime of the process.

## Requirements

- **Go 1.26+** (only needed to build from source)
- **A running [Ollama](https://ollama.com/) server** with at least one model pulled
- **Docker** (optional — only required to use the Docker tools)

```bash
# Install and start Ollama, then pull the default model:
ollama pull qwen3:0.6b
```

## Installation

```bash
git clone https://github.com/1funct0ry/gizmo.git
cd gizmo
make build        # produces ./bin/gizmo
```

## Usage

```bash
make run                       # build and launch the REPL
# or, after building:
./bin/gizmo

# point at a different model or Ollama endpoint:
./bin/gizmo --model llama3.2 --base-url http://localhost:11434/v1
```

| Flag         | Short | Default                     | Description                                   |
|--------------|-------|-----------------------------|-----------------------------------------------|
| `--base-url` | `-u`  | `http://localhost:11434/v1` | Base URL for the Ollama OpenAI-compatible API |
| `--model`    | `-m`  | `qwen3:0.6b`                | Model to use for the agent                    |

Once inside the REPL, type your request at the `Gizmo>` prompt. Type `/exit` or `/quit` (or press `Ctrl-D`) to leave. Command history is written to `gizmo.history` in the working directory.

```text
Gizmo> list the running docker containers and tell me which one exposes the most ports
  [result] ...
...
```

## Available tools

The agent can call any of these tools autonomously while answering.

### Filesystem
| Tool                | Description                                            |
|---------------------|--------------------------------------------------------|
| `read_file`         | Read and return the contents of a text file from disk. |
| `write_file`        | Write a string to a text file on disk.                 |
| `current_directory` | Return the current working directory.                  |

### Shell
| Tool                | Description                                                                         |
|---------------------|-------------------------------------------------------------------------------------|
| `run_shell_command` | Execute a shell command on the local machine and return its combined stdout/stderr. |

### Docker
| Tool                | Description                                                                 |
|---------------------|-----------------------------------------------------------------------------|
| `docker_images`     | List all Docker images on the local machine.                                |
| `docker_containers` | List containers with status and ports (`all=true` to include stopped ones). |
| `docker_logs`       | Fetch the last N lines of a container's logs (default 50).                  |
| `docker_inspect`    | Return the full JSON inspect output for a container.                        |
| `docker_exec`       | Run a command inside a running container and return its output.             |
| `docker_pull`       | Pull an image from a registry onto the local machine.                       |
| `docker_stop`       | Stop a running container (graceful SIGTERM, then SIGKILL after 10s).        |
| `docker_compose_ps` | Show the status of services in a Docker Compose project.                    |


## Building & development

```bash
make build        # build to bin/gizmo (injects Version/BuildTime via ldflags)
make run          # build then run the REPL
make test         # go test -v ./...
make fmt          # go fmt ./...
make vet          # go vet ./...
make clean        # remove bin/
make help         # list all targets
```

> **Note:** there are no tests in the repo yet, so `make test` currently exercises nothing.

### Adding a tool

Registering a tool is a one-line change. Implement the `Tool` interface
(`Name`, `Description`, `Execute(args string) string`, `Parameters() map[string]any`)
under `internal/toolcall/tools/<group>/`, then register it in `agent.New()`:

```go
registry.AddTool(&yourpkg.YourTool{})
```

The agent does not auto-discover tools — if it isn't added there, the model can't call it. Tools report failures as human-readable `"Error: ..."` strings (not Go errors), since that string is fed straight back to the model.


## Tech stack

- [Go](https://go.dev/) 1.26
- [openai-go](https://github.com/openai/openai-go) — chat + tool-calling client (pointed at Ollama)
- [spf13/cobra](https://github.com/spf13/cobra) — CLI framework
- [ergochat/readline](https://github.com/ergochat/readline) — REPL input
- [moby/moby](https://github.com/moby/moby) — Docker client

## License

Released under the [MIT License](LICENSE.md). © 2026 `One Functory`.
