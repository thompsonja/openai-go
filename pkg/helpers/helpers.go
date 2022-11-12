package helpers

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

const openaiImageFilesizeLimitMb = 4

func CreateMultipartFormData(formData map[string]string) (string, io.Reader, error) {
	b := new(bytes.Buffer)
	w := multipart.NewWriter(b)
	defer w.Close()

	for fieldName, data := range formData {
		// local files have prefix '@'
		if strings.HasPrefix(data, "@") {
			data = data[1:]
			f, err := os.Open(data)
			if err != nil {
				return "", nil, fmt.Errorf("os.Open: %v", err)
			}
			defer f.Close()
			p, err := w.CreateFormFile(fieldName, data)
			if err != nil {
				return "", nil, fmt.Errorf("w.CreateFormFile: %v", err)
			}
			io.Copy(p, f)
		} else {
			if err := w.WriteField(fieldName, data); err != nil {
				return "", nil, fmt.Errorf("w.WriteField: %v", err)
			}
		}
	}
	return w.FormDataContentType(), b, nil
}

func DalleRequestResponse(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid http status: %d", resp.StatusCode)
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}
	return responseBody, nil
}

func DownloadPng(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("http.Get: %v", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received unexpected status code: %v", response.StatusCode)
	}

	f, err := os.CreateTemp("", "fetched*.png")
	if err != nil {
		return "", fmt.Errorf("os.CreateTemp: %v", err)
	}

	_, err = io.Copy(f, response.Body)
	if err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}

	return f.Name(), nil
}

func VerifyPngs(pngPaths []string) error {
	foundPng := false
	var expectedWidth, expectedHeight int

	for _, pngPath := range pngPaths {
		f, err := os.Open(pngPath)
		if err != nil {
			return fmt.Errorf("os.Open: %v", err)
		}

		fi, err := f.Stat()
		if err != nil {
			return fmt.Errorf("f.Stat: %v", err)
		}
		if fi.Size() > openaiImageFilesizeLimitMb*1024*1024 {
			return fmt.Errorf("image size too large, must be under %d MB", openaiImageFilesizeLimitMb)
		}

		image, err := png.Decode(f)
		if err != nil {
			return fmt.Errorf("image must be valid png, png.Decode: %v", err)
		}
		width := image.Bounds().Dx()
		height := image.Bounds().Dy()
		if width != height {
			return fmt.Errorf("found non-square image with dimensions %dx%d", width, height)
		}

		if !foundPng {
			foundPng = true
			expectedWidth = width
			expectedHeight = height
		} else {
			if width != expectedWidth || height != expectedHeight {
				return fmt.Errorf("dimensions of all images must match, got both (%dx%d) and (%dx%d)", width, height, expectedWidth, expectedHeight)
			}
		}
	}

	return nil
}
