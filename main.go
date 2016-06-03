package main

import (
	"flag"
	"fmt"
	"github.com/domeos/alarm/cron"
	"github.com/domeos/alarm/g"
	"github.com/domeos/alarm/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	if *version {
		fmt.Println("falcon-alarm: ", g.FALCON_ALARM_VERSION)
		fmt.Println("domeos-alarm: ", g.DOMEOS_VERSION)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	g.ParseConfig(*cfg)
	g.InitRedisConnPool()
	g.InitDatabase()

	go http.Start()

	go cron.ReadHighEvent()
	go cron.ReadLowEvent()
	go cron.CombineSms()
	go cron.CombineMail()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		g.RedisConnPool.Close()
		os.Exit(0)
	}()

	select {}
}
