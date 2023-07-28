package scanner

import (
	"github.com/antchfx/htmlquery"
	"github.com/iahfdoa/crawlsForBeauty/util"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
	"strings"
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
	err = s.selection()
	if err != nil {
		return err
	}
	var AllOfResults int
	s.wg.Add()
	go func(U, xpath string) {
		defer func() {
			close(s.imgHtmlChan)
			s.wg.Done()
		}()

		page := 1
		for {
			limit, err := strconv.Atoi(s.request["limit"].(string))
			if err != nil {
				break
			}
			u, _ := url.JoinPath(U, strconv.Itoa(page))
			imgUrl, err := s.getByXpathUrl(u, xpath, s.imgHtmlChan)
			if err != nil {
				continue
			}
			s.lock.Lock()
			if len(imgUrl) == 0 || len(imgUrl) > limit || AllOfResults > limit {
				s.lock.Unlock()
				break
			}
			s.lock.Unlock()
			page++
		}
	}(s.request["url"].(string)+s.request["path"].(string), s.request["img_path"].(string))
	s.wg.Add()
	go func(xpath string) {
		defer func() {
			close(s.imgUrlChan)

			s.wg.Done()
		}()
		for img := range s.imgHtmlChan {
			func(img string) {
				imgUrl, err := s.getByXpathUrl(img, xpath, s.imgUrlChan)
				if err != nil {
					return
				}
				s.lock.Lock()
				AllOfResults += len(imgUrl)
				s.lock.Unlock()
			}(img)
		}

	}(s.request["img_url"].(string))

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
