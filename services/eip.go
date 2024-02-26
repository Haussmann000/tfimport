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
	addresses []types.Address
}

func (v EipOutput) Describe() (addresses []types.Address, err error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}
	client := ec2.NewFromConfig(cfg)
	result, err := client.DescribeAddresses(context.TODO(), &ec2.DescribeAddressesInput{})
	resources := result.Addresses
	for _, resource := range resources {
		v.addresses = append(v.addresses, resource)
	}
	return v.addresses, err
}

func (v EipOutput) OutputFile(eips []types.Address) (result []lib.Result, err error) {
	results := []lib.Result{}
	for _, eip := range eips {
		result := lib.Result{}
		result.Id = eip.PublicIp
		results = append(results, result)
	}
	lib.OutputFile(lib.EIP_RESOUCE, results)
	return results, err
}
