package chat

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/thompsonja/openai-go/pkg/helpers"
	"github.com/thompsonja/openai-go/pkg/usage"
)

const (
	chatCompletionsEndpoint = "https://api.openai.com/v1/chat/completions"
)

type CreateChatCompletionRequest struct {
	Model            string         `json:"model"`
	Messages         []string       `json:"messages"`
	Temperature      int            `json:"temperature,omitempty"`
	TopP             int            `json:"top_p,omitempty"`
	N                int            `json:"n,omitempty"`
	Stream           bool           `json:"stream,omitempty"`
	Stop             string         `json:"stop,omitempty"`
	MaxTokens        int            `json:"max_tokens,omitempty"`
	PresencePenalty  int            `json:"presence_penalty,omitempty"`
	FrequencyPenalty int            `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int `json:"logit_bias,omitempty"`
	User             string         `json:"user,omitempty"`
}

type createChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type createChatCompletionResponse struct {
	Index        int                         `json:"index"`
	Message      createChatCompletionMessage `json:"message"`
	FinishReason string                      `json:"finish_reason"`
}

type CreateChatCompletionResponse struct {
	ID      string                         `json:"id"`
	Object  string                         `json:"object"`
	Created int                            `json:"created"`
	Choices []createChatCompletionResponse `json:"choices"`
	Usage   usage.Usage                    `json:"usage"`
}

type API struct {
	requester *helpers.HttpRequester
}

func New(apiKey string) *API {
	return &API{
		requester: helpers.New(apiKey),
	}
}

func (a *API) Create(ctx context.Context, req *CreateChatCompletionRequest) (*CreateChatCompletionResponse, error) {
	resp, err := a.requester.SendHttpRequest(ctx, "POST", chatCompletionsEndpoint, "application/json", req)
	if err != nil {
		return nil, fmt.Errorf("a.requester.SendHttpRequest: %v", err)
	}

	var response *CreateChatCompletionResponse
	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return response, nil
}
