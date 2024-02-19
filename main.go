package main

import (
	"fmt"

	lib "github.com/Haussmann000/tfimport/lib"
	service "github.com/Haussmann000/tfimport/services"
)

func main() {
	vpc := service.Vpc{}
	result, err := vpc.Describe()
	if err != nil {
		fmt.Println(err)
	}
	subnet := service.Subnet{}
	subnetresult, err := subnet.Describe()
	if err != nil {
		fmt.Println(err)
	}

	output := lib.Output{}
	output.OutputFile(lib.VPC_RESOUCE, result)
	output.OutputFile(lib.SUBNET_RESOUCE, subnetresult)

	if err != nil {
		fmt.Println(err)
	}

}
