package completion

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/thompsonja/openai-go/pkg/helpers"
)

const (
	completionsEndpoint = "https://api.openai.com/v1/completions"
)

type CreateCompletionsRequest struct {
	Model            string `json:"model"`
	Prompt           string `json:"prompt"`
	Suffix           string `json:"suffix,omitempty"`
	MaxTokens        int    `json:"max_tokens"`
	Temperature      int    `json:"temperature"`
	TopP             int    `json:"top_p"`
	N                int    `json:"n"`
	Stream           bool   `json:"stream"`
	LogProbs         int    `json:"log_probs"`
	Echo             bool   `json:"echo"`
	PresencePenalty  int    `json:"presence_penalty"`
	FrequencyPenalty int    `json:"frequency_penalty"`
}

type Model struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	OwnedBy string `json:"owned_by"`
}

type API struct {
	requester *helpers.HttpRequester
}

func New(apiKey string) *API {
	return &API{
		requester: helpers.New(apiKey),
	}
}

func (a *API) Create(ctx context.Context, req *CreateCompletionsRequest) ([]*Model, error) {
	resp, err := a.requester.SendHttpRequest(ctx, "POST", completionsEndpoint, "application/json", req)
	if err != nil {
		return nil, fmt.Errorf("a.requester.SendHttpRequest: %v", err)
	}

	var models []*Model
	if err := json.Unmarshal(resp, &models); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return models, nil
}
