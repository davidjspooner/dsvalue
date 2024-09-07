package reflected

import (
	"fmt"
	"reflect"

	"github.com/davidjspooner/dsvalue/pkg/value"
)

type reflectedSimpleImpl struct {
	rValue reflect.Value
	source value.Source
}

var _ value.Simple = &reflectedSimpleImpl{}
var _ value.ModifiableValue = &reflectedSimpleImpl{}
var _ Reflected = &reflectedSimpleImpl{}

func (o *reflectedSimpleImpl) Source() value.Source {
	return o.source
}
func (o *reflectedSimpleImpl) Kind() value.Kind {
	rk := o.rValue.Kind()
	switch rk {
	case reflect.String:
		return value.StringKind
	case reflect.Bool:
		return value.BoolKind
	case reflect.Invalid:
		return value.UnknownKind
	default:
		if rk < reflect.Array {
			return value.NumberKind
		}
	}
	return value.UnknownKind
}

func (o *reflectedSimpleImpl) Interface() interface{} {
	return o.rValue.Interface()
}

func (o *reflectedSimpleImpl) String() string {
	return fmt.Sprintf("%v", o.rValue.Interface())
}

func (o *reflectedSimpleImpl) SetValue(value value.Value) error {
	return fmt.Errorf("not implemented - reflectedSimpleImpl.SetValue")
}

func (o *reflectedSimpleImpl) WithoutSource() interface{} {
	return o.Interface()
}

func (o *reflectedSimpleImpl) CompareTo(other value.Simple) (int, error) {
	return 0, fmt.Errorf("not implemented - reflectedSimpleImpl.CompareTo")
}
