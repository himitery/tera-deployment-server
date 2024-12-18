package usecases

type EventProcessor interface {
	Register() error
	Close() error
}
