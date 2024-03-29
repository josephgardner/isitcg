package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/josephgardner/isitcg/internal/isitcg"
)

const (
	TMPL_INDEX   = "index"
	TMPL_RESULTS = "results"
)

type renders interface {
	Index(w http.ResponseWriter, p isitcg.Product)
	Results(w http.ResponseWriter, r isitcg.Results)
}

type rendersHtml struct {
	views map[string]*template.Template
}

func (r *rendersHtml) Index(w http.ResponseWriter, p isitcg.Product) {
	r.render(TMPL_INDEX, w, p)
}

func (r *rendersHtml) Results(w http.ResponseWriter, res isitcg.Results) {
	r.render(TMPL_RESULTS, w, res)
}

func (r *rendersHtml) render(name string, w http.ResponseWriter, data any) {
	if v, ok := r.views[name]; !ok {
		http.Error(w, fmt.Sprintf("View not found: %v", name), http.StatusInternalServerError)
	} else {
		v.Execute(w, data)
	}
}

var _ renders = (*rendersHtml)(nil)

func renderer() renders {
	return &rendersHtml{map[string]*template.Template{
		TMPL_INDEX:   loadTemplate(TMPL_INDEX),
		TMPL_RESULTS: loadTemplate(TMPL_RESULTS),
	}}
}

func loadTemplate(name string) *template.Template {
	return template.Must(template.ParseFiles(
		"./templates/base.html",
		fmt.Sprintf("./templates/%s.html", name),
	))
}
