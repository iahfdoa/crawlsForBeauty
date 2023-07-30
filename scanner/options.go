package scanner

import (
	"context"
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
	Limit  int
	Output string
	Type   int
}

type Scanner struct {
	client          *http.Client
	request         map[string]interface{}
	imgHtmlChan     chan string
	imgUrlChan      chan string
	output          string
	limiter         *ratelimit.Limiter
	wg              *sizedwaitgroup.SizedWaitGroup
	lock            sync.Mutex
	NumChan         chan struct{}
	UrlCallBackFunc func(str string) (string, error)
}

func NewScanner(options *Options) (*Scanner, error) {
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
	request := make(map[string]interface{})
	limit := strconv.Itoa(options.Limit)
	scan.client = options.Client
	scan.request = request
	scan.imgHtmlChan = make(chan string)
	scan.imgUrlChan = make(chan string)
	scan.NumChan = make(chan struct{})
	scan.output = options.Output
	scan.limiter = limiter
	scan.wg = wg
	request["limit"] = limit
	request["tag"] = options.Tag
	request["type"] = options.Type

	return scan, nil
}

func (s *Scanner) SetType(t int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.request["type"] = t

}

func (s *Scanner) SetOutput(output string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.output = output
}
