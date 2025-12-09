package gemini

type Part struct {
	Text string `json:"text"`
}

type Contents struct {
	Parts []Part `json:"parts"`
}

type AnalyzeTextReq struct {
	Contents []Contents `json:"contents"`
}

type AnalyzeTextResp struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content Content `json:"content"`
}

type Content struct {
	Parts []Part `json:"parts"`
}
