package errors

import (
	"fmt"
	"io"
)

// base Error.
type base struct {
	msg string
	*stack
}

func (e *base) Error() string {
	return e.msg
}

func (e *base) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, e.msg)
			e.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, e.msg)
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", e.msg)
	}
}

// withCode Error.
type withCode struct {
	err   error
	code  int
	cause error
	*stack
}

func (w *withCode) Error() string { return fmt.Sprintf("%v", w) }

func (w *withCode) Cause() error { return w.cause }

func (w *withCode) Unwrap() error { return w.cause }

// withStack Error.
type withStack struct {
	error
	*stack
}

func (w *withStack) Cause() error { return w.error }

func (w *withStack) Unwrap() error {
	if e, ok := w.error.(interface{ Unwrap() error }); ok {
		return e.Unwrap()
	}

	return w.error
}

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v", w.Cause())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, w.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", w.Error())
	}
}

// withMessage Error.
type withMessage struct {
	cause error
	msg   string
}

func (w *withMessage) Error() string { return w.msg }

func (w *withMessage) Cause() error { return w.cause }

func (w *withMessage) Unwrap() error { return w.cause }

func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v\n", w.Cause())
			_, _ = io.WriteString(s, w.msg)
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, w.Error())
	}
}
