package views

import "speech-processing-service/internal/entity"

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	ErrorCode int    `json:"code"`
	Msg       string `json:"message"`
}

type GetAllTopicsResponse struct {
	Topics []struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		PhotoURL    string `json:"photo_url"`
	} `json:"topics"`
}

func NewGetAllToicsResponse(topics []entity.Topic) GetAllTopicsResponse {
	var topicsResp GetAllTopicsResponse
	for _, topic := range topics {
		topicsResp.Topics = append(topicsResp.Topics, struct {
			ID          int    `json:"id"`
			Title       string `json:"title"`
			Description string `json:"description"`
			PhotoURL    string `json:"photo_url"`
		}{
			ID:          topic.ID,
			Title:       topic.Title,
			Description: topic.Description,
			PhotoURL:    topic.PhotoURL,
		})
	}

	return topicsResp
}

type GetTopicQuestionsResponse struct {
	Questions []Question `json:"questions"`
}

type Question struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

func NewGetTopicQuestionsResponse(questions []entity.Question) GetTopicQuestionsResponse {
	var questionsResp GetTopicQuestionsResponse
	for _, question := range questions {
		questionsResp.Questions = append(questionsResp.Questions, Question{
			ID:   question.ID,
			Text: question.Text,
		})
	}

	return questionsResp
}

type StartSessionResponse struct {
	Session struct {
		ID    string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
		Topic struct {
			ID        int `json:"id" example:"1"`
			Questions []struct {
				ID   int    `json:"id" example:"101"`
				Text string `json:"text" example:"What is your name?"`
			} `json:"questions"`
		} `json:"topic"`
	} `json:"session"`
}

func NewStartSessionResponse(session *entity.Session) StartSessionResponse {
	var sessionResp StartSessionResponse

	sessionResp.Session.ID = session.ID
	sessionResp.Session.Topic.ID = session.TopicID

	for _, question := range session.Questions {
		sessionResp.Session.Topic.Questions = append(sessionResp.Session.Topic.Questions, struct {
			ID   int    `json:"id" example:"101"`
			Text string `json:"text" example:"What is your name?"`
		}{
			ID:   question.ID,
			Text: question.Text,
		})
	}

	return sessionResp
}

type CompleteSessionResp struct {
	OverallLevel string `json:"overall_level"`
	TopWords     []struct {
		Word  string `json:"word"`
		Level string `json:"level"`
	} `json:"top_words"`
	GrammarIssues []struct {
		Sentence          string `json:"sentence"`
		Explanation       string `json:"explanation"`
		CorrectedSentence string `json:"corrected_sentence"`
	} `json:"grammar_issues"`
	RephraseSuggestions []struct {
		Original   string `json:"original"`
		Suggestion string `json:"suggestion"`
	} `json:"rephrase_suggestions"`
	OverallFeedback string `json:"overall_feedback"`
}

func NewCompleteSessionResp(result *entity.AnalyzeTextResult) CompleteSessionResp {
	analyzeTextResp := CompleteSessionResp{
		TopWords: make([]struct {
			Word  string `json:"word"`
			Level string `json:"level"`
		}, 0, len(result.TopWords)),
		GrammarIssues: make([]struct {
			Sentence          string `json:"sentence"`
			Explanation       string `json:"explanation"`
			CorrectedSentence string `json:"corrected_sentence"`
		}, 0, len(result.GrammarIssues)),
		RephraseSuggestions: make([]struct {
			Original   string `json:"original"`
			Suggestion string `json:"suggestion"`
		}, 0, len(result.RephraseSuggestions)),
	}

	analyzeTextResp.OverallLevel = result.OverallLevel

	for _, word := range result.TopWords {
		analyzeTextResp.TopWords = append(analyzeTextResp.TopWords, struct {
			Word  string `json:"word"`
			Level string `json:"level"`
		}{
			Word:  word.Words,
			Level: word.Level,
		})
	}

	for _, issue := range result.GrammarIssues {
		analyzeTextResp.GrammarIssues = append(analyzeTextResp.GrammarIssues, struct {
			Sentence          string `json:"sentence"`
			Explanation       string `json:"explanation"`
			CorrectedSentence string `json:"corrected_sentence"`
		}{
			Sentence:          issue.Sentence,
			Explanation:       issue.Explanation,
			CorrectedSentence: issue.CorrectedSentence,
		})
	}

	for _, suggestion := range result.RephraseSuggestions {
		analyzeTextResp.RephraseSuggestions = append(analyzeTextResp.RephraseSuggestions, struct {
			Original   string `json:"original"`
			Suggestion string `json:"suggestion"`
		}{
			Original:   suggestion.Original,
			Suggestion: suggestion.Suggestion,
		})
	}

	analyzeTextResp.OverallFeedback = result.OverallFeedback

	return analyzeTextResp
}

