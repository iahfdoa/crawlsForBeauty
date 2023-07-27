package scanner

import (
	"net/http"
	"testing"
)

func TestScanner_Scan(t *testing.T) {
	client := &http.Client{}
	options := &Options{
		Client: client,
		Rate:   30,
		Tag:    "meitui",
		Url:    "https://mzt8.com/",
		Limit:  200,
		Output: "2/",
	}
	s, err := NewScanner(options)
	if err != nil {
		return
	}
	err = s.Scan()
	if err != nil {
		return
	}
}
