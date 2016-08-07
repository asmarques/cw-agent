package agent

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func getAutoScalingGroup(awsSession *session.Session, instanceID string) (string, error) {
	svc := ec2.New(awsSession)
	params := &ec2.DescribeTagsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("resource-id"),
				Values: []*string{aws.String(instanceID)},
			},
			{
				Name:   aws.String("key"),
				Values: []*string{aws.String("aws:autoscaling:groupName")},
			},
		},
	}

	resp, err := svc.DescribeTags(params)
	if err != nil {
		return "", err
	}

	if len(resp.Tags) == 0 {
		return "", fmt.Errorf("no autoscaling group found for instance %s", instanceID)
	}

	return *resp.Tags[0].Value, nil
}

func getInstanceID(awsSession *session.Session) (string, error) {
	svc := ec2metadata.New(awsSession)
	instanceID, err := svc.GetMetadata("instance-id")
	if err != nil {
		return "", fmt.Errorf("unable to retrieve instance ID from metadata service: %s", err)
	}
	return instanceID, nil
}

func getRegion(awsSession *session.Session) (string, error) {
	svc := ec2metadata.New(awsSession)
	region, err := svc.Region()
	if err != nil {
		return "", fmt.Errorf("unable to determine region from metadata service: %s", err)
	}
	return region, nil
}
