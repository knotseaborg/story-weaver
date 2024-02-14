package gpt

// image represents the response structure for DALL-E image generation.
type image struct {
	Created int    `json:"created"` // Timestamp indicating when the image was created
	Data    []data `json:"data"`    // Array of data containing revised prompt and image URL
}

// data contains information about the revised prompt and the URL of the generated image.
type data struct {
	RevisedPrompt string `json:"revised_prompt"` // Revised prompt used for image generation
	URL           string `json:"url"`            // URL of the generated image
}

// completion represents the response structure for GPT completions.
type completion struct {
	Choices []choice `json:"choices"` // Array of choices containing completion messages
}

// choice represents a single choice containing a completion message.
type choice struct {
	Message message `json:"message"` // Message object representing the completion
}

// message represents a completion message with its role and content.
type message struct {
	Role    string `json:"role"`    // Role of the message (e.g., user, model)
	Content string `json:"content"` // Content of the message
}

// QueryPlan represents a query plan structure for executing queries.
type QueryPlan struct {
	Query      string       `json:"query"`      // The query string
	Searchable []Searchable `json:"searchable"` // List of searchable entities
}

// Searchable represents an entity that can be searched within a query plan.
type Searchable struct {
	Name        string // Name of the searchable entity
	Keyword     string // Keyword associated with the entity
	Description string // Description of the searchable entity
}
