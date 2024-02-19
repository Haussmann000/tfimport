package service

import (
	"context"
	"log"

	"github.com/Haussmann000/tfimport/lib"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type Vpc struct {
	tf_resource_name string
	tf_resouce_id    string
}

func (v Vpc) Describe() (output []lib.Result, err error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}
	client := ec2.NewFromConfig(cfg)
	result, err := client.DescribeVpcs(context.TODO(), &ec2.DescribeVpcsInput{})
	vpcs := result.Vpcs
	results := []lib.Result{}
	for _, subnet := range vpcs {
		result := lib.Result{}
		result.Id = subnet.VpcId
		results = append(results, result)
	}
	return results, err
}
