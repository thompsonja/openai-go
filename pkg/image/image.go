package image

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/thompsonja/openai-go/pkg/helpers"
)

const (
	Small int = iota
	Medium
	Large
)

const (
	createEndpoint    = "https://api.openai.com/v1/images/generations"
	editEndpoint      = "https://api.openai.com/v1/images/edits"
	variationEndpoint = "https://api.openai.com/v1/images/variations"
)

var (
	sizes = map[string]string{
		"small":          "256x256",
		"medium":         "512x512",
		"large":          "1024x1024",
		"dalle-3-square": "1024x1024",
		"dalle-3-wide":   "1792x1024",
		"dalle-3-tall":   "1024x1792",
	}
)

type API struct {
	requester *helpers.HttpRequester
}

func New(apiKey string) *API {
	return &API{
		requester: helpers.New(apiKey),
	}
}

type CreateRequest struct {
	Model          string `json:"model,omitempty"`
	Prompt         string `json:"prompt,omitempty"`
	Quality        string `json:"quality,omitempty"`
	N              int    `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
}

func (a *API) Create(ctx context.Context, req *CreateRequest) ([]byte, error) {
	if _, ok := sizes[req.Size]; !ok {
		return nil, fmt.Errorf("invalid size input: %s", req.Size)
	}
	req.Size = sizes[req.Size]
	if req.Model != "dall-e-2" && req.Model != "dall-e-3" {
		return nil, fmt.Errorf("invalid model: %s", req.Model)
	}
	if req.Model != "dall-e-3" && req.Quality == "hd" {
		return nil, fmt.Errorf("hd quality is supported only by dall-e-3")
	}
	if req.ResponseFormat == "" {
		req.ResponseFormat = "b64_json"
	}
	return a.requester.SendHttpRequest(ctx, "POST", createEndpoint, "application/json", req)
}

type EditRequest struct {
	Prompt         string `json:"prompt,omitempty"`
	N              int    `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	Image          string `json:"image,omitempty"`
	Mask           string `json:"mask,omitempty"`
}

func (a *API) Edit(ctx context.Context, req *EditRequest) ([]byte, error) {
	if _, ok := sizes[req.Size]; !ok {
		return nil, fmt.Errorf("invalid size input: %s", req.Size)
	}
	req.Size = sizes[req.Size]
	return a.requester.SendHttpRequest(ctx, "POST", editEndpoint, "application/json", req)
}

func (a *API) EditWithUrls(ctx context.Context, req *EditRequest, imageUrl, maskUrl string) ([]byte, error) {
	sizeStr, ok := sizes[req.Size]
	if !ok {
		return nil, fmt.Errorf("invalid size input: %s", req.Size)
	}

	image, err := helpers.DownloadPng(imageUrl)
	if err != nil {
		return nil, fmt.Errorf("helpers.DownloadPng: %v", err)
	}
	defer os.Remove(image)
	mask, err := helpers.DownloadPng(maskUrl)
	if err != nil {
		return nil, fmt.Errorf("helpers.DownloadPng: %v", err)
	}
	defer os.Remove(mask)

	if err := helpers.VerifyPngs([]string{image, mask}); err != nil {
		return nil, fmt.Errorf("helpers.VerifyPngs: %v", err)
	}

	fd := map[string]string{
		"prompt":          req.Prompt,
		"n":               strconv.Itoa(req.N),
		"size":            sizeStr,
		"response_format": "b64_json",
		"image":           fmt.Sprintf("@%s", image),
		"mask":            fmt.Sprintf("@%s", mask),
	}

	ct, body, err := helpers.CreateMultipartFormData(fd)
	if err != nil {
		return nil, fmt.Errorf("helpers.CreateMultipartFormData: %v", err)
	}

	return a.requester.SendHttpRequest(ctx, "POST", editEndpoint, ct, body)
}

func (a *API) Variation(ctx context.Context, req *EditRequest) ([]byte, error) {
	if _, ok := sizes[req.Size]; !ok {
		return nil, fmt.Errorf("invalid size input: %s", req.Size)
	}
	return a.requester.SendHttpRequest(ctx, "POST", variationEndpoint, "application/json", req)
}

func (a *API) VariationWithUrls(ctx context.Context, imageUrl, size string, n int) ([]byte, error) {
	sizeStr, ok := sizes[size]
	if !ok {
		return nil, fmt.Errorf("invalid size input: %s", size)
	}

	p, err := helpers.DownloadPng(imageUrl)
	if err != nil {
		return nil, fmt.Errorf("helpers.DownloadPng: %v", err)
	}
	defer os.Remove(p)

	if err := helpers.VerifyPngs([]string{p}); err != nil {
		return nil, fmt.Errorf("helpers.VerifyPngs: %v", err)
	}

	fd := map[string]string{
		"n":               strconv.Itoa(n),
		"size":            sizeStr,
		"response_format": "b64_json",
		"image":           fmt.Sprintf("@%s", p),
	}

	ct, body, err := helpers.CreateMultipartFormData(fd)
	if err != nil {
		return nil, fmt.Errorf("helpers.CreateMultipartFormData: %v", err)
	}

	return a.requester.SendHttpRequest(ctx, "POST", variationEndpoint, ct, body)
}
