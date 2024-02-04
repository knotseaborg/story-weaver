package search

// Keyword search

type search struct {
	Search []item `json:"search"`
}

type item struct {
	ID          string
	Label       string `json:"label"`
	Description string `json:"description"`
}

// For sparQL queries

type queryResult struct {
	Results result `json:"results"`
}

type result struct {
	Bindings []map[string]struct {
		XMLLang string `json:"xml:lang"`
		Type    string `json:"type"`
		Value   string `json:"value"`
	} `json:"bindings"`
}
