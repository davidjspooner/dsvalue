package dspath

import (
	"fmt"

	"github.com/davidjspooner/dsvalue/pkg/key"
	"github.com/davidjspooner/dsvalue/pkg/value"
)

type VisitType int

const (
	AtCollectionStart VisitType = iota
	AtCollectionEnd
	AtLeaf
)

type WalkError string

func (e WalkError) Error() string {
	return string(e)
}

const (
	ErrSkipContents   WalkError = "skip contents"
	ErrSkipRestOfWalk WalkError = "skip rest of walk"
)

func walk(aValue value.Value, path Path, visitFn func(p Path, v value.Value, vt VisitType) error) (err error) {
	pathLen := len(path)
	kind := aValue.Kind()
	switch kind {
	case value.MapKind, value.ArrayKind:
		err = visitFn(path, aValue, AtCollectionStart)
		defer func() {
			r := recover()
			if r != nil {
				err = fmt.Errorf("panic: %v", r)
			}
			path = path[:pathLen]
			err2 := visitFn(path, aValue, AtCollectionEnd)
			if err == nil {
				err = err2
			}
		}()
		if err == ErrSkipContents {
			return nil
		}
		if aValue.Kind() == value.MapKind {
			//visit map children
			m := aValue.(value.Map)
			path = append(path[:pathLen], nil)
			err = m.ForEach(func(k key.Interface, child value.Value) error {
				path[pathLen] = k
				err = walk(child, path, visitFn)
				return err
			})
		} else {
			//visit array children
			array := aValue.(value.Array)
			path = append(path[:pathLen], key.Value[int]{})
			err = array.ForEach(func(index key.Interface, child value.Value) error {
				path[pathLen] = index
				err = walk(child, path, visitFn)
				return err
			})
		}

	default:
		err = visitFn(path, aValue, AtLeaf)
	}
	return err
}

func Walk(value value.Value, f func(p Path, v value.Value, vt VisitType) error) error {
	path := Path{}
	err := walk(value, path, f)
	switch err {
	case ErrSkipContents, ErrSkipRestOfWalk:
		return nil
	default:
		return err
	}
}
