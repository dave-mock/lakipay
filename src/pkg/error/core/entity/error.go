package error

type Error struct {
	Type    string
	Message string
}

func (err *Error) Error() string {
	return err.Message
}
