package isitcg

import (
	"encoding/json"
	"io/ioutil"
)

type rules struct {
	Rules []Rule
}

type Rule struct {
	Name        string
	Description string
	Result      string `yaml:"Result"`
	BlogUrl     string
	Rank        int      `yaml:"Rank"`
	Ingredients []string `yaml:"Ingredients"`
}

func LoadRules(filename string) ([]Rule, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var rules rules
	err = json.Unmarshal(b, &rules)
	if err != nil {
		return nil, err
	}
	return rules.Rules, err
}
