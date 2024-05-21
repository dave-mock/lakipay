package usecase

type Interactor interface {
	Info(msg string)
	Warning(msg string)
	Error(msg string)
}
