package dspath

import (
	"fmt"

	"github.com/davidjspooner/dsvalue/pkg/key"
	"github.com/davidjspooner/dsvalue/pkg/value"
)

type comparison struct {
	parent               *comparison
	pair                 pair
	result               int
	comparisonFilterFunc ComparisonFilterFunc
}

var _ value.Array = &comparison{}
var _ value.Map = &comparison{}

type ComparisonFilterFunc func(path Path, left, right value.Value, r int, err error) (int, error)

func (c *comparison) Kind() value.Kind {

	lKind, rKind := c.pair.Kinds()
	if lKind != rKind {
		return value.UnknownKind
	}
	return lKind
}
func (c *comparison) Source() value.Source {
	return c.pair.Source()
}
func (c *comparison) WithoutSource() interface{} {
	return c.pair.WithoutSource()
}

func (c *comparison) Field(p key.Interface) (value.Value, error) {
	child := &comparison{
		parent:               c,
		comparisonFilterFunc: c.comparisonFilterFunc,
	}
	err := c.pair.Field(p, &child.pair)
	return child, err

}

func (c *comparison) ForEach(f func(index key.Interface, value value.Value) error) error {
	return fmt.Errorf("not implemented - comparison.ForEach")
}

func (c *comparison) Length() (int, error) {
	return 0, fmt.Errorf("not implemented - comparison.Length")
}

func (c *comparison) Index(index key.Interface) (value.Value, error) {
	child := &comparison{
		parent:               c,
		comparisonFilterFunc: c.comparisonFilterFunc,
	}
	err := c.pair.Index(index, &child.pair)
	return child, err
}

func (c *comparison) compare(p Path, vt VisitType) (err error) {
	if c.comparisonFilterFunc != nil {
		defer func() {
			if vt == AtCollectionEnd || vt == AtLeaf || (err != nil && err != ErrSkipContents) {
				c.result, err = c.comparisonFilterFunc(p, c.pair.left, c.pair.right, c.result, err)
			}
		}()
	}
	var leftKind, rightKind value.Kind
	if c.pair.left != nil {
		leftKind = c.pair.left.Kind()
	}
	if c.pair.right != nil {
		rightKind = c.pair.right.Kind()
	}
	if leftKind != rightKind {
		c.result = int(leftKind - rightKind)
		return nil
	}

	comparisonFunc, ok := comparisonFunc[leftKind]
	if !ok {
		return fmt.Errorf("comparison not implemented for %s", leftKind)
	}

	c.result, err = comparisonFunc(c.pair.left, c.pair.right)
	return err
}

func (c *comparison) visitFn(p Path, v value.Value, vt VisitType) error { //TODO consider usefullness of err paramater
	pair, ok := v.(*comparison)
	if !ok {
		return fmt.Errorf("expected comparison value, got %T", v)
	}
	if pair == nil {
		return nil
	}
	err := c.compare(p, vt)
	if c.comparisonFilterFunc != nil {
		if vt == AtCollectionEnd || vt == AtLeaf {
			c.result, err = c.comparisonFilterFunc(p, pair.pair.left, pair.pair.right, c.result, err)
		}
	}
	return err
}

func Compare(left, right value.Value, comparisonFilterFunc ComparisonFilterFunc) (int, error) {
	c := &comparison{
		pair: pair{
			left:  left,
			right: right,
		},
		comparisonFilterFunc: comparisonFilterFunc,
	}
	err := Walk(c, c.visitFn)
	return c.result, err
}
