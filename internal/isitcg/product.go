package isitcg

import "strings"

type product struct {
	Name        string
	Ingredients string
}

func NewProduct(name, ingredients string) *product {
	return &product{
		Name:        name,
		Ingredients: ingredients,
	}
}

func (p *product) Parts() []string {
	parts := strings.Split(p.Ingredients, ",")
	for i, part := range parts {
		part = strings.TrimSpace(part)
		part = strings.Trim(part, ".")
		parts[i] = part
	}
	return parts
}
