package runner

import (
	"github.com/projectdiscovery/gologger"
	"time"
)

const (
	Version = "v1.1.0"
	Name    = "crawlsForBeauty"
	banner  = `
                            __     ______           ____                   __       
  ______________ __      __/ /____/ ____/___  _____/ __ )___  ____ ___  __/ /___  __
 / ___/ ___/ __ / | /| / / / ___/ /_  / __ \/ ___/ __  / _ \/ __ / / / / __/ / / /
/ /__/ /  / /_/ /| |/ |/ / (__  ) __/ / /_/ / /  / /_/ /  __/ /_/ / /_/ / /_/ /_/ / 
\___/_/   \__,_/ |__/|__/_/____/_/    \____/_/  /_____/\___/\__,_/\__,_/\__/\__, /    v1.1.0
                                                                           /____/   
`
)

// showBanner is used to show the banner to the user
func showBanner() {
	gologger.Print().Msgf("%s\n", banner)
	gologger.Print().Msgf("慎用。你要为自己的行为负责\n")
	gologger.Print().Msgf("开发者不承担任何责任，也不对任何误用或损坏负责.\n\n")
	time.Sleep(time.Second)
}
