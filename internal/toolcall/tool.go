package toolcall

type Tool interface {
	Name() string
	Description() string
	Execute(args string) string
	Parameters() map[string]any
}
