package isitcg

type results struct {
	Hash        string
	ProductName string
	Result      string
	Remainder   []string
	Matches     []rule
}

func NewMatchResults(remainder []string) *results {
	return &results{
		Remainder: remainder,
		Result:    "good",
	}
}
