package agent

import (
	"time"

	"math"

	"github.com/asmarques/cw-agent/stats"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

const (
	dimensionMountPath         = "MountPath"
	metricMemoryAvailable      = "MemoryAvailable"
	metricMemoryUsed           = "MemoryUsed"
	metricMemoryUtilization    = "MemoryUtilization"
	metricSwapAvailable        = "SwapAvailable"
	metricSwapUsed             = "SwapUsed"
	metricSwapUtilization      = "SwapUtilization"
	metricDiskSpaceAvailable   = "DiskSpaceAvailable"
	metricDiskSpaceUsed        = "DiskSpaceUsed"
	metricDiskSpaceUtilization = "DiskSpaceUtilization"
	metricDataMaxSize          = 20
)

type metric struct {
	name       string
	unit       string
	value      float64
	dimensions map[string]string
}

func (a *Agent) putMetrics() error {
	now := time.Now()

	diskMetrics, err := a.getDiskMetrics()
	if err != nil {
		return err
	}

	memMetrics, err := a.getMemoryMetrics()
	if err != nil {
		return err
	}

	var metricData []*cloudwatch.MetricDatum
	var dimensions []*cloudwatch.Dimension

	for key, value := range a.baseDimensions {
		dimensions = append(dimensions, &cloudwatch.Dimension{
			Name:  aws.String(key),
			Value: aws.String(value),
		})
	}

	for _, metric := range append(diskMetrics, memMetrics...) {
		if math.IsNaN(metric.value) {
			continue
		}

		var extraDimensions []*cloudwatch.Dimension
		for name, value := range metric.dimensions {
			dimensions = append(dimensions, &cloudwatch.Dimension{
				Name:  aws.String(name),
				Value: aws.String(value),
			})
		}

		for _, dimension := range dimensions {
			metricData = append(metricData,
				&cloudwatch.MetricDatum{
					MetricName: aws.String(metric.name),
					Dimensions: append(extraDimensions, dimension),
					Timestamp:  aws.Time(now),
					Unit:       aws.String(metric.unit),
					Value:      aws.Float64(metric.value),
				})
		}
	}

	// Split metric data accross multiple requests according to the API limits
	size := len(metricData)
	for i := 0; i < size; i += metricDataMaxSize {
		_, err = a.svc.PutMetricData(&cloudwatch.PutMetricDataInput{
			MetricData: metricData[i:min(i+metricDataMaxSize, size)],
			Namespace:  aws.String(a.config.Namespace),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Agent) getMemoryMetrics() ([]*metric, error) {
	metrics := []*metric{}
	memUnit := a.config.getMemoryUnit()

	if a.config.hasMemMetricsSelected() {
		mem, err := stats.GetMemoryStats()
		if err != nil {
			return nil, err
		}

		if a.config.AllMetrics || a.config.MemAvailable {
			metrics = append(metrics, &metric{
				name:  metricMemoryAvailable,
				unit:  memUnit.name,
				value: float64(mem.MemoryAvailable) / memUnit.value,
			})
		}

		if a.config.AllMetrics || a.config.MemUsed {
			metrics = append(metrics, &metric{
				name:  metricMemoryUsed,
				unit:  memUnit.name,
				value: float64(mem.MemoryUsed) / memUnit.value,
			})
		}

		if a.config.AllMetrics || a.config.MemUtilization {
			metrics = append(metrics, &metric{
				name:  metricMemoryUtilization,
				unit:  cloudwatch.StandardUnitPercent,
				value: float64(mem.MemoryUtilization),
			})
		}

		if a.config.AllMetrics || a.config.SwapAvailable {
			metrics = append(metrics, &metric{
				name:  metricSwapAvailable,
				unit:  memUnit.name,
				value: float64(mem.SwapAvailable) / memUnit.value,
			})
		}

		if a.config.AllMetrics || a.config.SwapUsed {
			metrics = append(metrics, &metric{
				name:  metricSwapUsed,
				unit:  memUnit.name,
				value: float64(mem.SwapUsed) / memUnit.value,
			})
		}

		if a.config.AllMetrics || a.config.SwapUtilization {
			metrics = append(metrics, &metric{
				name:  metricSwapUtilization,
				unit:  cloudwatch.StandardUnitPercent,
				value: float64(mem.SwapUtilization),
			})
		}
	}

	return metrics, nil
}

func (a *Agent) getDiskMetrics() ([]*metric, error) {
	metrics := []*metric{}
	diskUnit := a.config.getDiskUnit()

	if a.config.hasDiskMetricsSelected() {
		paths := a.config.getDiskPaths()

		for _, path := range paths {
			disk, err := stats.GetDiskStats(path)
			if err != nil {
				return nil, err
			}

			dimensions := map[string]string{dimensionMountPath: disk.MountPath}

			if a.config.AllMetrics || a.config.DiskAvailable {
				metrics = append(metrics, &metric{
					name:       metricDiskSpaceAvailable,
					unit:       diskUnit.name,
					value:      float64(disk.DiskSpaceAvailable) / diskUnit.value,
					dimensions: dimensions,
				})
			}

			if a.config.AllMetrics || a.config.DiskUsed {
				metrics = append(metrics, &metric{
					name:       metricDiskSpaceUsed,
					unit:       diskUnit.name,
					value:      float64(disk.DiskSpaceUsed) / diskUnit.value,
					dimensions: dimensions,
				})
			}

			if a.config.AllMetrics || a.config.DiskUtilization {
				metrics = append(metrics, &metric{
					name:       metricDiskSpaceUtilization,
					unit:       cloudwatch.StandardUnitPercent,
					value:      float64(disk.DiskSpaceUtilization),
					dimensions: dimensions,
				})
			}
		}
	}

	return metrics, nil
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
