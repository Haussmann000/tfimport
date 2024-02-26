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
	Resource_name string
	Subnets       []types.Subnet
}

func (v SubnetOutput) NewOutput(resource_name string) (*SubnetOutput, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}
	client := ec2.NewFromConfig(cfg)
	result, err := client.DescribeSubnets(context.TODO(), &ec2.DescribeSubnetsInput{})
	resources := result.Subnets
	var subnets []types.Subnet
	subnets = append(subnets, resources...)
	return &SubnetOutput{
		Resource_name: resource_name,
		Subnets:       subnets,
	}, err
}

func (v SubnetOutput) OutputFile(resources []types.Subnet) (result []lib.Result, err error) {
	results := []lib.Result{}
	for _, resource := range resources {
		result := lib.Result{}
		result.Id = resource.SubnetId
		results = append(results, result)
	}
	lib.OutputFile(v.Resource_name, results)
	return results, err
}
