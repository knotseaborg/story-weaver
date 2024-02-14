package search

// search represents the structure of search results returned from keyword searches.
type search struct {
	Search []item `json:"search"` // Array of items found in the search
}

// item represents an individual search result item.
type item struct {
	ID          string // Unique identifier of the item
	Label       string `json:"label"`       // Label of the item
	Description string `json:"description"` // Description of the item
}

// queryResult represents the structure of SPARQL query results.
type queryResult struct {
	Results result `json:"results"` // Results of the SPARQL query
}

// result represents the structure of individual query results.
type result struct {
	Bindings []map[string]struct {
		XMLLang string `json:"xml:lang"` // Language of the result value
		Type    string `json:"type"`     // Type of the result value
		Value   string `json:"value"`    // Value of the result
	} `json:"bindings"` // Array of result bindings
}
