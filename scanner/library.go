package scanner

import (
	"errors"
	"fmt"
	"strings"
)

func (s *Scanner) selection() error {
	switch s.request["type"].(type) {
	case int:
		switch s.request["type"] {
		case 3:
			switch {
			case s.request["tag"] != "":
				s.request["path"] = fmt.Sprintf("/cat/%s?page=", s.request["tag"])
			default:
				s.request["path"] = fmt.Sprintf("/cat/7?page=")
			}
			s.request["url"] = "https://www.photos18.com"
			s.request["img_path"] = "//div[@class=\"card\"]//a[@class=\"visited\"]/@href"
			s.request["img_url"] = "//div[@class=\"fullwidth\"]//a/@href"
		case 2:
			tag := s.request["tag"].(string)
			tagInfo := strings.Split(tag, ":")
			if len(tagInfo) < 2 {
				return errors.New("error tag")
			}
			var topRange string

			var categories string
			switch tagInfo[0] {
			case "toplist":
				topRange = "1y"
			default:
				topRange = ""
			}
			switch tagInfo[1] {
			case "5":
				categories = "111"
			case "96280":
				categories = "111"
			default:
				categories = "101"
			}
			sorting := tagInfo[0]
			id := tagInfo[1]
			switch id {
			case "0":
				s.request["path"] = "/search?categories=110&purity=010&atleast=1920x1080&topRange=1y&sorting=toplist&order=desc&ai_art_filter=1&page="
			default:
				s.request["path"] = fmt.Sprintf("/search?q=id:%s&atleast=1920x1080&categories=%s&purity=010&sorting=%s&order=desc&topRange=%s&page=", id, categories, sorting, topRange)
			case "-1":
				s.request["path"] = fmt.Sprintf("/search?categories=110&purity=100&atleast=1920x1080&topRange=1M&sorting=toplist&order=desc&page=")
			}

			s.request["url"] = "https://wallhaven.cc"
			s.request["img_path"] = "//section[@class='thumb-listing-page']//a[@class='preview']/@href"
			s.request["img_url"] = ""
			s.UrlCallBackFunc = func(str string) (string, error) {
				if len(str) <= 6 {
					return "", errors.New("字符串长度小于6位")
				}
				length := len(str)
				name := str[length-6:]
				url := fmt.Sprintf("https://w.wallhaven.cc/full/%s/wallhaven-%s.jpg", name[:2], name)
				return url, nil
			}
		case 1:
			switch {
			case s.request["tag"] != "":
				s.request["path"] = fmt.Sprintf("/pictures/%s", s.request["tag"])
			default:
				s.request["path"] = "/pictures/flexible-porn"
			}
			s.request["url"] = "https://sexroom.xxx/"
			s.request["img_path"] = "//*[@id=\"list_albums_common_albums_list_items\"]//a[@class='img']/@href"
			s.request["img_url"] = "//*[@id=\"aniimated-thumbnials\"]/a/@href"
		default:
			switch {
			case s.request["tag"] == "xiaoxinggan" || s.request["tag"] == "xiaotianmei":
				s.request["path"] = fmt.Sprintf("/topic/%s/page/", s.request["tag"])
			case s.request["tag"] != "" && !strings.HasPrefix(s.request["tag"].(string), "xiao"):
				s.request["path"] = fmt.Sprintf("/tag/%s/page/", s.request["tag"])
			default:
				s.request["path"] = fmt.Sprintf("/page/")
			}
			s.request["url"] = "https://mzt8.com"
			s.request["img_path"] = "//a[contains(@class, 'entry-thumbnail')]/@href"
			s.request["img_url"] = "//div[contains(@class, 'entry themeform')]//img/@src"
		}
	default:
		return errors.New("错误的类型")
	}
	return nil
}
