package models

// Story has information of the story created
type Story struct {
	ID             int    `json:"ID"`
	Title          string `json:"Title"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	StartParagraph int    `json:"-"`
	EndParagraph   int    `json:"-"`
}

// Paragraph has information of paragraph of the story
type Paragraph struct {
	ID            int      `json:"-"`
	StartSentence int      `json:"-"`
	EndSentence   int      `json:"-"`
	Sentences     []string `json:"sentences"`
}

// DetailedStory has detailed information (paragraphs) of the story
type DetailedStory struct {
	ID             int         `json:"id"`
	Title          string      `json:"title"`
	StartParagraph int         `json:"-"`
	EndParagraph   int         `json:"-"`
	CreatedAt      string      `json:"created_at"`
	UpdatedAt      string      `json:"updated_at"`
	Paragraphs     []Paragraph `json:"paragraphs"`
}

// PostResponse is response for a new word added to the story
type PostResponse struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	CurrentSentence string `json:"current_sentence"`
}
