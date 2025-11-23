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
	current := strings.Builder{}
	depth := 0

	for _, ch := range p.Ingredients {
		if ch == '(' {
			depth++
			current.WriteRune(ch)
		} else if ch == ')' {
			depth--
			current.WriteRune(ch)
		} else if ch == ',' && depth == 0 {
			part := strings.TrimSpace(current.String())
			part = strings.Trim(part, ".")
			if len(part) > 0 {
				parts = append(parts, part)
			}
			current.Reset()
		} else {
			current.WriteRune(ch)
		}
	}

	// Don't forget the last part
	part := strings.TrimSpace(current.String())
	part = strings.Trim(part, ".")
	if len(part) > 0 {
		parts = append(parts, part)
	}

	return parts
}
