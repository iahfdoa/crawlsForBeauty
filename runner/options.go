package runner

import (
	"github.com/gosuri/uiprogress"
	"github.com/gosuri/uiprogress/util/strutil"
	"github.com/iahfdoa/crawlsForBeauty/scanner"
	"github.com/iahfdoa/crawlsForBeauty/util"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	"os"
	"time"
)

type Options struct {
	Rate      int
	Tag       int
	Type      int
	Limit     int
	Output    string
	Proxy     string
	ProxyAuth string
	Version   bool
}

type Runner struct {
	Options *Options
	Scanner *scanner.Scanner
	Bar     *uiprogress.Bar
}

func ParserOptions() *Options {
	options := &Options{}
	set := goflags.NewFlagSet()
	set.SetDescription("crawlsForBeauty 美女图片爬虫工具")
	set.CreateGroup("Input", "输入",
		set.IntVar(&options.Type, "type", 0, "图库(0,1)"),
		set.IntVarP(&options.Tag, "tag", "t", 0, "类型 （1,2,3,4,5,6）"),
	)
	set.CreateGroup("Config", "配置",
		set.IntVarP(&options.Limit, "limit", "l", 600, "至少获取多少图片"),
		set.StringVarP(&options.Output, "output", "o", "output/", "输出位置"),
		set.StringVarP(&options.Proxy, "proxy", "p", "", "代理支持（socks5、http）"),
		set.StringVarP(&options.ProxyAuth, "proxy-auth", "pa", "", "代理认证（user:password）"),
		set.IntVarP(&options.Rate, "rate", "r", 40, "线程"),
	)
	set.CreateGroup("Version", "版本",
		set.BoolVarP(&options.Version, "version", "v", false, "显示版本"),
	)
	_ = set.Parse()
	showBanner()
	if options.Version {
		gologger.Print().Msgf("%s version: %s", Name, Version)
		os.Exit(0)
	}
	return options
}

func NewRunner(options *Options) (*Runner, error) {
	uiprogress.Start() // 开始进度条
	bar := uiprogress.AddBar(options.Limit)
	bar.Width = 100
	bar.AppendCompleted()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return strutil.PadLeft(b.TimeElapsedString(), 2, ' ')
	})
	uiprogress.RefreshInterval = 100 * time.Millisecond
	tagFunc := func(t, i int) string {
		switch t {
		case 1:
			switch i {
			case 1:
				return "stockings-porn"
			case 2:
				return "foot-fetish-porn"
			case 3:
				return "housewife-porn"
			case 4:
				return "teacher-porn"
			case 5:
				return "teen-porn"
			case 6:
				return "masturbation-porn"
			default:
				return "homemade-porn"
			}
		default:
			switch i {
			case 1:
				return "meitui"
			case 2:
				return "meixiong"
			case 3:
				return "meitun"
			case 4:
				return "shenyan"
			case 5:
				return "xiaoxinggan"
			case 6:
				return "xiaotianmei"
			default:
				return ""
			}

		}
	}
	if options.Output == "" {
		options.Output = "output/"
	}
	newScanner, err := scanner.NewScanner(&scanner.Options{
		Client: util.NewClient(options.Proxy, options.ProxyAuth),
		Rate:   options.Rate,
		Tag:    tagFunc(options.Type, options.Tag),
		Type:   options.Type,
		Limit:  options.Limit,
		Output: options.Output,
	})
	if err != nil {
		return nil, err
	}
	runner := &Runner{
		Options: options,
		Scanner: newScanner,
		Bar:     bar,
	}
	return runner, nil
}
