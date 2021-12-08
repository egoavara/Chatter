package chq

import (
	"fmt"
	"strings"
)

type (
	BuilderError struct {
		stack   []error
		restore []func() error
	}
	StackError struct {
		base          string
		stack         []error
		restoreFailed []error
	}
)

func NewStackError() *BuilderError {
	return &BuilderError{
		stack: []error{},
	}
}
func (s *BuilderError) Handle(handler func() error, restoreHandler ...func() error) *BuilderError {
	if err := handler(); err != nil {
		s.stack = append(s.stack, err)
	} else {
		s.restore = append(s.restore, restoreHandler...)
	}
	return s
}
func (s *BuilderError) Build(base string) *StackError {
	if len(s.stack) > 0 {
		err := &StackError{
			base:          base,
			stack:         s.stack,
			restoreFailed: []error{},
		}
		for _, v := range s.restore {
			if e := v(); e != nil {
				err.restoreFailed = append(err.restoreFailed, e)
			}
		}
		return err
	}
	return nil
}

const _LINE_FORMAT = "%-7s : %s\n"

func (s *StackError) Error() string {
	builder := new(strings.Builder)
	builder.WriteString(fmt.Sprintf(_LINE_FORMAT, "error", s.base))
	builder.WriteString(fmt.Sprintf(_LINE_FORMAT, "stack", s.stack[0].Error()))
	for _, e := range s.stack[1:] {
		builder.WriteString(fmt.Sprintf(_LINE_FORMAT, "", e.Error()))
	}
	if len(s.restoreFailed) > 0 {
		builder.WriteString(fmt.Sprintf(_LINE_FORMAT, "restore", s.restoreFailed[0].Error()))
		for _, e := range s.restoreFailed[1:] {
			builder.WriteString(fmt.Sprintf(_LINE_FORMAT, "", e.Error()))
		}
	}
	return builder.String()
}
