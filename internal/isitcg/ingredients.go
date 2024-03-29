package isitcg

import (
	"encoding/json"
	"regexp"
	"sort"
	"strings"
)

type IngredientHandler interface {
	CreateHash(product string, ingredients string) string
	ResultsFromHash(hash string) Results
	ResultsFromProduct(product Product) Results
	ProductFromHash(hash string) Product
}

type defaultIngredientHandler struct {
	rules []Rule
}

var _ IngredientHandler = (*defaultIngredientHandler)(nil)

func NewIngredientHandler(rules []Rule) IngredientHandler {
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

func (h defaultIngredientHandler) ResultsFromHash(hash string) Results {
	product := h.ProductFromHash(hash)
	results := h.ResultsFromProduct(product)
	results.Hash = hash
	return results
}

func matchAny(str string, matchWords []string) bool {

	low := strings.ToLower(str)
	parts := strings.Split(low, "/")
	parts = append([]string{low}, parts...)

	regex := regexp.MustCompile(`(\[.*?\])|(\(.*?\))|\W`)
	for _, part := range parts {
		for _, matchWord := range matchWords {
			s1 := regex.ReplaceAllString(matchWord, "")
			s2 := regex.ReplaceAllString(part, "")
			if strings.EqualFold(strings.ToLower(s1), s2) {
				return true
			}
		}
	}
	return false
}

func (h defaultIngredientHandler) ResultsFromProduct(product Product) Results {
	results := Results{
		ProductName: product.Name,
		Matches:     []Rule{},
		Remainder:   product.Parts(),
		Result:      "good",
	}

	if len(results.ProductName) != 0 {
		results.SearchResult = results.ProductName
	} else {
		results.SearchResult = "curly girl method"
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
