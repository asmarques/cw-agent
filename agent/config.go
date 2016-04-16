package agent

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// Config defines the configuration to use with an Agent
type Config struct {
	Region     string
	Hostname   string
	Namespace  string
	Interval   int64
	RunOnce    bool
	AllMetrics bool

	// Memory
	MemUnit                string
	MemSharedBuffersAsUsed bool
	MemAvailable           bool
	MemUsed                bool
	MemUtilization         bool
	SwapAvailable          bool
	SwapUsed               bool
	SwapUtilization        bool

	// Disk
	DiskUnit        string
	DiskAvailable   bool
	DiskUsed        bool
	DiskUtilization bool
	DiskPaths       string
}

type unit struct {
	name  string
	value float64
}

var validUnits = map[string]unit{
	"B":  {name: cloudwatch.StandardUnitBytes, value: 1},
	"KB": {name: cloudwatch.StandardUnitKilobytes, value: 1 << 10},
	"MB": {name: cloudwatch.StandardUnitMegabytes, value: 1 << 20},
	"GB": {name: cloudwatch.StandardUnitGigabytes, value: 1 << 30},
	"TB": {name: cloudwatch.StandardUnitTerabytes, value: 1 << 40},
}

func (c *Config) validate() error {
	var ok bool

	_, ok = validUnits[c.MemUnit]
	if !ok {
		return fmt.Errorf("invalid unit for memory reporting: %s", c.MemUnit)
	}

	_, ok = validUnits[c.DiskUnit]
	if !ok {
		return fmt.Errorf("invalid unit for disk reporting: %s", c.DiskUnit)
	}

	if !(c.hasMemMetricsSelected() || c.hasDiskMetricsSelected()) {
		return fmt.Errorf("no metrics selected")
	}

	return nil
}

func (c *Config) hasMemMetricsSelected() bool {
	return c.AllMetrics || c.MemAvailable || c.MemUsed || c.MemUtilization ||
		c.SwapAvailable || c.SwapUsed || c.SwapUtilization
}

func (c *Config) hasDiskMetricsSelected() bool {
	return c.AllMetrics || c.DiskAvailable || c.DiskUsed || c.DiskUtilization
}

func (c *Config) getDiskUnit() unit {
	return validUnits[c.DiskUnit]
}

func (c *Config) getMemoryUnit() unit {
	return validUnits[c.MemUnit]
}

func (c *Config) getDiskPaths() []string {
	return strings.Split(c.DiskPaths, ",")
}
