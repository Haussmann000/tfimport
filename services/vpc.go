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
	vpcs []types.Vpc
}

func (v VpcOutput) Describe() (vpcs []types.Vpc, err error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}
	client := ec2.NewFromConfig(cfg)
	result, err := client.DescribeVpcs(context.TODO(), &ec2.DescribeVpcsInput{})
	resources := result.Vpcs
	for _, resource := range resources {
		v.vpcs = append(v.vpcs, resource)
	}
	return v.vpcs, err
}

func (v VpcOutput) OutputFile(vpcs []types.Vpc) (result []lib.Result, err error) {
	results := []lib.Result{}
	for _, vpc := range vpcs {
		result := lib.Result{}
		result.Id = vpc.VpcId
		results = append(results, result)
	}
	lib.OutputFile(lib.VPC_RESOUCE, results)
	return results, err
}

func (v VpcOutput) OutputTfFile(vpcs []types.Vpc) error {
	lib.OutputTfFile(vpcs)
	return nil
}
