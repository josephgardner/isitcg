package isitcg

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// RuleTests holds a set of tests for rules
type RuleTests struct {
	Tests []RuleTest `yaml:"Tests"`
}

// RuleTest holds a single test for rules
type RuleTest struct {
	TestName          string   `yaml:"TestName"`
	Ingredients       string   `yaml:"Ingredients"`
	Rules             []Rule   `yaml:"Rules"`
	ExpectedResult    string   `yaml:"ExpectedResult"`
	ExpectedMatches   []Rule   `yaml:"ExpectedMatches"`
	ExpectedRemainder []string `yaml:"ExpectedRemainder"`
}

func (r RuleTest) String() string {
	return r.TestName
}

func NewRuleTests() *RuleTests {
	yamlFile, err := ioutil.ReadFile("testdata/rule-tests.yml")
	if err != nil {
		panic(err)
	}

	ruleTests := &RuleTests{}
	err = yaml.Unmarshal(yamlFile, ruleTests)
	if err != nil {
		panic(err)
	}

	return &RuleTests{ruleTests.Tests}
}

func TestRuleMatchesIngredient(t *testing.T) {
	for _, c := range NewRuleTests().Tests {
		t.Run(c.TestName, func(tt *testing.T) {
			// Act
			actual := NewIngredientHandler(c.Rules).ResultsFromProduct(Product{
				Name:        "test",
				Ingredients: c.Ingredients,
			})

			// Assert
			assert.NotNil(tt, actual)
			assert.Equal(tt, c.ExpectedResult, actual.Result, "expected result")
			assert.Equal(tt, c.ExpectedRemainder, actual.Remainder, "expected remainder")
			require.Equal(tt, len(c.ExpectedMatches), len(actual.Matches), "match length")
			for i := 0; i < len(c.ExpectedMatches); i++ {
				expectMatch := c.ExpectedMatches[i]
				actualMatch := actual.Matches[i]
				assert.Equal(tt, expectMatch.Result, actualMatch.Result)
				assert.Equal(tt, expectMatch.Ingredients, actualMatch.Ingredients)
			}
		})
	}
}
