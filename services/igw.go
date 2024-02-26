package service

import (
	"context"
	"log"

	"github.com/Haussmann000/tfimport/lib"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type IgwOutput struct {
	Resource_name string
	Igws          []types.InternetGateway
}

func (v IgwOutput) NewOutput(resource_name string) (*IgwOutput, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}
	client := ec2.NewFromConfig(cfg)
	result, err := client.DescribeInternetGateways(context.TODO(), &ec2.DescribeInternetGatewaysInput{})
	resources := result.InternetGateways
	var igws []types.InternetGateway
	igws = append(igws, resources...)
	return &IgwOutput{
		Resource_name: resource_name,
		Igws:          igws,
	}, err
}

func (v IgwOutput) OutputFile(resources []types.InternetGateway) (result []lib.Result, err error) {
	results := []lib.Result{}
	for _, resource := range resources {
		result := lib.Result{}
		result.Id = resource.InternetGatewayId
		results = append(results, result)
	}
	lib.OutputFile(v.Resource_name, results)
	return results, err
}
