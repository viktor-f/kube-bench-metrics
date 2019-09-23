package main

import (
	"os"
	"time"

	"github.com/ibrokethecloud/kube-bench-metrics/metrics"
	"github.com/ibrokethecloud/kube-bench-metrics/wrapper"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// VERSION refers to the current build id.
// Will be override by commit hash from the ld directives.
var VERSION = "0.0.1"

// main code execution //
func main() {
	logrus.SetLevel(logrus.DebugLevel)

	app := cli.NewApp()
	app.Name = "wrapper"
	app.Version = VERSION
	app.Action = runWrapper
	app.HideVersion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "check, c",
			Value: "",
			Usage: "A comma delimited list of checks to run",
		},
		cli.StringFlag{
			Name:  "versionOverride",
			Value: "",
			Usage: "Version of checks to run",
		},
		cli.IntFlag{
			Name:  "delay, d",
			Value: 60,
			Usage: "Delay in minutes between scheduled runs",
		},
		cli.StringFlag{
			Name:  "nodeType, n",
			Value: "",
			Usage: "Node classification: Possible values are node and master",
		},
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func runWrapper(ctx *cli.Context) (err error) {

	// Need to organize this.
	// For now just parse the yaml and generate the metrics
	delay := ctx.Int("delay")
	nodeType := ctx.String("nodeType")
	// Run kube-bench cli in a predefined interval
	go func() {
		for {
			w := wrapper.NewWrapper(nodeType)
			for _, node := range w.NodeType {
				err := w.RunBenchMarking(ctx, node)
				if err != nil {
					logrus.Error(err)
				}
				metrics.GenerateMetrics(*w)
			}

			// Configurable value of time.
			// Else default of 60 mins will
			// be used.
			time.Sleep(time.Duration(delay) * time.Minute)
		}
	}()

	metrics.ServeMetrics()
	return nil
}