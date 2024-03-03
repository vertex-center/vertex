package port

type LogsService interface {
	Push(content string) error
}
