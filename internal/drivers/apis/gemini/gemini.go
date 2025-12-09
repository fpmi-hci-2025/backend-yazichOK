package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"speech-processing-service/internal/config"
	"speech-processing-service/internal/errs"
)

type Gemini struct {
	client *http.Client
	cfg    *config.ExternalAPI
}

func New(cfg *config.ExternalAPI) Gemini {
	return Gemini{
		client: &http.Client{},
		cfg:    cfg,
	}
}

func (g *Gemini) GetTranscriptionURL() string {
	return g.cfg.URL + "/models/gemini-2.0-flash:generateContent?key=" + g.cfg.APIKey
}

func (g *Gemini) AnalyzeText(ctx context.Context, prompt string) (string, error) {
	reqBody := AnalyzeTextReq{
		Contents: []Contents{
			{
				Parts: []Part{
					{
						Text: prompt,
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	if err := encoder.Encode(reqBody); err != nil {
		return "", errs.New(errs.ErrMarshalingJSON, err.Error())
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		g.GetTranscriptionURL(),
		&buf,
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Response body:", string(bodyBytes))
		return "", fmt.Errorf("failed to analyze text: %s", resp.Status)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errs.New(errs.ErrDecodingJSON, err.Error())
	}

	var result AnalyzeTextResp
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", err
	}

	if len(result.Candidates) > 0 && len(result.Candidates[0].Content.Parts) > 0 {
		return result.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("no text found in response")
}
