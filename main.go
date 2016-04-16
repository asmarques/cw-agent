package main

import (
	"flag"
	"log"

	"github.com/asmarques/cw-agent/agent"
)

var config agent.Config

func init() {
	flag.StringVar(&config.Region, "region", "", "name of the AWS region to publish metrics to (defaults to the region of the executing EC2 instance)")
	flag.StringVar(&config.Hostname, "hostname", "", "name used to identify the current host (defaults to the instance-id of the executing EC2 instance)")
	flag.StringVar(&config.Namespace, "namespace", "System/Linux", "namespace to use when publishing metrics")
	flag.Int64Var(&config.Interval, "interval", 5, "interval in minutes between publishing metrics")
	flag.BoolVar(&config.AllMetrics, "all-metrics", false, "reports all available memory and disk metrics")
	flag.BoolVar(&config.RunOnce, "once", false, "report metrics once and exit")

	// Memory metrics
	flag.StringVar(&config.MemUnit, "mem-unit", "MB", "specifies unit for memory and swap reporting: B, KB, MB, GB or TB")
	flag.BoolVar(&config.MemSharedBuffersAsUsed, "mem-shared-buffers-as-used", false, "reports shared and buffered memory as used memory")
	flag.BoolVar(&config.MemAvailable, "mem-avail", false, "reports available memory")
	flag.BoolVar(&config.MemUsed, "mem-used", false, "reports used memory")
	flag.BoolVar(&config.MemUtilization, "mem-util", false, "reports memory utilization percentage")
	flag.BoolVar(&config.SwapAvailable, "swap-avail", false, "reports swap available")
	flag.BoolVar(&config.SwapUsed, "swap-used", false, "reports swap used")
	flag.BoolVar(&config.SwapUtilization, "swap-util", false, "reports swap utilization percentage")

	// Disk metrics
	flag.StringVar(&config.DiskUnit, "disk-unit", "GB", "specifies unit for disk space reporting: B, KB, MB, GB or TB")
	flag.BoolVar(&config.DiskAvailable, "disk-avail", false, "reports disk space available")
	flag.BoolVar(&config.DiskUsed, "disk-used", false, "reports disk space used")
	flag.BoolVar(&config.DiskUtilization, "disk-util", false, "reports disk space utilization percentage")
	flag.StringVar(&config.DiskPaths, "disk-paths", "/", "specifies a comma separated list of paths to report disk space for")
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	a, err := agent.New(&config)
	if err != nil {
		log.Fatal(err)
	}

	err = a.Run()
	if err != nil {
		log.Fatal(err)
	}
}
