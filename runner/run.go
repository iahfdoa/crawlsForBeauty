package runner

import "github.com/gosuri/uiprogress"

func (r *Runner) Run() error {
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
