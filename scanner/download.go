package scanner

import (
	"io"
	"net/http"
	"os"
	"strings"
)

func download(client *http.Client, url, path string) error {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Referer", url)
	text, err := client.Do(request)
	if err != nil {
		return err
	}
	if text.StatusCode != 200 {
		return nil
	}
	defer text.Body.Close()
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, text.Body)
	return err
}

func getText(client *http.Client, url string) (string, error) {
	response, err := client.Get(url)
	if err != nil {
		return "", err // Return empty string and the error.
	}

	defer response.Body.Close()
	buffer := strings.Builder{}

	_, err = io.Copy(&buffer, response.Body)
	if err != nil {
		return "", err // Return empty string and the error.
	}

	return buffer.String(), nil // Return the text as a string and no error.
}
