package runner

import (
	"encoding/json"
	"fmt"
	"github.com/iahfdoa/crawlsForBeauty/util"
	"github.com/iahfdoa/crawlsForBeauty/web/dao"
	"github.com/iahfdoa/crawlsForBeauty/web/router"
	"github.com/robfig/cron"
	"golang.org/x/exp/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

const imagesJson = `config/config.json`

var mutex sync.Mutex

func (r *Runner) ApiMain() {
	go func() {
		r.Scanner.SetType(r.Options.Type)
		r.Scanner.SetOutput(filepath.Join(r.Options.ApiOutput, fmt.Sprintf("%d", r.Options.Type)))
		//uiprogress.Start() // 开始进度条
		go func() {
			for range r.Scanner.NumChan {
				r.Bar.Incr()
			}
		}()
		err := r.Scanner.Scan()
		if err != nil {
			panic(err)
		}

		//uiprogress.Stop()
		r.InitConfig()
	}()
	r.InitConfig()
	c := cron.New()
	dao.NewOptions(r.Options.Type)
	c.AddFunc("@every 1h", func() {
		err := util.CreateDir(filepath.Join(r.Options.ApiOutput, strconv.Itoa(r.Options.Type)))
		if err != nil {
			return
		}
		err = util.CreateDir("config")
		if err != nil {
			return
		}
		_, err = os.Create("config/config.json")
		if err != nil {
			return
		}

		// 使用互斥锁更新 imagePaths 切片
		mutex.Lock()
		defer mutex.Unlock()

		imagePaths0, err := util.GetAllFilesInDir(filepath.Join(r.Options.ApiOutput, strconv.Itoa(0)))
		if err != nil {
			return
		}
		imagePaths1, err := util.GetAllFilesInDir(filepath.Join(r.Options.ApiOutput, strconv.Itoa(1)))
		if err != nil {
			return
		}
		imagePaths2, err := util.GetAllFilesInDir(filepath.Join(r.Options.ApiOutput, strconv.Itoa(2)))
		if err != nil {
			return
		}
		m := make(map[string][]string)
		m["0"] = imagePaths0
		m["1"] = imagePaths1
		m["2"] = imagePaths2
		data, err := json.Marshal(m)
		if err != nil {
			return

		}
		err = os.WriteFile(imagesJson, data, 0644)
		if err != nil {
			return

		}
	}) // 每隔1小时更新一次
	c.AddFunc("@every 5h", func() {
		rand.Seed(uint64(time.Now().UnixNano())) // 设置随机数种子，以保证每次生成的随机数都不同

		// 生成一个范围在0到2之间的随机整数
		choice := rand.Intn(3)
		r.Scanner.SetType(choice)
		r.Scanner.SetOutput(filepath.Join(r.Options.ApiOutput, fmt.Sprintf("%d", choice)))
		//uiprogress.Start() // 开始进度条
		go func() {
			for range r.Scanner.NumChan {
				r.Bar.Incr()
			}
		}()
		err := r.Scanner.Scan()
		if err != nil {
			return
		}

		//uiprogress.Stop()
		r.InitConfig()
	})
	c.Start()
	ro := router.Router()

	err := ro.Run(fmt.Sprintf(":%d", r.Options.ApiPort))
	if err != nil {
		panic(err)
	}
}

func (r *Runner) InitConfig() {
	err := util.CreateDir(filepath.Join(r.Options.ApiOutput, strconv.Itoa(r.Options.Type)))
	if err != nil {
		panic(err)
	}
	err = util.CreateDir("config")
	if err != nil {
		panic(err)
	}
	_, err = os.Create("config/config.json")
	if err != nil {
		panic(err)
	}

	// 使用互斥锁更新 imagePaths 切片
	mutex.Lock()
	defer mutex.Unlock()

	// 清空原有的 imagePaths 切片
	imagePaths0, err := util.GetAllFilesInDir(filepath.Join(r.Options.ApiOutput, strconv.Itoa(0)))
	if err != nil {
		return
	}
	imagePaths1, err := util.GetAllFilesInDir(filepath.Join(r.Options.ApiOutput, strconv.Itoa(1)))
	if err != nil {
		return
	}
	imagePaths2, err := util.GetAllFilesInDir(filepath.Join(r.Options.ApiOutput, strconv.Itoa(2)))
	if err != nil {
		return
	}
	m := make(map[string][]string)
	m["0"] = imagePaths0
	m["1"] = imagePaths1
	m["2"] = imagePaths2
	data, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal imagePaths: %s", err))

	}
	err = os.WriteFile(imagesJson, data, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to save imagePaths to config file: %s", err))

	}

}
