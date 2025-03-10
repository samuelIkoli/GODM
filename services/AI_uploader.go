package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html/charset"
)

func UploadHelper(from string) (string, error) {
	// logic to help upload document to AI can be here
	resp, err := http.Get(from)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch file: %s", resp.Status)
	}

	buf := new(bytes.Buffer)
	_, buffErr := io.Copy(buf, resp.Body)
	if buffErr != nil {
		return "", buffErr
	}

	utf8Reader, err := charset.NewReader(buf, "")
	if err != nil {
		return "", fmt.Errorf("error converting encoding to UTF-8: %w", err)
	}

	body, err := io.ReadAll(utf8Reader)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
