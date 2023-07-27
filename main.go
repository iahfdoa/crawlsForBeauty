package main

import (
	"fmt"
	"github.com/iahfdoa/crawlsForBeauty/runner"
	"os"
	"os/signal"
)

func main() {
	run, err := runner.NewRunner(runner.ParserOptions())
	if err != nil {
		os.Exit(0)
	}
	if run == nil {
		os.Exit(0)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			os.Exit(1)
		}
	}()
	err = run.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
