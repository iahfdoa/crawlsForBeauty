package scanner

import (
	"github.com/antchfx/htmlquery"
	"github.com/iahfdoa/crawlsForBeauty/util"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

func (s *Scanner) getByXpathUrl(url, xpathExpr string, out chan string) ([]string, error) {
	text, err := getText(s.client, url)
	if err != nil {
		return nil, err
	}

	doc, err := htmlquery.Parse(strings.NewReader(text))
	if err != nil {
		return nil, err
	}

	var urls []string
	attrs := htmlquery.Find(doc, xpathExpr)
	for _, node := range attrs {
		attrValue := htmlquery.InnerText(node)
		if attrValue != "" {
			urls = append(urls, attrValue)
		}
	}
	for _, u := range urls {
		out <- u
	}
	return urls, nil
}

func (s *Scanner) Scan() (err error) {
	if s.request["tag"] != "" {
		s.request["url"], err = url.JoinPath(s.request["url"], "tag", s.request["tag"])
		if err != nil {
			return err
		}
	}
	s.wg.Add()
	go func(U, xpath string) {
		defer func() {
			close(s.imgHtmlChan)
			s.wg.Done()
		}()
		var numberOfResults int
		page := 1
		for {
			limit, err := strconv.Atoi(s.request["limit"])
			if err != nil {
				break
			}
			u, _ := url.JoinPath(U, "page", strconv.Itoa(page))
			imgUrl, err := s.getByXpathUrl(u, "//a[contains(@class, 'entry-thumbnail')]/@href", s.imgHtmlChan)
			if err != nil {
				continue
			}

			if len(imgUrl) == 0 || len(imgUrl) > limit || numberOfResults > limit {
				break
			}
			page++
			numberOfResults += len(imgUrl) * 20
		}
	}(s.request["url"], "//a[contains(@class, 'entry-thumbnail')]/@href")
	s.wg.Add()
	go func() {
		defer func() {
			close(s.imgUrlChan)
			s.wg.Done()
		}()
		var numberOfResults int
		wg := sync.WaitGroup{}
		for img := range s.imgHtmlChan {
			wg.Add(1)
			go func(img string) {
				defer func() {
					wg.Done()
				}()
				imgUrl, err := s.getByXpathUrl(img, "//div[contains(@class, 'entry themeform')]//img/@src", s.imgUrlChan)
				if err != nil {
					return
				}
				numberOfResults += len(imgUrl)

			}(img)
		}

		wg.Wait()
	}()

	for u := range s.imgUrlChan {
		parse, err := url.Parse(u)
		if err != nil {
			continue
		}
		output := filepath.Join(s.output, path.Base(parse.Path))

		s.limiter.Take()
		s.wg.Add()
		go func(u, output string) {
			defer func() {
				s.NumChan <- struct{}{}
				s.wg.Done()
			}()
			if util.FileExists(output) {
				return
			}
			err := download(s.client, u, output)
			if err != nil {
				return
			}

		}(u, output)

	}
	s.wg.Wait()
	close(s.NumChan)
	return nil
}
