package deepgram

type TranscribeTextReq struct {
	URL string `json:"url"`
}

type Results struct {
	Channels []Channels `json:"channels"`
}

type Channels struct {
	Alternatives []Alternatives `json:"alternatives"`
}

type Alternatives struct {
	Transcript string `json:"transcript"`
}

type TranscribeTextResp struct {
	Results Results `json:"results"`
}
