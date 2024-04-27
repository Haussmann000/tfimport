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
	ResourceName string
	Vpcs         []types.Vpc
	lib.Output
}

func (v *Output) NewOutput(resource_name string) (*VpcOutput, error) {
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
		ResourceName: resource_name,
		Vpcs:         vpcs,
	}, err
}

func (v VpcOutput) OutputFile(resources []types.Vpc) (result []lib.Result, err error) {
	results := []lib.Result{}
	for _, resource := range resources {
		result := lib.Result{}
		result.Id = resource.VpcId
		results = append(results, result)
	}
	lib.OutputFile(v.ResourceName, results)
	return results, err
}

func (v VpcOutput) OutputTfFile(result []types.Vpc) (err error) {
	var tf lib.VpcTfOutput
	for i, r := range result {
		tf = lib.VpcTfOutput{CidrBlock: *r.CidrBlock, Index: i, Tags: r.Tags}
		err := lib.OutputTfFile(tf, lib.AWS_VPC_TMPL)
		if err != nil {
			return err
		}
	}
	return err
}
