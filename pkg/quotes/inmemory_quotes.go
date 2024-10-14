package quotes

type InMemoryQuotes struct {
	quotes []string
}

func NewInMemoryQuotes() *InMemoryQuotes {
	return &InMemoryQuotes{
		quotes: []string{
			"All saints who remember to keep and do these sayings, walking in obedience to the commandments, shall receive health in their navel and marrow to their bones",

			"And shall find wisdom and great treasures of knowledge, even hidden treasures",

			"And shall run and not be weary, and shall walk and not faint",

			"And I, the Lord, give unto them a promise, that the destroying angel shall pass by them, as the children of Israel, and not slay them",
		},
	}
}

func (q *InMemoryQuotes) Get(i int) []byte {
	return []byte(q.quotes[i])
}

func (q *InMemoryQuotes) Len() int {
	return len(q.quotes)
}
