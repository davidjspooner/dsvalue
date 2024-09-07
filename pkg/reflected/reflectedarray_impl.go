package reflected

import (
	"fmt"
	"reflect"

	"github.com/davidjspooner/dsvalue/pkg/key"
	"github.com/davidjspooner/dsvalue/pkg/value"
)

type reflectedArrayImpl struct {
	rValue reflect.Value
	source value.Source
}

var _ value.ModifiableArray = &reflectedArrayImpl{}
var _ Reflected = &reflectedArrayImpl{}

func (o *reflectedArrayImpl) Source() value.Source {
	return o.source
}
func (o *reflectedArrayImpl) Kind() value.Kind {
	return value.ArrayKind
}

func (o *reflectedArrayImpl) Index(index key.Interface) (value.Value, error) {
	length := o.rValue.Len()

	nIndex, ok := index.(key.Value[int])
	if !ok {
		return nil, fmt.Errorf("expected key.Value[int], but got %T", index)
	}

	safeIndex, err := value.NormalizeIndex(nIndex.X, length)

	if err != nil {
		return nil, err
	}
	child := o.rValue.Index(safeIndex)
	return NewReflectedObject(child, o.source)
}
func (o *reflectedArrayImpl) Length() (int, error) {
	rk := o.rValue.Kind()
	switch rk {
	case reflect.Array, reflect.Slice:
		return o.rValue.Len(), nil
	default:
		return 0, fmt.Errorf("Length not supported for %s (%s)", o.Kind(), rk)
	}
}

func (o *reflectedArrayImpl) Interface() interface{} {
	return o.rValue.Interface()
}

func (o *reflectedArrayImpl) SetValue(value value.Value) error {
	return fmt.Errorf("not implemented - SetValue")
}

func (o *reflectedArrayImpl) SetIndex(index key.Interface, value value.Value) error {
	return fmt.Errorf("not implemented - SetIndex")
}

func (o *reflectedArrayImpl) Append(value value.Value) (key.Interface, error) {
	return key.Value[int]{}, fmt.Errorf("not implemented - Append")
}

func (o *reflectedArrayImpl) ForEach(f func(index key.Interface, value value.Value) error) error {
	i := key.Value[int]{}
	for i.X = 0; i.X < o.rValue.Len(); i.X++ {
		child, err := NewReflectedObject(o.rValue.Index(i.X), o.source)
		if err != nil {
			return err
		}
		if err = f(i, child); err != nil {
			return err
		}
	}
	return nil
}

func (o *reflectedArrayImpl) WithoutSource() interface{} {
	return o.Interface()
}
