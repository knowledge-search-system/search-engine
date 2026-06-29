package xerrors

import "fmt"

type opError struct {
	op  string
	err error
}

func (e *opError) Error() string {
	return fmt.Sprintf("%s: %s", e.op, e.err)
}

func (e *opError) Unwrap() error {
	return e.err
}

func Wrap(op string, err error) error {
	if err == nil {
		return nil
	}
	return &opError{op: op, err: err}
}
