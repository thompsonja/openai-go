package client

import (
	"github.com/thompsonja/openai-go/pkg/image"
)

type Client struct {
	Image *image.API
}

func New(apiKey string) *Client {
	return &Client{
		Image: image.New(apiKey),
	}
}
