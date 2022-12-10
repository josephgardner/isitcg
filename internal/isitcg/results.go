package isitcg

type Results struct {
	Hash         string
	ProductName  string
	SearchResult string
	Result       string
	Remainder    []string
	Matches      []Rule
}

func NewMatchResults(remainder []string) *Results {
	return &Results{
		Remainder: remainder,
		Result:    "good",
	}
}
