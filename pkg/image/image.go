package image

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
		"small":  "256x256",
		"medium": "512x512",
		"large":  "1024x1024",
	}
)

type API struct {
	apiKey string
}

func New(apiKey string) *API {
	return &API{
		apiKey: apiKey,
	}
}

func (a *API) Create(prompt, size string, n int) ([]byte, error) {
	sizeStr, ok := sizes[size]
	if !ok {
		return nil, fmt.Errorf("invalid size input: %s", size)
	}
	vals := map[string]interface{}{
		"prompt":          prompt,
		"n":               n,
		"size":            sizeStr,
		"response_format": "b64_json",
	}

	jv, err := json.Marshal(vals)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %v", err)
	}

	req, err := http.NewRequest("POST", createEndpoint, bytes.NewBuffer(jv))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.apiKey))

	return helpers.DalleRequestResponse(req)
}

func (a *API) Edit(imageUrl, maskUrl, prompt, size string, n int) ([]byte, error) {
	sizeStr, ok := sizes[size]
	if !ok {
		return nil, fmt.Errorf("invalid size input: %s", size)
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
		"prompt":          prompt,
		"n":               strconv.Itoa(n),
		"size":            sizeStr,
		"response_format": "b64_json",
		"image":           fmt.Sprintf("@%s", image),
		"mask":            fmt.Sprintf("@%s", mask),
	}

	ct, body, err := helpers.CreateMultipartFormData(fd)
	if err != nil {
		return nil, fmt.Errorf("helpers.CreateMultipartFormData: %v", err)
	}

	form := url.Values{}
	form.Add("prompt", prompt)
	form.Add("n", strconv.Itoa(n))
	form.Add("size", sizeStr)
	form.Add("response_format", "b64_json")
	form.Add("image", imageUrl)
	form.Add("mask", maskUrl)

	req, err := http.NewRequest("POST", editEndpoint, body)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %v", err)
	}

	req.Header.Add("Content-Type", ct)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.apiKey))

	return helpers.DalleRequestResponse(req)
}

func (a *API) Variation(imageUrl, size string, n int) ([]byte, error) {
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

	req, err := http.NewRequest("POST", variationEndpoint, body)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %v", err)
	}

	req.Header.Add("Content-Type", ct)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.apiKey))

	return helpers.DalleRequestResponse(req)
}
