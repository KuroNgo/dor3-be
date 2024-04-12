package google

import "github.com/robfig/cron"

func init() {
	Cron.Start()
}

var Cron = cron.New()
