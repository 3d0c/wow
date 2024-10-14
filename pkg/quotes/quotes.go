package quotes

type Quotes interface {
	Get(i int) []byte
	Len() int
}
