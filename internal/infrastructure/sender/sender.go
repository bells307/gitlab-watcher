package sender

// Sender Отправитель сообщений
type Sender interface {
	SendMessage(string) error
}
