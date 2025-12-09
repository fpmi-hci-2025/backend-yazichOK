package entity

type Topic struct {
	ID          int
	Title       string
	Description string
	PhotoURL    string
	Questions   []Question
}

type Question struct {
	ID   int
	Text string
}

type Session struct {
	ID        string
	TopicID   int
	Questions []Question
}

type TopWord struct {
	Words string `json:"words"`
	Level string `json:"level"`
}

type GrammarIssue struct {
	Sentence          string `json:"sentence"`
	Explanation       string `json:"explanation"`
	CorrectedSentence string `json:"corrected_sentence"`
}

type RephraseSuggestion struct {
	Original   string `json:"original"`
	Suggestion string `json:"suggestion"`
}

type AnalyzeTextResult struct {
	OverallLevel        string               `json:"overall_level"`
	TopWords            []TopWord            `json:"top_words"`
	GrammarIssues       []GrammarIssue       `json:"grammar_issues"`
	RephraseSuggestions []RephraseSuggestion `json:"rephrase_suggestions"`
	OverallFeedback     string               `json:"overall_feedback"`
}

type ArticlePreview struct {
	ID            int
	ImageURL      string
	Level         string
	MinutesToRead int
	Title         string
}

type Article struct {
	ID         int
	ImageURL   string
	Content    string
	Title      string
	Level      string
	Minutes    int
	Vocabulary []VocabularyWord
	Rules      []GrammarRule
}

type VocabularyWord struct {
	ID           int
	Word         string
	PartOfSpeech string
	Meaning      string
}

type GrammarRule struct {
	ID      int
	Name    string
	Example string
	Note    string
}

type WordCollection struct {
	ID                string // UUID
	UserID            int
	Name              string
	ImageURL          string
	TotalWordsCount   int
	LearnedWordsCount int
	CurrentStreakDays int
	LongestStreakDays int
	LastStudiedAt     *string
	AISuggestions     []AISuggestion
	CreatedAt         string
	UpdatedAt         string
}

type AISuggestion struct {
	Word        string `json:"word"`
	Translation string `json:"translation"`
	Reason      string `json:"reason"`
}

type UserWord struct {
	ID             string
	CollectionID   string
	Word           string
	Translation    string
	Example        *string
	NextReviewDate string
	ReviewCount    int
	CreatedAt      string
	UpdatedAt      string
}

type WordCollectionDetail struct {
	ID                string
	Name              string
	ImageURL          string
	TotalWordsCount   int
	LearnedWordsCount int
	CurrentStreakDays int
	LongestStreakDays int
	LastStudiedAt     *string
	UserWords         []UserWord
	AISuggestions     []AISuggestion
	CreatedAt         string
	UpdatedAt         string
}
