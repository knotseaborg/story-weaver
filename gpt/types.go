package gpt

// FOr Dall-e

type image struct {
	Created int    `json:"created"`
	Data    []data `json:"data"`
}

type data struct {
	RevisedPrompt string `json:"revised_prompt"`
	URL           string `json:"url"`
}

// For GPT-3.5

type completion struct {
	Choices []choice `json:"choices"`
}

type choice struct {
	Message message `json:"message"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// For query plan

type QueryPlan struct {
	Query      string       `json:"query"`
	Searchable []Searchable `json:"searchable"`
}

type Searchable struct {
	Name        string
	Keyword     string
	Description string
}
