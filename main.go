package main

import (
	"fmt"

	"github.com/Haussmann000/tfimport/service"
)

func main() {
	result, err := service.DescribeMyVpcs()
	fmt.Println(result)
}
