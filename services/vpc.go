package service

import (
	"context"
	"log"

	"github.com/Haussmann000/tfimport/lib"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type VpcOutput struct {
	Resource_name string
	Vpcs          []types.Vpc
}

func (v VpcOutput) NewOutput(resource_name string) (*VpcOutput, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}
	client := ec2.NewFromConfig(cfg)
	result, err := client.DescribeVpcs(context.TODO(), &ec2.DescribeVpcsInput{})
	resources := result.Vpcs
	var vpcs []types.Vpc
	vpcs = append(vpcs, resources...)
	return &VpcOutput{
		Resource_name: resource_name,
		Vpcs:          vpcs,
	}, err
}

func (v VpcOutput) OutputFile(resources []types.Vpc) (result []lib.Result, err error) {
	results := []lib.Result{}
	for _, resource := range resources {
		result := lib.Result{}
		result.Id = resource.VpcId
		results = append(results, result)
	}
	lib.OutputFile(v.Resource_name, results)
	return results, err
}
