package value

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"golang.org/x/exp/constraints"
)

type genericSimple[T any] struct {
	value  T
	source Source
}

//-------------------------------------------

type stringImpl genericSimple[string]

var _ String = (&stringImpl{})

func (s *stringImpl) Source() Source {
	return s.source
}
func (s *stringImpl) Kind() Kind {
	return StringKind
}
func (s *stringImpl) String() string {
	return s.value
}
func (s *stringImpl) StringOrError() (string, error) {
	return s.value, nil
}
func (s *stringImpl) WithoutSource() interface{} {
	return s.value
}
func (s *stringImpl) CompareTo(other Simple) (int, error) {
	if other, ok := other.(String); ok {
		return strings.Compare(s.value, other.String()), nil
	}
	return 0, fmt.Errorf("cannot compare string to %T", other)
}

func NewString(value string, source Source) String {
	s := &stringImpl{value, source}
	return s
}

//-------------------------------------------

type boolImpl genericSimple[bool]

func (b *boolImpl) Kind() Kind {
	return BoolKind
}
func (b *boolImpl) Source() Source {
	return b.source
}
func (b *boolImpl) String() string {
	return strconv.FormatBool(b.value)
}
func (b *boolImpl) Bool() (bool, error) {
	return b.value, nil
}
func (b *boolImpl) WithoutSource() interface{} {
	return b.value
}
func (b *boolImpl) CompareTo(other Simple) (int, error) {
	if other, ok := other.(Bool); ok {
		otherB, err := other.Bool()
		if err != nil {
			return 0, err
		}
		if b.value == otherB {
			return 0, nil
		}
		if b.value {
			return 1, nil
		}
		return -1, nil
	}
	return 0, fmt.Errorf("cannot compare bool to %T", other)
}

func NewBool(value bool, source Source) Bool {
	return &boolImpl{value, source}
}

//-------------------------------------------

type numberImpl genericSimple[string]

func (n *numberImpl) Int(bits int) (int64, error) {
	return strconv.ParseInt(n.value, 10, bits)
}
func (n *numberImpl) Float(bits int) (float64, error) {
	return strconv.ParseFloat(n.value, bits)
}
func (n *numberImpl) Unsigned(bits int) (uint64, error) {
	return strconv.ParseUint(n.value, 10, bits)
}
func (n *numberImpl) Complex(bits int) (complex128, error) {
	return strconv.ParseComplex(n.value, bits)
}
func (n *numberImpl) Source() Source {
	return n.source
}
func (n *numberImpl) Kind() Kind {
	return NumberKind
}
func (n *numberImpl) String() string {
	return n.value
}
func (n *numberImpl) WithoutSource() interface{} {
	return n.value
}

func (n *numberImpl) CompareTo(other Simple) (int, error) {
	if other, ok := other.(Number); ok {
		otherF, err := other.Float(64)
		if err != nil {
			return 0, err
		}
		f, err := n.Float(64)
		if err != nil {
			return 0, err
		}
		if f == otherF {
			return 0, nil
		}
		if f < otherF {
			return -1, nil
		}
		return 1, nil
	}
	return 0, fmt.Errorf("cannot compare number to %T", other)
}

func NewNumber(value string, source Source) Number {
	_, err := strconv.ParseFloat(value, 64) //todo handle complex
	if err != nil {
		panic(fmt.Errorf("invalid number: %s", value))
	}
	n := &numberImpl{value, source}
	return n
}

func NewInt[T constraints.Integer](value T, source Source) Number {
	n := &numberImpl{strconv.FormatInt(int64(value), 10), source}
	return n
}

func NewUnsigned[T constraints.Unsigned](value T, source Source) Number {
	n := &numberImpl{strconv.FormatUint(uint64(value), 10), source}
	return n
}

func NewFloat[T constraints.Float](value T, source Source) Number {
	bits := int(unsafe.Sizeof(value) * 8)
	n := &numberImpl{strconv.FormatFloat(float64(value), 'f', -1, bits), source}
	return n
}

func NewComplex[T constraints.Complex](value T, source Source) Number {
	bits := int(unsafe.Sizeof(value) * 8)
	return &numberImpl{strconv.FormatComplex(complex128(value), 'f', -1, bits), source}
}

//-------------------------------------------

type Null struct {
	source Source
}

func (n *Null) Kind() Kind {
	return NullKind
}

func (n *Null) Source() Source {
	return n.source
}

func (n *Null) WithoutSource() interface{} {
	return nil
}

func NewNull(source Source) *Null {
	n := &Null{source}
	return n
}
