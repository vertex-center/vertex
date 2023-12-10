package port

type LogsAdapter interface {
	Push(content string) error
}
