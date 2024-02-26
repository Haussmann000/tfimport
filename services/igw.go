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
	igws []types.InternetGateway
}

func (v IgwOutput) Describe() (vpcs []types.InternetGateway, err error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}
	client := ec2.NewFromConfig(cfg)
	result, err := client.DescribeInternetGateways(context.TODO(), &ec2.DescribeInternetGatewaysInput{})
	resources := result.InternetGateways
	for _, resource := range resources {
		v.igws = append(v.igws, resource)
	}
	return v.igws, err
}

func (v IgwOutput) OutputFile(igws []types.InternetGateway) (result []lib.Result, err error) {
	results := []lib.Result{}
	for _, igw := range igws {
		result := lib.Result{}
		result.Id = igw.InternetGatewayId
		results = append(results, result)
	}
	lib.OutputFile(lib.IGW_RESOUCE, results)
	return results, err
}
