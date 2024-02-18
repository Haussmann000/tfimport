package main

import (
	"fmt"

	lib "github.com/Haussmann000/tfimport/lib"
	service "github.com/Haussmann000/tfimport/services"
)

func main() {
	result, err := service.DescribeMyVpcs()
	lib.OutputFile(result)
	if err != nil {
		fmt.Println(err)
	}

}
