package runner

import "github.com/gosuri/uiprogress"

func (r *Runner) Run() error {
	uiprogress.Start() // 开始进度条
	go func() {
		for range r.Scanner.NumChan {
			r.Bar.Incr()
		}
	}()
	err := r.Scanner.Scan()
	if err != nil {
		return err
	}
	uiprogress.Stop()
	return nil
}
