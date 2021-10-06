package model

type ErrNotFound struct {
	Message string
}

func (err ErrNotFound) Error() string {
	return err.Message
}