type ArticlePreview struct {
	ID            int    `json:"id"`
	ImageURL      string `json:"image_url"`
	Level         string `json:"level"`
	MinutesToRead int    `json:"minutes_to_read"`
	Title         string `json:"title"`
}

type ArticlesPreviewData struct {
	Articles []ArticlePreview `json:"articles"`
}

type VocabularyWord struct {
	ID           int    `json:"id"`
	Word         string `json:"word"`
	PartOfSpeech string `json:"part_of_speech"`
	Meaning      string `json:"meaning"`
}

type GrammarRuleItem struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Example string `json:"example"`
	Note    string `json:"note"`
}

type ArticleDetails struct {
	ID         int               `json:"id"`
	ImageURL   string            `json:"image_url"`
	Content    string            `json:"content"`
	Title      string            `json:"title"`
	Level      string            `json:"level"`
	Minutes    int               `json:"minutes"`
	Vocabulary []VocabularyWord  `json:"vocabulary"`
	Rules      []GrammarRuleItem `json:"rules"`
}

type ArticleData struct {
	Article ArticleDetails `json:"article"`
}

func NewArticlesPreviewResponse(articles []entity.ArticlePreview) ArticlesPreviewData {
	previews := make([]ArticlePreview, 0, len(articles))
	for _, article := range articles {
		previews = append(previews, ArticlePreview{
			ID:            article.ID,
			ImageURL:      article.ImageURL,
			Level:         article.Level,
			MinutesToRead: article.MinutesToRead,
			Title:         article.Title,
		})
	}
	return ArticlesPreviewData{
		Articles: previews,
	}
}

func NewArticleResponse(article entity.Article) ArticleData {
	vocabulary := make([]VocabularyWord, 0, len(article.Vocabulary))
	for _, word := range article.Vocabulary {
		vocabulary = append(vocabulary, VocabularyWord{
			ID:           word.ID,
			Word:         word.Word,
			PartOfSpeech: word.PartOfSpeech,
			Meaning:      word.Meaning,
		})
	}

	rules := make([]GrammarRuleItem, 0, len(article.Rules))
	for _, rule := range article.Rules {
		rules = append(rules, GrammarRuleItem{
			ID:      rule.ID,
			Name:    rule.Name,
			Example: rule.Example,
			Note:    rule.Note,
		})
	}

	return ArticleData{
		Article: ArticleDetails{
			ID:         article.ID,
			ImageURL:   article.ImageURL,
			Content:    article.Content,
			Title:      article.Title,
			Level:      article.Level,
			Minutes:    article.Minutes,
			Vocabulary: vocabulary,
			Rules:      rules,
		},
	}
}

type WordCollectionResponse struct {
	ID                string `json:"id"` // UUID
	Name              string `json:"name"`
	ImageURL          string `json:"image_url"`
	TotalWords        int    `json:"total_words"`
	LearnedWords      int    `json:"learned_words"`
	CurrentStreakDays int    `json:"current_streak_days"`
	CreatedAt         string `json:"created_at"`
}

