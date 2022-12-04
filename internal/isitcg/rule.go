package isitcg

type rule struct {
	Name        string
	Description string
	Result      string `yaml:"Result"`
	BlogUrl     string
	Rank        int      `yaml:"Rank"`
	Ingredients []string `yaml:"Ingredients"`
}
