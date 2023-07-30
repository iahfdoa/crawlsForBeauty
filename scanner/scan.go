package scanner

import (
	"context"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/iahfdoa/crawlsForBeauty/util"
	"github.com/projectdiscovery/gologger"
	_ "golang.org/x/image/webp"
	"image"
	"image/png"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func (s *Scanner) getByXpathUrl(ctx context.Context, url, xpathExpr string, out chan string) ([]string, error) {
	text, err := getText(s.client, url)
	if err != nil {
		return nil, err
	}
	//gologger.Debug().Msg(text)
	doc, err := htmlquery.Parse(strings.NewReader(text))
	if err != nil {
		return nil, err
	}

	var urls []string
	attrs := htmlquery.Find(doc, xpathExpr)
	if len(attrs) == 0 {
		gologger.Debug().Msg("没有任何数据找到")
	}
	//fmt.Println(attrs)
	for _, node := range attrs {
		attrValue := htmlquery.InnerText(node)
		if attrValue != "" {
			urls = append(urls, attrValue)
		}
	}
	for _, u := range urls {
		select {
		case out <- u:
			//gologger.Debug().Msg(u)
		case <-ctx.Done():
			return nil, nil
		}
	}
	return urls, nil
}

func (s *Scanner) Scan() (err error) {
	err = s.selection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	limit, err := strconv.Atoi(s.request["limit"].(string))
	var AllOfResults int
	s.wg.Add()
	go func(U, xpath string) {
		defer func() {
			close(s.imgHtmlChan)
			s.wg.Done()
		}()

		page := 0
		for {
			page++
			if err != nil {
				break
			}
			var u string
			switch s.request["type"] {
			case 3:
				u = U + strconv.Itoa(page)
			case 2:
				u = U + strconv.Itoa(page)
			default:
				u, _ = url.JoinPath(U, strconv.Itoa(page))
			}

			imgUrl, err := s.getByXpathUrl(ctx, u, xpath, s.imgHtmlChan)
			if err != nil {
				gologger.Debug().Msg(err.Error())
				continue
			}
			//for _, u := range imgUrl {
			//	gologger.Debug().Msg(u)
			//}
			if len(imgUrl) == 0 || len(imgUrl) > limit || AllOfResults > limit {
				cancel()
				break
			}

		}
	}(s.request["url"].(string)+s.request["path"].(string), s.request["img_path"].(string))
	s.wg.Add()
	go func(xpath string) {
		defer func() {
			close(s.imgUrlChan)

			s.wg.Done()
		}()
		for img := range s.imgHtmlChan {
			switch s.request["type"] {
			case 2:
				select {
				case s.imgUrlChan <- img:
					s.lock.Lock()
					AllOfResults += 1
					s.lock.Unlock()
					if AllOfResults > limit {
						cancel()
						return
					}
				case <-ctx.Done():
					return
				}
			default:
				if s.request["type"] == 3 {
					img = fmt.Sprintf("%s%s", s.request["url"], img)
				}
				imgUrl, _ := s.getByXpathUrl(ctx, img, xpath, s.imgUrlChan)
				if len(imgUrl) == 0 || len(imgUrl) > limit || AllOfResults > limit {
					cancel()
					return
				}
				s.lock.Lock()
				AllOfResults += len(imgUrl)
				s.lock.Unlock()
			}

		}

	}(s.request["img_url"].(string))

	for u := range s.imgUrlChan {
		if s.UrlCallBackFunc != nil {
			u, err = s.UrlCallBackFunc(u)
			if err != nil {
				gologger.Debug().Msg(err.Error())
				continue
			}
		}
		parse, err := url.Parse(u)
		if err != nil {
			gologger.Debug().Msg(err.Error())
			continue
		}
		var webpToPng func(io.Writer, io.Reader) error
		filename := path.Base(parse.Path)
		ext := path.Ext(filename)
		switch ext {
		case ".webp":
			filename = strings.TrimSuffix(filename, ext) + ".png"
			webpToPng = func(writer io.Writer, reader io.Reader) error {
				// Decode the WebP image from the []byte data
				webpImage, _, err := image.Decode(reader)
				if err != nil {
					return fmt.Errorf("error decoding WebP image: %s", err)
				}

				// Encode the WebP image as PNG and write it to the file
				err = png.Encode(writer, webpImage)
				if err != nil {
					return fmt.Errorf("error encoding WebP image to PNG: %s", err)
				}

				return nil
			}
		}
		output := filepath.Join(s.output, filename)
		u = strings.Split(u, "?")[0]

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

			gologger.Debug().Msgf("开始下载:%s", u)
			err := download(s.client, u, output, webpToPng)
			if err != nil {
				gologger.Debug().Msg(err.Error())
				os.Remove(output)
				return
			}

		}(u, output)

	}
	s.wg.Wait()
	close(s.NumChan)
	return nil
}
