package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/josephgardner/isitcg/internal/isitcg"
)

// Cache for trending products
type trendingCache struct {
	mu      sync.RWMutex
	entries map[string]cacheEntry
}

type cacheEntry struct {
	products  []TrendingProduct
	expiresAt time.Time
}

var productCache = &trendingCache{
	entries: make(map[string]cacheEntry),
}

const cacheTTL = 5 * time.Minute

// TrendingProduct represents a product from the analytics API
type TrendingProduct struct {
	AmazonURL   string   `json:"amazon_url"`
	AmazonImage string   `json:"amazon_image"`
	AmazonTitle string   `json:"amazon_title"`
	Brand       string   `json:"brand"`
	Badges      []string `json:"badges"`
	StarRating  float64  `json:"star_rating"`
	ReviewCount int      `json:"review_count"`
}

// TrendingResponse is the API response structure
type TrendingResponse struct {
	Products []TrendingProduct `json:"products"`
}

// TrendingData is passed to the trending template
type TrendingData struct {
	AnalyticsURL string
	Products     []TrendingProduct
	SearchQuery  string
	Error        bool
}

const (
	TMPL_INDEX    = "index"
	TMPL_RESULTS  = "results"
	TMPL_TRENDING = "trending"
)

type renders interface {
	Index(w http.ResponseWriter, p isitcg.Product)
	Results(w http.ResponseWriter, r isitcg.Results)
	Trending(w http.ResponseWriter, searchQuery string, noCache bool)
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

func (r *rendersHtml) Trending(w http.ResponseWriter, searchQuery string, noCache bool) {
	data := TrendingData{
		AnalyticsURL: r.analyticsURL,
		SearchQuery:  searchQuery,
	}

	// Fetch products from analytics API (with caching)
	products, err := r.fetchTrendingProducts(searchQuery, noCache)
	if err != nil {
		data.Error = true
	} else {
		data.Products = products
	}

	r.render(TMPL_TRENDING, w, data)
}

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func (r *rendersHtml) fetchTrendingProducts(searchQuery string, noCache bool) ([]TrendingProduct, error) {
	// Only cache the default (no query) response to avoid unbounded cache growth
	useCache := searchQuery == "" && !noCache
	cacheKey := "trending:default"

	// Check cache first
	if useCache {
		if products, ok := productCache.get(cacheKey); ok {
			return products, nil
		}
	}

	// Fetch from API
	apiURL := fmt.Sprintf("%s/api/products/popular?limit=40&source=trending", r.analyticsURL)
	if searchQuery != "" {
		apiURL += "&q=" + url.QueryEscape(searchQuery)
	}

	resp, err := httpClient.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var trendingResp TrendingResponse
	if err := json.NewDecoder(resp.Body).Decode(&trendingResp); err != nil {
		return nil, err
	}

	// Only cache the default response
	if searchQuery == "" {
		productCache.set(cacheKey, trendingResp.Products)
	}

	return trendingResp.Products, nil
}

func (c *trendingCache) get(key string) ([]TrendingProduct, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[key]
	if !ok || time.Now().After(entry.expiresAt) {
		return nil, false
	}
	return entry.products, true
}

func (c *trendingCache) set(key string, products []TrendingProduct) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = cacheEntry{
		products:  products,
		expiresAt: time.Now().Add(cacheTTL),
	}
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
		"formatRating": func(rating float64) string {
			return fmt.Sprintf("%.1f", rating)
		},
		"fullStars": func(rating float64) int {
			return int(math.Floor(rating))
		},
		"hasHalfStar": func(rating float64) bool {
			return math.Mod(rating, 1) >= 0.3
		},
		"emptyStars": func(rating float64) int {
			full := int(math.Floor(rating))
			hasHalf := 0
			if math.Mod(rating, 1) >= 0.3 {
				hasHalf = 1
			}
			return 5 - full - hasHalf
		},
		"formatCount": func(count int) string {
			if count >= 1000 {
				return fmt.Sprintf("%.1fK", float64(count)/1000)
			}
			return fmt.Sprintf("%d", count)
		},
		"hasBadge": func(badges []string, badge string) bool {
			for _, b := range badges {
				if b == badge {
					return true
				}
			}
			return false
		},
		"seq": func(n int) []int {
			s := make([]int, n)
			for i := range s {
				s[i] = i
			}
			return s
		},
		"productsWithRatings": func(products []TrendingProduct) []TrendingProduct {
			var filtered []TrendingProduct
			for _, p := range products {
				if p.StarRating > 0 {
					filtered = append(filtered, p)
				}
			}
			return filtered
		},
	}
	return template.Must(template.New("base.html").Funcs(funcMap).ParseFiles(
		"./templates/base.html",
		fmt.Sprintf("./templates/%s.html", name),
	))
}
