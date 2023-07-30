package scanner

import (
	"io"
	"net/http"
	"os"
	"strings"
)

func download(client *http.Client, url, path string, webpToPng func(dst io.Writer, src io.Reader) error) error {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Referer", url)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
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
	if webpToPng != nil {
		err = webpToPng(out, text.Body)
	} else {

		_, err = io.Copy(out, text.Body)
	}
	return err
}

func getText(client *http.Client, url string) (string, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("Referer", url)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	response, err := client.Do(request)
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
