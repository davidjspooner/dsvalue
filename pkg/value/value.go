package value

type Value interface {
	Kind() Kind
	Source() Source
	WithoutSource() interface{}
}

type ModifiableValue interface {
	SetValue(value Value) error
}
