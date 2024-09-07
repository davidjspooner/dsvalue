package dsformat

import (
	"io"

	"github.com/davidjspooner/dsvalue/pkg/value"
)

type FormatOption func(Interface) error

type Interface interface {
	Description() string
	NewWithOptions(options ...FormatOption) (Interface, error)
	NewEncoder(writer io.Writer, options ...FormatOption) (Encoder, error)
	NewDecoder(reader io.Reader, options ...FormatOption) (Decoder, error)
}

type Encoder interface {
	Encode(source value.Value) error
}
type Decoder interface {
	Decode(target value.Source) error
}
