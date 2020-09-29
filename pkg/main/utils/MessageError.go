package utils

import "fmt"

type MessageError struct {
	StatusCode   int
	ErrorMessage string
}

func (e *MessageError) Error() string {
	return fmt.Sprintf(e.ErrorMessage)
}
