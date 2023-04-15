package completion

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/thompsonja/openai-go/pkg/helpers"
	"github.com/thompsonja/openai-go/pkg/usage"
)

const (
	completionsEndpoint = "https://api.openai.com/v1/completions"
)

type CreateCompletionRequest struct {
	Model            string         `json:"model"`
	Prompt           string         `json:"prompt,omitempty"`
	Suffix           string         `json:"suffix,omitempty"`
	MaxTokens        int            `json:"max_tokens,omitempty"`
	Temperature      int            `json:"temperature,omitempty"`
	TopP             int            `json:"top_p,omitempty"`
	N                int            `json:"n,omitempty"`
	Stream           bool           `json:"stream,omitempty"`
	LogProbs         int            `json:"log_probs,omitempty"`
	Echo             bool           `json:"echo,omitempty"`
	Stop             string         `json:"stop,omitempty"`
	PresencePenalty  int            `json:"presence_penalty,omitempty"`
	FrequencyPenalty int            `json:"frequency_penalty,omitempty"`
	BestOf           int            `json:"best_of,omitempty"`
	LogitBias        map[string]int `json:"logit_bias,omitempty"`
	User             string         `json:"user,omitempty"`
}

type logProbs struct {
	// TODO populate
}

type createCompletionResponse struct {
	Text         string   `json:"text"`
	Index        int      `json:"index"`
	LogProbs     logProbs `json:"logprobs"`
	FinishReason string   `json:"finish_reason"`
}

type CreateCompletionResponse struct {
	ID      string                     `json:"id"`
	Object  string                     `json:"object"`
	Created int                        `json:"created"`
	Model   string                     `json:"model"`
	Choices []createCompletionResponse `json:"choices"`
	Usage   usage.Usage                `json:"usage"`
}

type API struct {
	requester *helpers.HttpRequester
}

func New(apiKey string) *API {
	return &API{
		requester: helpers.New(apiKey),
	}
}

func (a *API) Create(ctx context.Context, req *CreateCompletionRequest) (*CreateCompletionResponse, error) {
	resp, err := a.requester.SendHttpRequest(ctx, "POST", completionsEndpoint, "application/json", req)
	if err != nil {
		return nil, fmt.Errorf("a.requester.SendHttpRequest: %v", err)
	}

	var response *CreateCompletionResponse
	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return response, nil
}
