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
