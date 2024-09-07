package reflected

import (
	"fmt"
	"reflect"

	"github.com/davidjspooner/dsvalue/pkg/key"
	"github.com/davidjspooner/dsvalue/pkg/value"
)

type reflectedMapImpl struct {
	rValue reflect.Value
	source value.Source
}

var _ value.ModifiableMap = &reflectedMapImpl{}
var _ Reflected = &reflectedMapImpl{}

func (o *reflectedMapImpl) Source() value.Source {
	return o.source
}
func (o *reflectedMapImpl) Kind() value.Kind {
	return value.MapKind
}

func (o *reflectedMapImpl) Field(k key.Interface) (value.Value, error) {

	safeKey, ok := k.(key.Value[string])
	if !ok {
		return nil, fmt.Errorf("expected key.Value[string], but got %T", k)
	}

	child := o.rValue.MapIndex(reflect.ValueOf(safeKey.X))
	if !child.IsValid() {
		return nil, fmt.Errorf("Field %q not found", k)
	}
	return NewReflectedObject(child, o.source)
}

func (o *reflectedMapImpl) Length() (int, error) {
	rk := o.rValue.Kind()
	return 0, fmt.Errorf("Length not supported for %s (%s)", o.Kind(), rk)
}

func (o *reflectedMapImpl) Interface() interface{} {
	return o.rValue.Interface()
}

func (o *reflectedMapImpl) SetValue(value value.Value) error {
	return fmt.Errorf("not implemented - SetValue")
}

func (o *reflectedMapImpl) SetField(k key.Interface, value value.Value) error {
	return fmt.Errorf("not implemented - SetField")
}

func (o *reflectedMapImpl) ForEach(f func(index key.Interface, value value.Value) error) error {
	oKeys := o.rValue.MapKeys()
	i := o.rValue.Interface()
	_ = i
	for _, k := range oKeys {
		reflectedChild := o.rValue.MapIndex(k)
		child, err := NewReflectedObject(reflectedChild, o.source)
		i2 := reflectedChild.Interface()
		_ = i2
		if err != nil {
			return err
		}
		ks := k.String()
		if err = f(key.Value[string]{X: ks}, child); err != nil {
			return err
		}
	}
	return nil
}
func (o *reflectedMapImpl) WithoutSource() interface{} {
	return o.Interface()
}
