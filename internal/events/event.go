package events

type EventEmitter interface {
	emit() error
}

func Emit(e EventEmitter) error {
	return e.emit()
}
