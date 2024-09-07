package key

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Value[T any] struct{ X T }

var _ Interface = Value[int]{0}

//func (p *Value) EvaluateFor(obj Value) (Value, error) {
//	m, ok := obj.(Map)
//	if !ok {
//		a, ok := obj.(Array)
//		if ok {
//			r := NewArray(nil, a.Source())
//			a.ForEach(func(index any, value Value) error {
//				childEval, err := p.EvaluateFor(value)
//				if err != nil {
//					return err
//				}
//				r.Append(childEval)
//				return nil
//			})
//			return r, nil
//		}
//
//		return nil, fmt.Errorf("expected map, got %T", obj)
//	}
//	return m.Field(p)
//}

var identPattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

func (p Value[T]) String() string {

	v := reflect.ValueOf(p.X)
	k := v.Kind()
	if k >= reflect.Bool && k <= reflect.Complex128 {
		return fmt.Sprintf("[%v]", p.X)
	}
	switch k {
	case reflect.String:
		s := fmt.Sprint(p.X)
		if identPattern.MatchString(s) {
			return fmt.Sprintf(".%s", s)
		}
		return fmt.Sprintf("[%q]", s)
	default:
		panic(fmt.Sprintf("unsupported type - %s", v.Type().String()))
	}
}

type Range struct {
	Start, End int
	Tail       bool
}

var _ Interface = &Range{}

func (p *Range) String() string {
	sb := strings.Builder{}
	sb.WriteString("[")
	if p.Start != 0 {
		sb.WriteString(strconv.Itoa(p.Start))
	}
	sb.WriteString(":")
	if !p.Tail {
		sb.WriteString(strconv.Itoa(p.End))
	}
	sb.WriteString("]")
	return sb.String()
}
