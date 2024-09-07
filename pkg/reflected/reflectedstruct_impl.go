package reflected

import (
	"fmt"
	"reflect"

	"github.com/davidjspooner/dsvalue/pkg/key"
	"github.com/davidjspooner/dsvalue/pkg/value"
)

type reflectedStructImpl struct {
	rValue reflect.Value
	source value.Source
}

var _ value.ModifiableMap = &reflectedStructImpl{}
var _ Reflected = &reflectedStructImpl{}

func (o *reflectedStructImpl) Source() value.Source {
	return o.source
}
func (o *reflectedStructImpl) Kind() value.Kind {
	return value.MapKind
}

func (o *reflectedStructImpl) Field(k key.Interface) (value.Value, error) {
	rk := o.rValue.Kind()
	return nil, fmt.Errorf("Field not supported for %s (%s)", o.Kind(), rk)
}

func (o *reflectedStructImpl) Length() (int, error) {
	rk := o.rValue.Kind()
	return 0, fmt.Errorf("Length not supported for %s (%s)", o.Kind(), rk)
}

func (o *reflectedStructImpl) Interface() interface{} {
	return o.rValue.Interface()
}

func (o *reflectedStructImpl) SetValue(value value.Value) error {
	return fmt.Errorf("not implemented - SetValue")
}

func (o *reflectedStructImpl) SetField(key key.Interface, value value.Value) error {
	return fmt.Errorf("not implemented - SetField")
}

func (o *reflectedStructImpl) ForEach(f func(index key.Interface, value value.Value) error) error {
	return fmt.Errorf("not implemented - ForEach")
}
func (o *reflectedStructImpl) WithoutSource() interface{} {
	return o.Interface()
}
