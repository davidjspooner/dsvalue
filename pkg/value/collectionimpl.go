package value

import (
	"fmt"

	"github.com/davidjspooner/dsvalue/pkg/key"
	"golang.org/x/exp/constraints"
)

// --------------------------------------
func NormalizeIndex(originalIndex int, length int) (int, error) {
	index := originalIndex
	if index < 0 {
		index = length + index
		if index < 0 {
			return 0, fmt.Errorf("index out of range: %d", originalIndex)
		}
	}
	if index >= length {
		return 0, fmt.Errorf("index out of range: %d", originalIndex)
	}
	return index, nil
}

// --------------------------------------

type genericArray[T any] struct {
	elements []T
	source   Source
}

func (a *genericArray[T]) Index(index key.Interface) (T, error) {

	iKey, err := index.(key.Value[int])
	if !err {
		var none T
		return none, fmt.Errorf("expected key.Value[int], but got %T", index)
	}

	if iKey.X < 0 || iKey.X >= len(a.elements) {
		var none T
		return none, fmt.Errorf("index out of range: %d", index)
	}
	return a.elements[iKey.X], nil
}

func (a *genericArray[T]) SetIndex(index key.Interface, value T) error {
	iKey, ok := index.(key.Value[int])
	if !ok {
		return fmt.Errorf("expected key.Value[int], but got %T", index)
	}
	fixedIndex, err := NormalizeIndex(iKey.X, len(a.elements))
	if err != nil {
		return err
	}
	a.elements[fixedIndex] = value
	return nil
}

func (a *genericArray[T]) Append(value T) (key.Interface, error) {
	a.elements = append(a.elements, value)
	last := key.Value[int]{X: len(a.elements) - 1}
	return last, nil
}

func (a *genericArray[T]) ForEach(f func(index key.Interface, value T) error) error {
	var v T
	index := key.Value[int]{X: 0}
	for index.X, v = range a.elements {
		if err := f(index, v); err != nil {
			return err
		}
	}
	return nil
}

// --------------------------------------

type arrayImpl struct {
	genericArray[Value]
}

var _ ModifiableArray = &arrayImpl{}

func (a *arrayImpl) Length() (int, error) {
	return len(a.elements), nil
}

func (a *arrayImpl) Source() Source {
	return a.source
}

func (a *arrayImpl) Kind() Kind {
	return ArrayKind
}

func (a *arrayImpl) SetValue(value Value) error {
	return fmt.Errorf("not implemented - SetValue")
}

func (a *arrayImpl) SetIndex(index key.Interface, value Value) error {
	return a.genericArray.SetIndex(index, value)
}

func (a *arrayImpl) Index(index key.Interface) (Value, error) {
	return a.genericArray.Index(index)
}

func (a *arrayImpl) Append(value Value) (key.Interface, error) {
	return a.genericArray.Append(value)
}

func (a *arrayImpl) ForEach(f func(index key.Interface, value Value) error) error {
	return a.genericArray.ForEach(func(index key.Interface, value Value) error {
		return f(index, value)
	})
}

func (a *arrayImpl) WithoutSource() interface{} {
	copy := make([]interface{}, len(a.elements))
	for i, v := range a.elements {
		copy[i] = v.WithoutSource()
	}
	return copy
}

func NewArray(elements []Value, source Source) ModifiableArray {
	return &arrayImpl{
		genericArray: genericArray[Value]{elements: elements, source: source},
	}
}

// --------------------------------------

type genericMap[K constraints.Ordered, T any] struct {
	elements map[K]T
	source   Source
}

func (m *genericMap[K, T]) Field(k key.Interface) (T, error) {
	var none T
	if m.elements == nil {
		return none, fmt.Errorf("field not found: %s", k)
	}
	checkedKey, ok := k.(key.Value[K])
	if !ok {
		return none, fmt.Errorf("expected key.Value[%T], but got %T", m.elements, k)
	}
	v, ok := m.elements[checkedKey.X]
	if !ok {
		return none, fmt.Errorf("field not found: %s", k)
	}
	return v, nil
}

func (m *genericMap[K, T]) SetField(k key.Interface, value T) error {
	if m.elements == nil {
		m.elements = make(map[K]T)
	}
	checkedKey, ok := k.(key.Value[K])
	if !ok {
		return fmt.Errorf("expected key.Value[%T], but got %T", m.elements, k)
	}

	m.elements[checkedKey.X] = value
	return nil
}

func (m *genericMap[K, T]) ForEach(f func(index key.Interface, value T) error) error {
	var k key.Value[K]
	for k.X = range m.elements {
		if err := f(k, m.elements[k.X]); err != nil {
			return err
		}
	}
	return nil
}

// --------------------------------------

type mapImpl struct {
	genericMap[string, Value]
}

var _ ModifiableMap = &mapImpl{}

func (m *mapImpl) Length() (int, error) {
	return len(m.elements), nil
}
func (m *mapImpl) Field(key key.Interface) (Value, error) {
	return m.genericMap.Field(key)
}
func (m *mapImpl) Source() Source {
	return m.source
}

func (m *mapImpl) Kind() Kind {
	return MapKind
}

func (m *mapImpl) SetField(key key.Interface, value Value) error {
	return m.genericMap.SetField(key, value)
}

func (m *mapImpl) SetValue(value Value) error {
	return fmt.Errorf("not implemented - SetValue")
}

func (m *mapImpl) ForEach(f func(index key.Interface, value Value) error) error {
	return m.genericMap.ForEach(func(index key.Interface, value Value) error {
		return f(index, value)
	})
}

func (m *mapImpl) WithoutSource() interface{} {
	copy := make(map[string]any, len(m.elements))
	for k, v := range m.elements {
		copy[k] = v.WithoutSource()
	}
	return copy
}

func NewMap(elements map[string]Value, source Source) Map {
	return &mapImpl{
		genericMap: genericMap[string, Value]{elements: elements, source: source},
	}
}
