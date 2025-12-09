package views

type StartSessionRequest struct {
	TopicID int `json:"topic_id"`
}

type AddWordToCollectionRequest struct {
	Word        string  `json:"word"`
	Translation string  `json:"translation"`
	Example     *string `json:"example"`
}
