package service

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func DescribeMyVpcs() (output VpcResults, err error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}
	client := ec2.NewFromConfig(cfg)
	result, err := client.DescribeVpcs(context.TODO(), &ec2.DescribeVpcsInput{})
	vpcs := result.Vpcs
	output = OutputResult(vpcs)
	return output, err
}

func OutputResult(vpcs []types.Vpc) (VpcResults VpcResults) {
	vpcresults := VpcResults
	for _, vpc := range vpcs {
		vpcresults = append(vpcresults, vpc)
	}
	return vpcresults
}

type ImportBlock struct {
	Id string
	To string
}

// type VpcResult struct {
// 	VpcId     string
// 	CidrBlock string
// }

type VpcResults []types.Vpc
