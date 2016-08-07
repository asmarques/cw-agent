package agent

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// Agent is used to retrieve and report metrics from the current host to CloudWatch
type Agent struct {
	config         *Config
	baseDimensions map[string]string
	svc            *cloudwatch.CloudWatch
}

const (
	awsRegionEnv                 = "AWS_REGION"
	instanceIDDimensionKey       = "InstanceId"
	autoScalingGroupDimensionKey = "AutoScalingGroupName"
	hostnameDimensionKey         = "Hostname"
)

// New creates a new Agent based on the supplied configuration
func New(config *Config) (*Agent, error) {
	err := config.validate()
	if err != nil {
		return nil, fmt.Errorf("configuration error: %s", err)
	}

	awsSession := session.New()

	// Try to determine the AWS region using the following sources:
	// - Supplied Agent configuration
	// - Environment variable (AWS_REGION)
	// - EC2 metadata of current host
	if config.Region != "" {
		awsSession.Config.Region = aws.String(config.Region)
	} else if region := os.Getenv(awsRegionEnv); region != "" {
		awsSession.Config.Region = aws.String(region)
	} else {
		region, err := getRegion(awsSession)
		if err != nil {
			return nil, err
		}
		awsSession.Config.Region = aws.String(region)
	}

	// Configure dimensions to aggregate by when reporting metrics
	dimensions := make(map[string]string)
	if config.Hostname != "" {
		dimensions[hostnameDimensionKey] = config.Hostname
	} else {
		instanceID, err := getInstanceID(awsSession)
		if err != nil {
			return nil, err
		}
		dimensions[instanceIDDimensionKey] = instanceID

		if config.AutoScaling {
			autoScalingGroup, err := getAutoScalingGroup(awsSession, instanceID)
			if err != nil {
				return nil, err
			}
			dimensions[autoScalingGroupDimensionKey] = autoScalingGroup
		}
	}

	return &Agent{
		config:         config,
		baseDimensions: dimensions,
		svc:            cloudwatch.New(awsSession),
	}, nil
}

// Run executes the Agent to report metrics once or according to a predefined interval
func (a *Agent) Run() error {
	err := a.putMetrics()
	if err != nil {
		return err
	}

	if a.config.RunOnce {
		return nil
	}

	metricsTicker := time.NewTicker(time.Duration(a.config.Interval) * time.Minute)

	for {
		select {
		case <-metricsTicker.C:
			err = a.putMetrics()
			if err != nil {
				return err
			}
		}
	}
}
