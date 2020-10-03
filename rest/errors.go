package rest

import "fmt"

type (
	UserErr struct {
		message string
		inner   error
	}
	ServiceErr struct {
		message string
		inner   error
	}
)

func NewUserErr(message string, err error) error {
	return &UserErr{message: message, inner: err}
}

func (e *UserErr) Error() string {
	return fmt.Sprintf("%v: %v", e.message, e.inner)
}

func NewServiceErr(message string, err error) error {
	return &ServiceErr{message: message, inner: err}
}

func (e *ServiceErr) Error() string {
	if e.message == "" {
		return e.inner.Error()
	}
	return fmt.Sprintf("%v: %v", e.message, e.inner)
}
