package service

import (
	"context"
	"log"

	"github.com/Haussmann000/tfimport/lib"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EipOutput struct {
	Resource_name string
	Eips          []types.Address
}

func (v EipOutput) NewOutput(resource_name string) (*EipOutput, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}
	client := ec2.NewFromConfig(cfg)
	result, err := client.DescribeAddresses(context.TODO(), &ec2.DescribeAddressesInput{})
	resources := result.Addresses
	var subnets []types.Address
	subnets = append(subnets, resources...)
	return &EipOutput{
		Resource_name: resource_name,
		Eips:          subnets,
	}, err
}

func (v EipOutput) OutputFile(resources []types.Address) (result []lib.Result, err error) {
	results := []lib.Result{}
	for _, resource := range resources {
		result := lib.Result{}
		result.Id = resource.PublicIp
		results = append(results, result)
	}
	lib.OutputFile(v.Resource_name, results)
	return results, err
}
