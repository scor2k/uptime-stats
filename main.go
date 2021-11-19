package main

import "github.com/op/go-logging"
import "github.com/urfave/cli/v2"
import "os"

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
