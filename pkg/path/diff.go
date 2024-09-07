package path

import (
	"fmt"

	"github.com/davidjspooner/dsvalue/pkg/key"
	"github.com/davidjspooner/dsvalue/pkg/value"
)

type diff struct {
	pair                  pair
	differenceHandlerFunc func(p Path, left, right value.Value) error
}

var _ value.Array = &comparison{}
var _ value.Map = &comparison{}

func (d *diff) Kind() value.Kind {
	lKind, rKind := d.pair.Kinds()
	if lKind != rKind {
		return value.UnknownKind
	}
	return lKind
}
func (d *diff) Source() value.Source {
	return d.pair.Source()
}
func (d *diff) WithoutSource() interface{} {
	return d.pair.WithoutSource()
}
func (d *diff) Field(p key.Interface) (value.Value, error) {
	child := &diff{
		differenceHandlerFunc: d.differenceHandlerFunc,
	}
	err := d.pair.Field(p, &child.pair)
	return child, err
}

func (d *diff) forEachArray(f func(index key.Interface, value value.Value) error) error {
	leftArray, ok := d.pair.left.(value.Array)
	if !ok {
		return fmt.Errorf("left is not an array")
	}
	rightArray, ok := d.pair.right.(value.Array)
	if !ok {
		return fmt.Errorf("right is not an array")
	}
	leftLength, _ := leftArray.Length()
	rightLength, _ := rightArray.Length()
	count := Max(leftLength, rightLength)
	i := key.Value[int]{}
	for i.X = 0; i.X < count; i.X++ {
		child := &diff{
			differenceHandlerFunc: d.differenceHandlerFunc,
		}
		if i.X < leftLength {
			child.pair.left, _ = leftArray.Index(i)
		}
		if i.X < rightLength {
			child.pair.right, _ = rightArray.Index(i)
		}
		err := f(i, child)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *diff) forEachMap(f func(index key.Interface, value value.Value) error) error {
	leftMap, ok := d.pair.left.(value.Map)
	if !ok {
		return fmt.Errorf("left is not a map")
	}
	rightMap, ok := d.pair.right.(value.Map)
	if !ok {
		return fmt.Errorf("right is not a map")
	}

	//var leftKeys, rightKeys []string
	//leftMap.ForEach(func(k key.Interface, v value.Value) error {
	//	leftKeys = append(leftKeys, k.String())
	//	return nil
	//})
	//rightMap.ForEach(func(k key.Interface, v value.Value) error {
	//	rightKeys = append(rightKeys, k.String())
	//	return nil
	//})

	err := leftMap.ForEach(func(k key.Interface, v value.Value) error {
		child := &diff{
			differenceHandlerFunc: d.differenceHandlerFunc,
		}

		//TODO remove this
		leftMapReal := leftMap.WithoutSource()
		rightMapReal := rightMap.WithoutSource()
		_, _ = leftMapReal, rightMapReal

		child.pair.left = v
		child.pair.right, _ = rightMap.Field(k)
		return f(k, child)
	})
	if err != nil {
		return err
	}
	err = rightMap.ForEach(func(k key.Interface, v value.Value) error {
		child := &diff{
			differenceHandlerFunc: d.differenceHandlerFunc,
		}

		//TODO remove this
		leftMapReal := leftMap.WithoutSource()
		rightMapReal := rightMap.WithoutSource()
		_, _ = leftMapReal, rightMapReal

		child.pair.right = v
		child.pair.left, _ = leftMap.Field(k)
		if child.pair.left == nil {
			//we must have seen this key in the left map
			return nil
		}
		return f(k, child)
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *diff) ForEach(f func(index key.Interface, value value.Value) error) error {
	leftKind, rightKind := d.pair.Kinds()
	if leftKind != rightKind {
		d.differenceHandlerFunc(nil, d.pair.left, d.pair.right)
		return nil
	}

	//TODO remove this
	leftReal := d.pair.left.WithoutSource()
	rightReal := d.pair.right.WithoutSource()
	_, _ = leftReal, rightReal

	switch leftKind {
	case value.ArrayKind:
		return d.forEachArray(f)
	case value.MapKind:
		return d.forEachMap(f)
	default:
		return fmt.Errorf("unsupported type for iterating %q", leftKind)
	}
}
func (d *diff) Length() (int, error) {
	return 0, fmt.Errorf("not implemented - diff.Length")
}
func (d *diff) Index(index key.Interface) (value.Value, error) {
	child := &diff{
		differenceHandlerFunc: d.differenceHandlerFunc,
	}
	err := d.pair.Index(index, &child.pair)
	return child, err
}

func diffVisitFn(p Path, v value.Value, visitType VisitType) error {

	d, ok := v.(*diff)
	if !ok {
		return fmt.Errorf("expected *diff, got %T", v)
	}

	lKind, rKind := d.pair.Kinds()
	if lKind != rKind {
		return d.differenceHandlerFunc(p, d.pair.left, d.pair.right)
	}
	switch visitType {
	case AtCollectionStart:
		leftArray, ok := d.pair.left.(value.Array)
		if ok {
			rightArray, ok := d.pair.right.(value.Array)
			if ok {
				leftLength, _ := leftArray.Length()
				rightArray, _ := rightArray.Length()
				if leftLength != rightArray {
					err := d.differenceHandlerFunc(p, d.pair.left, d.pair.right)
					if err == nil {
						err = ErrSkipContents
					}
				}
			}
		}
		//pass - all good
	case AtCollectionEnd:
		//pass - all good
	case AtLeaf:
		if d.pair.left == nil || d.pair.right == nil {
			if d.pair.left != d.pair.right {
				err := d.differenceHandlerFunc(p, d.pair.left, d.pair.right)
				if err != nil {
					return err
				}
			}
		} else {
			leftString := d.pair.left.(value.Simple).String()
			rightString := d.pair.right.(value.Simple).String()
			if leftString != rightString {
				err := d.differenceHandlerFunc(p, d.pair.left, d.pair.right)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func Diff(left, right value.Value, differenceHandlerFunc func(p Path, left, right value.Value) error) error {
	d := &diff{
		pair: pair{
			left:  left,
			right: right,
		},
		differenceHandlerFunc: differenceHandlerFunc,
	}
	err := Walk(d, diffVisitFn)
	return err
}
