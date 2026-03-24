package functions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// --- Response structs ---

type Thumbnail struct {
	Src    string `json:"src"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type MetaURL struct {
	Scheme   string `json:"scheme"`
	Netloc   string `json:"netloc"`
	Hostname string `json:"hostname"`
	Favicon  string `json:"favicon"`
	Path     string `json:"path"`
}

type NewsArticle struct {
	Title         string    `json:"title"`
	URL           string    `json:"url"`
	Description   string    `json:"description"`
	Age           string    `json:"age"`
	PageAge       string    `json:"page_age"`
	PageFetched   string    `json:"page_fetched"`
	Source        string    `json:"source"`
	BreakingNews  bool      `json:"breaking_news"`
	Thumbnail     Thumbnail `json:"thumbnail"`
	MetaURL       MetaURL   `json:"meta_url"`
	ExtraSnippets []string  `json:"extra_snippets,omitempty"`
}

type NewsQuery struct {
	Original             string `json:"original"`
	ShowStrictWarning    bool   `json:"show_strict_warning"`
	IsNavigational       bool   `json:"is_navigational"`
	IsNewsBreaking       bool   `json:"is_news_breaking"`
	SpellcheckOff        bool   `json:"spellcheck_off"`
	Country              string `json:"country"`
	BadResults           bool   `json:"bad_results"`
	ShouldFallback       bool   `json:"should_fallback"`
	PostalCode           string `json:"postal_code"`
	City                 string `json:"city"`
	HeaderCountry        string `json:"header_country"`
	MoreResultsAvailable bool   `json:"more_results_available"`
	State                string `json:"state"`
}

type NewsSearchResponse struct {
	Query     NewsQuery     `json:"query"`
	Articles  []NewsArticle `json:"results"`
	FetchedAt time.Time     `json:"fetched_at"`
}

// --- Search options ---

type SearchOptions struct {
	// Freshness: "pd" (24h), "pw" (7d), "pm" (31d), "py" (1y), or custom "YYYY-MM-DDtoYYYY-MM-DD"
	Freshness  string
	Count      int
	Offset     int
	Country    string
	SearchLang string
}

var defaultCyberKeywords = []string{
	"ransomware", "DDoS attack", "data breach", "cyberattack",
	"malware", "phishing", "zero-day", "vulnerability exploit",
}

// SearchCyberNews queries the Brave News Search API for cybersecurity news.
// If keywords is empty, a default set of cybersecurity terms is used.
func SearchCyberNews(apiKey string, keywords []string, opts *SearchOptions) (*NewsSearchResponse, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api key is required")
	}

	if len(keywords) == 0 {
		keywords = defaultCyberKeywords
	}

	if opts == nil {
		opts = &SearchOptions{
			Freshness: "pd", // last 24 hours by default
			Count:     20,
		}
	}

	query := strings.Join(keywords, " OR ")

	params := url.Values{}
	params.Set("q", query)
	params.Set("count", fmt.Sprintf("%d", max(1, min(opts.Count, 50))))

	if opts.Freshness != "" {
		params.Set("freshness", opts.Freshness)
	}
	if opts.Offset > 0 {
		params.Set("offset", fmt.Sprintf("%d", opts.Offset))
	}
	if opts.Country != "" {
		params.Set("country", opts.Country)
	}
	if opts.SearchLang != "" {
		params.Set("search_lang", opts.SearchLang)
	}

	endpoint := "https://api.search.brave.com/res/v1/news/search?" + params.Encode()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("X-Subscription-Token", apiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var raw struct {
		Query   NewsQuery     `json:"query"`
		Results []NewsArticle `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &NewsSearchResponse{
		Query:     raw.Query,
		Articles:  raw.Results,
		FetchedAt: time.Now().UTC(),
	}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
