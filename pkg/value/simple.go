package value

type Simple interface {
	Value
	String() string
	CompareTo(other Simple) (int, error)
}

type String interface {
	Simple
	StringOrError() (string, error)
}

type Bool interface {
	Simple
	Bool() (bool, error)
}

type Number interface {
	Simple
	Int(bits int) (int64, error)
	Float(bits int) (float64, error)
	Unsigned(bits int) (uint64, error)
	Complex(bits int) (complex128, error)
}
