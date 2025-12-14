package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/josephgardner/isitcg/internal/isitcg"
)

const (
	TMPL_INDEX    = "index"
	TMPL_RESULTS  = "results"
	TMPL_TRENDING = "trending"
)

type renders interface {
	Index(w http.ResponseWriter, p isitcg.Product)
	Results(w http.ResponseWriter, r isitcg.Results)
	Trending(w http.ResponseWriter)
}

type rendersHtml struct {
	views        map[string]*template.Template
	analyticsURL string
}

func (r *rendersHtml) Index(w http.ResponseWriter, p isitcg.Product) {
	r.render(TMPL_INDEX, w, struct {
		isitcg.Product
		AnalyticsURL string
	}{p, r.analyticsURL})
}

func (r *rendersHtml) Results(w http.ResponseWriter, res isitcg.Results) {
	r.render(TMPL_RESULTS, w, struct {
		isitcg.Results
		AnalyticsURL string
	}{res, r.analyticsURL})
}

func (r *rendersHtml) Trending(w http.ResponseWriter) {
	r.render(TMPL_TRENDING, w, map[string]string{
		"AnalyticsURL": r.analyticsURL,
	})
}

func (r *rendersHtml) render(name string, w http.ResponseWriter, data any) {
	if v, ok := r.views[name]; !ok {
		http.Error(w, fmt.Sprintf("View not found: %v", name), http.StatusInternalServerError)
	} else {
		v.Execute(w, data)
	}
}

var _ renders = (*rendersHtml)(nil)

func renderer(analyticsURL string) renders {
	return &rendersHtml{
		views: map[string]*template.Template{
			TMPL_INDEX:    loadTemplate(TMPL_INDEX),
			TMPL_RESULTS:  loadTemplate(TMPL_RESULTS),
			TMPL_TRENDING: loadTemplate(TMPL_TRENDING),
		},
		analyticsURL: analyticsURL,
	}
}

func loadTemplate(name string) *template.Template {
	funcMap := template.FuncMap{
		"div":  func(a, b int) int { return a / b },
		"mult": func(a, b int) int { return a * b },
	}
	return template.Must(template.New("base.html").Funcs(funcMap).ParseFiles(
		"./templates/base.html",
		fmt.Sprintf("./templates/%s.html", name),
	))
}
