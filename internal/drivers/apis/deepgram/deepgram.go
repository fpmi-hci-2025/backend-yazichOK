package deepgram

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

type Deepgram struct {
	client *http.Client
	cfg    *config.ExternalAPI
}

func New(cfg *config.ExternalAPI) Deepgram {
	return Deepgram{
		client: &http.Client{},
		cfg:    cfg,
	}
}

func (deepgram *Deepgram) GetTranscriptionURL() string {
	return deepgram.cfg.URL + "listen?model=nova-3&smart_format=true"
}

func (deepgram *Deepgram) GetTranscriptionHeaders() map[string]string {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Token " + deepgram.cfg.APIKey,
	}
	return headers
}

func (deepgram *Deepgram) TranscribeAudio(ctx context.Context, url string) (string, error) {
	reqBody := TranscribeTextReq{
		URL: url,
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(reqBody); err != nil {
		return "", errs.New(errs.ErrMarshalingJSON, err.Error())
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		deepgram.GetTranscriptionURL(),
		&buf,
	)
	if err != nil {
		return "", errs.New(errs.ErrExecutionRequest, err.Error())
	}

	for key, value := range deepgram.GetTranscriptionHeaders() {
		req.Header.Set(key, value)
	}

	resp, err := deepgram.client.Do(req)
	if err != nil {
		return "", errs.New(errs.ErrExecutionRequest, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errs.New(errs.ErrUnexpectedStatusCode, fmt.Sprintf("status_code:%d", resp.StatusCode))
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errs.New(errs.ErrDecodingJSON, err.Error())
	}

	var transcriptionResp TranscribeTextResp
	if err := json.Unmarshal(bodyBytes, &transcriptionResp); err != nil {
		return "", errs.New(errs.ErrDecodingJSON, err.Error())
	}

	return transcriptionResp.Results.Channels[0].Alternatives[0].Transcript, nil
}
