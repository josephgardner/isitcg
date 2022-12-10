package isitcg

import "strings"

type Product struct {
	Name        string
	Ingredients string
}

func NewProduct(name, ingredients string) *Product {
	return &Product{
		Name:        name,
		Ingredients: ingredients,
	}
}

func (p *Product) Parts() []string {
	parts := strings.Split(p.Ingredients, ",")
	for i, part := range parts {
		part = strings.TrimSpace(part)
		part = strings.Trim(part, ".")
		parts[i] = part
	}
	return parts
}
