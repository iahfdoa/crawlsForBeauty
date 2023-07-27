package scanner

import (
	"context"
	"fmt"
	"github.com/iahfdoa/crawlsForBeauty/util"
	"github.com/projectdiscovery/ratelimit"
	"github.com/remeh/sizedwaitgroup"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Options struct {
	Client *http.Client
	Rate   int
	Tag    string
	Url    string
	Limit  int
	Output string
}

type Scanner struct {
	client      *http.Client
	request     map[string]string
	imgHtmlChan chan string
	imgUrlChan  chan string
	output      string
	limiter     *ratelimit.Limiter
	wg          *sizedwaitgroup.SizedWaitGroup
	lock        sync.Mutex
	NumChan     chan struct{}
}

func NewScanner(options *Options) (*Scanner, error) {
	if options.Url == "" {
		return nil, fmt.Errorf("错误的目标: url、xpathExpr 不能为空")
	}
	if options.Output != "" {
		err := util.CreateDir(options.Output)
		if err != nil {

			return nil, err
		}
	}
	limiter := ratelimit.New(context.Background(), uint(options.Rate), time.Duration(1)*time.Second)
	wg := new(sizedwaitgroup.SizedWaitGroup)
	*wg = sizedwaitgroup.New(options.Rate)
	// 初始化
	scan := &Scanner{}
	request := make(map[string]string)
	request["tag"] = options.Tag
	request["url"] = options.Url
	limit := strconv.Itoa(options.Limit)
	request["limit"] = limit
	scan.client = options.Client
	scan.request = request
	scan.imgHtmlChan = make(chan string)
	scan.imgUrlChan = make(chan string)
	scan.NumChan = make(chan struct{})
	scan.output = options.Output
	scan.limiter = limiter
	scan.wg = wg
	return scan, nil
}
