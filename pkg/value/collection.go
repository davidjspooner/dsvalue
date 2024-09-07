package value

import "github.com/davidjspooner/dsvalue/pkg/key"

type Collection interface {
	Value
	ForEach(f func(index key.Interface, value Value) error) error
	Length() (int, error)
}

type Array interface {
	Collection
	Index(index key.Interface) (Value, error)
}

type ModifiableArray interface {
	Array
	ModifiableValue
	SetIndex(index key.Interface, value Value) error
	Append(value Value) (key.Interface, error)
}

type Map interface {
	Collection
	Field(key key.Interface) (Value, error)
}

type ModifiableMap interface {
	Map
	ModifiableValue
	SetField(key key.Interface, value Value) error
}

//---------------------------------------------------------------------
