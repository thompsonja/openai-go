package model

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/thompsonja/openai-go/pkg/helpers"
)

const (
	modelsEndpoint = "https://api.openai.com/v1/models"
)

type Permission struct {
}

type Model struct {
	ID         string       `json:"id"`
	Object     string       `json:"object"`
	OwnedBy    string       `json:"owned_by"`
	Permission []Permission `json:"permission"`
}

type API struct {
	requester *helpers.HttpRequester
}

func New(apiKey string) *API {
	return &API{
		requester: helpers.New(apiKey),
	}
}

func (a *API) List(ctx context.Context) ([]*Model, error) {
	resp, err := a.requester.SendHttpRequest(ctx, "GET", modelsEndpoint, "application/json", nil)
	if err != nil {
		return nil, fmt.Errorf("a.requester.SendHttpRequest: %v", err)
	}

	var models []*Model
	if err := json.Unmarshal(resp, &models); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return models, nil
}

func (a *API) GetModel(ctx context.Context, modelName string) (*Model, error) {
	resp, err := a.requester.SendHttpRequest(ctx, "GET", fmt.Sprintf("modelsEndpoint/%s", modelName), "application/json", nil)
	if err != nil {
		return nil, fmt.Errorf("a.requester.SendHttpRequest: %v", err)
	}

	var model *Model
	if err := json.Unmarshal(resp, &model); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return model, nil
}
