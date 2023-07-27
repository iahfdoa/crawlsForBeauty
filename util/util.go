package util

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CreateDir(dirPath string) error {
	outputFolder := filepath.Dir(dirPath)
	if mkdirErr := os.MkdirAll(outputFolder, 0700); mkdirErr != nil {
		return mkdirErr
	}

	return nil
}

func FolderExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return info.IsDir()
}
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return !info.IsDir()
}
func parseProxyAuth(auth string) (string, string, bool) {
	parts := strings.SplitN(auth, ":", 2)
	if len(parts) != 2 || parts[0] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// getProxyFunc 辅助函数：获取代理设置函数
func getProxyFunc(proxy, auth string) func(*http.Request) (*url.URL, error) {
	if proxy == "" {
		return nil
	}
	proxyURL, err := url.Parse(proxy)
	if err != nil {

		return nil
	}
	if auth != "" {
		username, password, ok := parseProxyAuth(auth)
		if !ok {
			return nil
		}
		proxyURL.User = url.UserPassword(username, password)
	}
	return http.ProxyURL(proxyURL)
}

func NewClient(proxy, auth string) *http.Client {
	transport := &http.Transport{
		Proxy:           getProxyFunc(proxy, auth),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true, MinVersion: tls.VersionTLS10},
	}
	return &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}
}
