package value

type KindClass int

const (
	UnknownClass = KindClass(0)
	SimpleClass  = KindClass(iota+1) << collectionShift
	CollectionClass

	collectionShift = 4
	collectionMask  = -1 << collectionShift
)

func (class KindClass) String() string {
	switch class {
	case UnknownClass:
		return "Unknown"
	case SimpleClass:
		return "Simple"
	case CollectionClass:
		return "Collection"
	default:
		return "Unknown"
	}
}

type Kind int

const (
	UnknownKind = Kind(iota) + Kind(UnknownClass)
)
const (
	NullKind = Kind(iota + SimpleClass)
	StringKind
	BoolKind
	NumberKind
)
const (
	ArrayKind = Kind(iota + CollectionClass)
	MapKind
)

func (kind Kind) String() string {
	switch kind {
	case UnknownKind:
		return "Unknown"
	case NullKind:
		return "Null"
	case StringKind:
		return "String"
	case BoolKind:
		return "Bool"
	case NumberKind:
		return "Number"
	case ArrayKind:
		return "Array"
	case MapKind:
		return "Map"
	default:
		return "Unknown"
	}
}

func (kind Kind) Class() Kind {
	return kind & collectionMask
}
