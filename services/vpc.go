package service

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// ...

func DescribeMyVpcs() (vpcs []types.Vpc, err error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}
	client := ec2.NewFromConfig(cfg)
	result, err := client.DescribeVpcs(context.TODO(), &ec2.DescribeVpcsInput{})
	vpcs = result.Vpcs

	return vpcs, err
}

type VpcResult struct {
	VpcId string
}
