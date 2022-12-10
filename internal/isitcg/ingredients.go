package isitcg

import (
	"encoding/json"
	"regexp"
	"sort"
	"strings"
)

type ingredientHandler interface {
	CreateHash(product string, ingredients string) string
	ResultsFromHash(hash string) results
	ResultsFromProduct(product Product) results
	ProductFromHash(hash string) Product
}

type defaultIngredientHandler struct {
	rules []Rule
}

var _ ingredientHandler = (*defaultIngredientHandler)(nil)

func NewIngredientHandler(rules []Rule) ingredientHandler {
	return &defaultIngredientHandler{
		rules: rules,
	}
}

func (h defaultIngredientHandler) CreateHash(productName, ingredients string) string {
	product := Product{Name: productName, Ingredients: ingredients}
	json, _ := json.Marshal(product)
	hash := Compress(string(json))
	return hash
}

func (h defaultIngredientHandler) ProductFromHash(hash string) Product {
	j := Decompress(hash)
	var product Product
	json.Unmarshal([]byte(j), &product)
	return product
}

func (h defaultIngredientHandler) ResultsFromHash(hash string) results {
	product := h.ProductFromHash(hash)
	results := h.ResultsFromProduct(product)
	results.Hash = hash
	return results
}

func matchAny(str string, matchWords []string) bool {
	// Compile the regular expression to use for matching
	re := regexp.MustCompile("(?i)(" + strings.Join(matchWords, "|") + ")")

	// Find all matches in the input string
	return re.MatchString(str)
}

func (h defaultIngredientHandler) ResultsFromProduct(product Product) results {
	results := results{
		ProductName: product.Name,
		Remainder:   product.Parts(),
		Result:      "good",
	}

	for _, cur := range h.rules {
		matches := intersect(results.Remainder, cur.Ingredients)

		if len(matches) == 0 {
			continue
		}

		newRule := Rule{
			Name:        cur.Name,
			Description: cur.Description,
			BlogUrl:     cur.BlogUrl,
			Result:      cur.Result,
			Rank:        cur.Rank,
			Ingredients: matches,
		}
		results.Matches = append(results.Matches, newRule)
		if cur.Result == "danger" {
			results.Result = "danger"
		} else if cur.Result == "warning" && results.Result == "good" {
			results.Result = "warning"
		}

		newRemainder := []string{}
		for _, r := range results.Remainder {
			if !contains(matches, r) {
				newRemainder = append(newRemainder, r)
			}
		}
		results.Remainder = newRemainder
	}

	sort.Slice(results.Matches, func(i, j int) bool {
		return results.Matches[i].Rank < results.Matches[j].Rank
	})

	return results
}

func intersect(left, right []string) []string {
	matches := []string{}
	for _, l := range left {
		if matchAny(l, right) {
			matches = append(matches, l)
		}
	}
	return matches
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
