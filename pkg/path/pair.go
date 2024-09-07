package dspath

import (
	"fmt"

	"github.com/davidjspooner/dsvalue/pkg/key"
	"github.com/davidjspooner/dsvalue/pkg/value"
	"golang.org/x/exp/constraints"
)

type pair struct {
	left  value.Value
	right value.Value
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

var comparisonFunc = map[value.Kind]func(left, right value.Value) (int, error){
	value.NullKind: func(left, right value.Value) (int, error) {
		return 0, nil
	},
	value.BoolKind: func(left, right value.Value) (int, error) {
		l := left.(value.Bool)
		r := right.(value.Bool)
		return l.CompareTo(r)
	},
	value.NumberKind: func(left, right value.Value) (int, error) {
		l := left.(value.Number)
		r := right.(value.Number)
		return l.CompareTo(r)
	},
	value.StringKind: func(left, right value.Value) (int, error) {
		l := left.(value.String)
		r := right.(value.String)
		return l.CompareTo(r)
	},
	value.ArrayKind: func(left, right value.Value) (int, error) {
		return 0, fmt.Errorf("not implemented - pairFunc[ArrayKind]")
	},
	value.MapKind: func(left, right value.Value) (int, error) {
		return 0, fmt.Errorf("not implemented - pairFunc[Map]")
	},
}

func (p *pair) Kind() value.Kind {
	leftKind, _ := p.Kinds()
	return leftKind
}

func (p *pair) Kinds() (value.Kind, value.Kind) {
	var lKind, rKind value.Kind
	if p.left != nil {
		lKind = p.left.Kind()
	}
	if p.right != nil {
		rKind = p.right.Kind()
	}
	return lKind, rKind
}
func (p *pair) Source() value.Source {
	var lSource, rSource value.Source
	if p.left != nil {
		lSource = p.left.Source()
	}
	if p.right != nil {
		rSource = p.right.Source()
	}
	if lSource == nil {
		return rSource
	}
	return lSource
}
func (p *pair) WithoutSource() interface{} {
	var pair [2]any
	if p.left != nil {
		pair[0] = p.left.WithoutSource()
	}
	if p.right != nil {
		pair[1] = p.right.WithoutSource()
	}

	return pair
}

func (p *pair) Field(k key.Interface, child *pair) error {
	var leftMap, rightMap value.Map
	if p.left != nil {
		leftMap = p.left.(value.Map)
	}
	if p.right != nil {
		rightMap = p.right.(value.Map)
	}
	if leftMap != nil {
		child.left, _ = leftMap.Field(k)
	}
	if rightMap != nil {
		child.right, _ = rightMap.Field(k)
	}
	return nil
}

func (p *pair) Length() (int, error) {
	return 0, fmt.Errorf("not implemented - pair.Length")
}

func (p *pair) Index(index key.Interface, child *pair) error {
	var leftArray, rightArray value.Array
	if p.left != nil {
		leftArray = p.left.(value.Array)
	}
	if p.right != nil {
		rightArray = p.right.(value.Array)
	}
	if leftArray != nil {
		child.left, _ = leftArray.Index(index)
	}
	if rightArray != nil {
		child.right, _ = rightArray.Index(index)
	}
	return nil
}
