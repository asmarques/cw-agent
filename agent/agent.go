package agent

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// Agent is used to retrieve and report metrics from the current host to CloudWatch
type Agent struct {
	config          *Config
	extraDimensions map[string]string
	svc             *cloudwatch.CloudWatch
}

const (
	awsRegionEnv           = "AWS_REGION"
	instanceIDMetadataKey  = "instance-id"
	instanceIDDimensionKey = "InstanceId"
	hostnameDimensionKey   = "Hostname"
)

// New creates a new Agent based on the supplied configuration
func New(config *Config) (*Agent, error) {
	err := config.validate()
	if err != nil {
		return nil, fmt.Errorf("configuration error: %s", err)
	}

	awsConfig := &aws.Config{}
	awsSession := session.New()
	svcMetadata := ec2metadata.New(awsSession)

	// Try to determine the AWS region for CloudWatch using the following sources:
	// - Supplied Agent configuration
	// - Environment variable (AWS_REGION)
	// - EC2 metadata of current host
	if config.Region != "" {
		awsConfig.Region = aws.String(config.Region)
	} else if region := os.Getenv(awsRegionEnv); region != "" {
		awsConfig.Region = aws.String(region)
	} else {
		region, err := svcMetadata.Region()
		if err != nil {
			return nil, fmt.Errorf("unable to determine region from metadata service: %s", err)
		}
		awsConfig.Region = aws.String(region)
	}

	extraDimensions := make(map[string]string)

	// Try to determine the name of the executing host to associate with the reported metrics
	if config.Hostname != "" {
		extraDimensions[hostnameDimensionKey] = config.Hostname
	} else {
		instanceID, err := svcMetadata.GetMetadata(instanceIDMetadataKey)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve key %s from metadata service: %s",
				instanceIDMetadataKey, err)
		}
		extraDimensions[instanceIDDimensionKey] = instanceID
	}

	return &Agent{
		config:          config,
		extraDimensions: extraDimensions,
		svc:             cloudwatch.New(awsSession, awsConfig),
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
