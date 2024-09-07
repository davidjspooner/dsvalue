package dspath

import "fmt"

type ErrInvalidPath struct {
	Path  string
	Inner error
}

func (e *ErrInvalidPath) Error() string {
	return fmt.Errorf("invalid path '%s': %s", e.Path, e.Inner).Error()
}

type ErrEvaluation struct {
	Path  string
	Inner error
}

func (e *ErrEvaluation) Error() string {
	return fmt.Errorf("error evaluating path '%s': %s", e.Path, e.Inner).Error()
}
