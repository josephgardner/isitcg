package isitcg

import "strings"

type Product struct {
	Name        string `json:"n"`
	Ingredients string `json:"i"`
}

func NewProduct(name, ingredients string) *Product {
	return &Product{
		Name:        name,
		Ingredients: ingredients,
	}
}

func (p *Product) Parts() []string {
	parts := make([]string, 0)
	for _, part := range strings.Split(p.Ingredients, ",") {
		part = strings.TrimSpace(part)
		part = strings.Trim(part, ".")
		if len(part) > 0 {
			parts = append(parts, part)
		}
	}
	return parts
}
