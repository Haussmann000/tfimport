package service

import (
	"context"
	"log"

	"github.com/Haussmann000/tfimport/lib"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type Subnet struct {
	tf_resource_name string
	tf_resouce_id    string
}

func (s Subnet) Describe() (output []lib.Result, err error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}
	client := ec2.NewFromConfig(cfg)
	result, err := client.DescribeSubnets(context.TODO(), &ec2.DescribeSubnetsInput{})
	subnets := result.Subnets
	results := []lib.Result{}
	for _, subnet := range subnets {
		result := lib.Result{}
		result.Id = subnet.SubnetId
		results = append(results, result)
	}

	return results, err
}
