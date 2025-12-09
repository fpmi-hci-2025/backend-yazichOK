package get_article_by_id

import (
	"context"
	"speech-processing-service/internal/drivers/storage"
	"speech-processing-service/internal/entity"
	"speech-processing-service/internal/errs"
)

type StorageProvider interface {
	GetArticleByID(ctx context.Context, id int) (storage.Article, error)
	GetArticleVocabulary(ctx context.Context, articleID int) ([]storage.ArticleVocabulary, error)
	GetArticleGrammarRules(ctx context.Context, articleID int) ([]storage.ArticleGrammarRule, error)
}

type URLGetter interface {
	GenerateUrl(ctx context.Context, imagePath string, isAnswer bool) (string, error)
}

type UseCase struct {
	storage   StorageProvider
	urlGetter URLGetter
}

func New(storage StorageProvider, urlGetter URLGetter) UseCase {
	return UseCase{
		storage:   storage,
		urlGetter: urlGetter,
	}
}

func (u *UseCase) GetArticle(ctx context.Context, id int) (entity.Article, error) {
	article, err := u.storage.GetArticleByID(ctx, id)
	if err != nil {
		return entity.Article{}, errs.New(errs.ErrUseCaseExecution, "u.storage.GetArticleByID: "+err.Error())
	}

	vocabulary, err := u.storage.GetArticleVocabulary(ctx, id)
	if err != nil {
		return entity.Article{}, errs.New(errs.ErrUseCaseExecution, "u.storage.GetArticleVocabulary: "+err.Error())
	}

	rules, err := u.storage.GetArticleGrammarRules(ctx, id)
	if err != nil {
		return entity.Article{}, errs.New(errs.ErrUseCaseExecution, "u.storage.GetArticleGrammarRules: "+err.Error())
	}

	vocabEntities := make([]entity.VocabularyWord, 0, len(vocabulary))
	for _, word := range vocabulary {
		vocabEntities = append(vocabEntities, entity.VocabularyWord{
			ID:           word.ID,
			Word:         word.Word,
			PartOfSpeech: word.PartOfSpeech,
			Meaning:      word.Meaning,
		})
	}

	ruleEntities := make([]entity.GrammarRule, 0, len(rules))
	for _, rule := range rules {
		ruleEntities = append(ruleEntities, entity.GrammarRule{
			ID:      rule.ID,
			Name:    rule.Name,
			Example: rule.Example,
			Note:    rule.Note,
		})
	}

	imageURL, err := u.urlGetter.GenerateUrl(ctx, article.ImageURL, false)
	if err != nil {
		return entity.Article{}, errs.Wrap("u.urlGetter.GenerateURl", err)
	}

	return entity.Article{
		ID:         article.ID,
		ImageURL:   imageURL,
		Content:    article.Content,
		Title:      article.Title,
		Level:      article.Level,
		Minutes:    article.MinutesToRead,
		Vocabulary: vocabEntities,
		Rules:      ruleEntities,
	}, nil
}
