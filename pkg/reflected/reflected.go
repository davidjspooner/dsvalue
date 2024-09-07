package reflected

import (
	"fmt"
	"reflect"

	"github.com/davidjspooner/dsvalue/pkg/value"
)

type Reflected interface {
	value.Value
}

func NewReflectedObject(rValue reflect.Value, source value.Source) (value.Value, error) {
	rk := rValue.Kind()
	object, ok := rValue.Interface().(value.Value)
	if ok {
		return object, nil
	}
	for rk == reflect.Ptr || rk == reflect.Interface {
		if rValue.IsNil() {
			return value.NewNull(source), nil
		}
		rValue = rValue.Elem()
		rk = rValue.Kind()
		object, ok := rValue.Interface().(value.Value)
		if ok {
			return object, nil
		}
	}
	switch rk {
	case reflect.Array, reflect.Slice:
		return &reflectedArrayImpl{
			rValue: rValue,
			source: source,
		}, nil
	case reflect.Map:
		return &reflectedMapImpl{
			rValue: rValue,
			source: source,
		}, nil
	case reflect.String:
		return &reflectedSimpleImpl{
			rValue: rValue,
			source: source,
		}, nil
	case reflect.Struct:
		return &reflectedStructImpl{
			rValue: rValue,
			source: source,
		}, nil
	default:
		if rk < reflect.Array {
			return &reflectedSimpleImpl{
				rValue: rValue,
				source: source,
			}, nil
		}
	}
	return nil, fmt.Errorf("unsupported kind: %s", rk)
}
