package main

import (
	"fmt"

	service "github.com/Haussmann000/tfimport/services"
)

func main() {
	vpcs := service.VpcOutput{}
	result, err := vpcs.Describe()
	if err != nil {
		fmt.Println(err)
	}
	subnets := service.SubnetOutput{}
	subnetresult, err := subnets.Describe()
	if err != nil {
		fmt.Println(err)
	}

	vpcs.OutputFile(result)
	subnets.OutputFile(subnetresult)
	vpcs.OutputTfFile(result)

	if err != nil {
		fmt.Println(err)
	}

}
