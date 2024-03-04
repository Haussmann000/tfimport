package main

import (
	"fmt"

	"github.com/Haussmann000/tfimport/lib"
	service "github.com/Haussmann000/tfimport/services"
)

func main() {
	result, err := DescribeServices()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func DescribeServices() (result []lib.Result, err error) {
	vpc, err := service.VpcOutput.NewOutput(service.VpcOutput{}, lib.VPC_RESOUCE)

	if err != nil {
		return nil, err
	}

	_, err = vpc.OutputFile(vpc.Vpcs)
	if err != nil {
		return nil, err
	}
	err = vpc.OutputTfFile(vpc.Vpcs)
	if err != nil {
		return nil, err
	}

	subnet, err := service.SubnetOutput.NewOutput(service.SubnetOutput{}, lib.SUBNET_RESOUCE)
	if err != nil {
		return nil, err
	}
	_, err = subnet.OutputFile(subnet.Subnets)
	if err != nil {
		return nil, err
	}
	eip, err := service.EipOutput.NewOutput(service.EipOutput{}, lib.EIP_RESOUCE)
	if err != nil {
		return nil, err
	}
	_, err = eip.OutputFile(eip.Eips)
	if err != nil {
		return nil, err
	}
	return result, err
}
