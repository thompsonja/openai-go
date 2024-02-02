package moderations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/thompsonja/openai-go/pkg/helpers"
)

const (
	moderationsEndpoint = "https://api.openai.com/v1/moderations"
)

type API struct {
	requester *helpers.HttpRequester
}

func New(apiKey string) *API {
	return &API{
		requester: helpers.New(apiKey),
	}
}

type ModerationRequest struct {
	Input string `json:"input"`
	Model string `json:"model,omitempty"`
}

type ModerationResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Results []struct {
		Flagged        bool               `json:"flagged"`
		Categories     map[string]bool    `json:"categories"`
		CategoryScores map[string]float32 `json:"category_scores"`
	} `json:"results"`
}

func (a *API) CreateModeration(ctx context.Context, req *ModerationRequest) (*ModerationResponse, error) {
	b, err := a.requester.SendHttpRequest(ctx, "POST", moderationsEndpoint, "application/json", req)
	if err != nil {
		return nil, fmt.Errorf("a.requester.SendHttpRequest: %w", err)
	}
	var data ModerationResponse
	if err = json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}
	return &data, nil
}

func GetAllFlaggedCategories(resp *ModerationResponse) []string {
	flaggedCategories := []string{}
	if len(resp.Results) > 0 {
		for _, result := range resp.Results {
			if !result.Flagged {
				continue
			}
			for k, v := range result.Categories {
				if !v {
					continue
				}
				flaggedCategories = append(flaggedCategories, k)
			}
		}
	}
	return flaggedCategories
}
