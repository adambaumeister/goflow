package fields

//
// Value represents the interface to flowRecord Field values
// Field values can be of many types but should always implement the same methods
type Value interface {
	ToString() string
	SetType(uint16)
	ToInt() int
	ToBytes() []byte
}
