package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/bcolucci/zerohash/pkg/monitoring"
	"github.com/bcolucci/zerohash/pkg/producer"
)

var runningMode string
var verbose bool

var productIDs = []string{"BTC-USD", "ETH-USD", "ETH-BTC"}

func main() {
	parseEnvVars()

	if !verbose {
		log.SetOutput(ioutil.Discard)
	}

	ctx, cancel := context.WithCancel(context.Background())

	updates, err := producer.ProduceCoinbaseUpdates(ctx, productIDs)
	if err != nil {
		log.Fatal(err)
	}

	points := producer.ProduceVWAPUpdates(ctx, updates, 200)

	if runningMode == "console" {
		time.AfterFunc(30*time.Second, cancel)
		for {
			point, open := <-points
			if !open {
				return
			}
			fmt.Printf("%s: %.4f\t(%d pts)\n", point.ProductId, point.VWAP, point.NbPoints)
		}
	}

	if runningMode == "server" {
		monitoring.StartMonitoringServer(points)
	}
}

func parseEnvVars() {
	runningMode = os.Getenv("RUNNING_MODE")
	if runningMode == "" {
		runningMode = "console"
	}
	if runningMode != "console" && runningMode != "server" {
		log.Panic("unknown running mode")
	}

	if os.Getenv("VERBOSE") != "" {
		verbose = true
	}
}