type CreateWordCollectionResponse struct {
	Data struct {
		Collection WordCollectionResponse `json:"collection"`
	} `json:"data"`
}

func NewCreateWordCollectionResponse(collection entity.WordCollection) CreateWordCollectionResponse {
	var resp CreateWordCollectionResponse
	resp.Data.Collection = WordCollectionResponse{
		ID:                collection.ID,
		Name:              collection.Name,
		ImageURL:          collection.ImageURL,
		TotalWords:        collection.TotalWordsCount,
		LearnedWords:      collection.LearnedWordsCount,
		CurrentStreakDays: collection.CurrentStreakDays,
		CreatedAt:         collection.CreatedAt,
	}
	return resp
}

type GetUserCollectionsResponse struct {
	Collections []WordCollectionResponse `json:"collections"`
}

func NewGetUserCollectionsResponse(collections []entity.WordCollection) GetUserCollectionsResponse {
	var resp GetUserCollectionsResponse
	resp.Collections = make([]WordCollectionResponse, 0, len(collections))

	for _, collection := range collections {
		resp.Collections = append(resp.Collections, WordCollectionResponse{
			ID:                collection.ID,
			Name:              collection.Name,
			ImageURL:          collection.ImageURL,
			TotalWords:        collection.TotalWordsCount,
			LearnedWords:      collection.LearnedWordsCount,
			CurrentStreakDays: collection.CurrentStreakDays,
			CreatedAt:         collection.CreatedAt,
		})
	}

	return resp
}

type WordCollectionDetailResponse struct {
	Collection WordCollectionDetail `json:"collection"`
}

type WordCollectionDetail struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	TotalWords    int               `json:"total_words"`
	UserWords     []UserWordDTO     `json:"user_words"`
	AISuggestions []AISuggestionDTO `json:"ai_suggestions"`
}

type UserWordDTO struct {
	ID             string  `json:"id"`
	Word           string  `json:"word"`
	Translation    string  `json:"translation"`
	Example        *string `json:"example"`
	NextReviewDate string  `json:"next_review_date"`
	ReviewCount    int     `json:"review_count"`
}

type AISuggestionDTO struct {
	Word        string `json:"word"`
	Translation string `json:"translation"`
	Reason      string `json:"reason"`
}

func NewWordCollectionDetailResponse(detail entity.WordCollectionDetail) WordCollectionDetailResponse {
	userWords := make([]UserWordDTO, 0, len(detail.UserWords))
	for _, word := range detail.UserWords {
		userWords = append(userWords, UserWordDTO{
			ID:             word.ID,
			Word:           word.Word,
			Translation:    word.Translation,
			Example:        word.Example,
			NextReviewDate: word.NextReviewDate,
			ReviewCount:    word.ReviewCount,
		})
	}

	aiSuggestions := make([]AISuggestionDTO, 0, len(detail.AISuggestions))
	for _, suggestion := range detail.AISuggestions {
		aiSuggestions = append(aiSuggestions, AISuggestionDTO{
			Word:        suggestion.Word,
			Translation: suggestion.Translation,
			Reason:      suggestion.Reason,
		})
	}

	return WordCollectionDetailResponse{
		Collection: WordCollectionDetail{
			ID:            detail.ID,
			Name:          detail.Name,
			TotalWords:    detail.TotalWordsCount,
			UserWords:     userWords,
			AISuggestions: aiSuggestions,
		},
	}
}

type AddWordToCollectionResponse struct {
	Word UserWordDTO `json:"word"`
}

func NewAddWordToCollectionResponse(word entity.UserWord) AddWordToCollectionResponse {
	return AddWordToCollectionResponse{
		Word: UserWordDTO{
			ID:             word.ID,
			Word:           word.Word,
			Translation:    word.Translation,
			Example:        word.Example,
			NextReviewDate: word.NextReviewDate,
			ReviewCount:    word.ReviewCount,
		},
	}
}
