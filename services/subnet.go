package service

import (
	"context"
	"log"

	"github.com/Haussmann000/tfimport/lib"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type SubnetOutput struct {
	subnets []types.Subnet
}

func (s SubnetOutput) Describe() (subnets []types.Subnet, err error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}
	client := ec2.NewFromConfig(cfg)
	result, err := client.DescribeSubnets(context.TODO(), &ec2.DescribeSubnetsInput{})
	resources := result.Subnets
	for _, resource := range resources {
		s.subnets = append(s.subnets, resource)
	}
	return s.subnets, err
}

func (v SubnetOutput) OutputFile(subnets []types.Subnet) (result []lib.Result, err error) {
	results := []lib.Result{}
	for _, subnet := range subnets {
		result := lib.Result{}
		result.Id = subnet.SubnetId
		results = append(results, result)
	}
	lib.OutputFile(lib.SUBNET_RESOUCE, results)
	return results, err
}
