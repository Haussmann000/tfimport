// aws_client.go
package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// ServiceClient provides methods to interact with AWS services
type ServiceClient interface {
	GetEC2Instances(ctx context.Context) ([]string, error)
	ListS3Buckets(ctx context.Context) ([]string, error)
	ListSQSQueues(ctx context.Context) ([]string, error)
	ListDynamoDBTables(ctx context.Context) ([]string, error)
	ListCloudWatchAlarms(ctx context.Context) ([]string, error)
}

// AWSClient implements ServiceClient interface
type AWSClient struct {
	EC2Client        *ec2.Client
	S3Client         *s3.Client
	SQSClient        *sqs.Client
	DynamoDBClient   *dynamodb.Client
	CloudWatchClient *cloudwatch.Client
}

func NewAWSClient(cfg aws.Config) *AWSClient {
	return &AWSClient{
		EC2Client:        ec2.NewFromConfig(cfg),
		S3Client:         s3.NewFromConfig(cfg),
		SQSClient:        sqs.NewFromConfig(cfg),
		DynamoDBClient:   dynamodb.NewFromConfig(cfg),
		CloudWatchClient: cloudwatch.NewFromConfig(cfg),
	}
}

func (c *AWSClient) GetEC2Instances(ctx context.Context) ([]string, error) {
	resp, err := c.EC2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, err
	}

	var instanceIDs []string
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			instanceIDs = append(instanceIDs, *instance.InstanceId)
		}
	}
	return instanceIDs, nil
}

func (c *AWSClient) ListS3Buckets(ctx context.Context) ([]string, error) {
	resp, err := c.S3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	var bucketNames []string
	for _, bucket := range resp.Buckets {
		bucketNames = append(bucketNames, *bucket.Name)
	}
	return bucketNames, nil
}

func (c *AWSClient) ListSQSQueues(ctx context.Context) ([]string, error) {
	resp, err := c.SQSClient.ListQueues(ctx, &sqs.ListQueuesInput{})
	if err != nil {
		return nil, err
	}

	var queueURLs []string
	for _, url := range resp.QueueUrls {
		queueURLs = append(queueURLs, url)
	}
	return queueURLs, nil
}

func (c *AWSClient) ListDynamoDBTables(ctx context.Context) ([]string, error) {
	resp, err := c.DynamoDBClient.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		return nil, err
	}

	var tableNames []string
	for _, name := range resp.TableNames {
		tableNames = append(tableNames, name)
	}
	return tableNames, nil
}

func (c *AWSClient) ListCloudWatchAlarms(ctx context.Context) ([]string, error) {
	resp, err := c.CloudWatchClient.DescribeAlarms(ctx, &cloudwatch.DescribeAlarmsInput{})
	if err != nil {
		return nil, err
	}

	var alarmNames []string
	for _, alarm := range resp.MetricAlarms {
		alarmNames = append(alarmNames, *alarm.AlarmName)
	}
	return alarmNames, nil
}
