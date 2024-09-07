package value

import (
	"fmt"
)

type unknownSource struct{}

func (u *unknownSource) String() string {
	return "<unknown>"
}

var UnknownSource Source = &unknownSource{}

type Position struct {
	Line   int
	Column int
}

type SourcePosition struct {
	position Position
	source   Source
}

func NewSourcePosition(source Source) *SourcePosition {
	sub, ok := source.(*SourcePosition)
	if ok {
		copy := *sub
		return &copy
	}
	return &SourcePosition{Position{1, 1}, source}
}

func (s *SourcePosition) String() string {
	if s.source == nil {
		return fmt.Sprintf("[Ln=%d,Col=%d]", s.position.Line, s.position.Column)
	}
	return fmt.Sprintf("%s [Ln=%d,Col=%d]", s.source.String(), s.position.Line, s.position.Column)
}

func (s *SourcePosition) Position() Position {
	return s.position
}

func (s *SourcePosition) Advance(delta Position) {
	s.position.Line += delta.Line
	if delta.Line > 0 {
		s.position.Column = 1
	}
	s.position.Column += delta.Column
}
