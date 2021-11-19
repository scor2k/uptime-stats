package main

import (
	"fmt"
	"github.com/op/go-logging"
	"net"
	"os"
	"time"
)
import "github.com/urfave/cli/v2"
import "github.com/tatsushid/go-fastping"

var log = logging.MustGetLogger("uptimestats")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} ▶ %{level:.8s} %{shortfunc} %{color:reset} %{message}`,
)

func main() {
	appLogs := logging.NewLogBackend(os.Stderr, "", 0)
	appLogsFormatter := logging.NewBackendFormatter(appLogs, format)
	appLogsLeveled := logging.AddModuleLevel(appLogsFormatter)
	appLogsLeveled.SetLevel(logging.INFO, "")
	logging.SetBackend(appLogsLeveled)

	app := &cli.App{
		Name:     "uptimestats",
		HelpName: "uptimstats",
		Usage:    "check service availability and generate reports",
		Flags:    []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:  "monitor",
				Usage: "monitor host uptime",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "addr",
						Value:    "127.0.0.1",
						Usage:    "addr for monitoring",
						Required: true,
					},
					&cli.BoolFlag{
						Name:  "debug",
						Value: false,
						Usage: "Enable debug logging",
					},
				},
				Action: func(c *cli.Context) error {
					if c.Bool("debug") {
						appLogsLeveled.SetLevel(logging.DEBUG, "")
					}
					log.Debug("addr: ", c.String("addr"))

					ping := fastping.NewPinger()
					_ = ping.AddIP(c.String("addr"))
					ping.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
						fmt.Printf("[%v] ▶ IP Addr: %s receive, RTT: %v\n", time.Now().Format(time.RFC3339), addr.String(), rtt)
					}
					ping.OnIdle = func() {
						log.Debug("finish")
					}
					err := ping.Run()
					if err != nil {
						log.Error(err)
					}

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
