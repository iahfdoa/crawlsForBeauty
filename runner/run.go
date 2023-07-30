package runner

import (
	"github.com/gosuri/uiprogress"
	"os"
)

func (r *Runner) Run() error {
	// api
	if r.Options.Api {
		r.ApiMain()
		os.Exit(0)
	}
	if !r.Options.Debug {
		uiprogress.Start()
	}
	go func() {
		for range r.Scanner.NumChan {
			r.Bar.Incr()
		}
	}()
	err := r.Scanner.Scan()
	if err != nil {
		return err
	}
	if !r.Options.Debug {
		uiprogress.Stop()
	}
	return nil
}
