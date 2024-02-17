package service

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// ...

func DescribeMyVpcs() {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}
	client := ec2.NewFromConfig(cfg)
	result := client.DescribeVpcs(context.TODO(), &ec2.DescribeVpcsInput{})
	return result
}

type Vpc struct {
	VpcId string
}
