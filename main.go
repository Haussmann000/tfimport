// main.go
package main

import (
	"context"
	"fmt"

	service "github.com/Haussmann000/tfimport/services"
	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	// Load AWS configuration from default sources
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Errorf("failed to load AWS config: %v", err))
	}

	// Initialize container
	container := service.NewContainer(cfg)
	s3, err := container.AWSClient.ListS3Buckets(context.Background())
	if err != nil {
		panic(fmt.Errorf("failed to fetch S3 buckets: %v", err))
	}
	ec2, err := container.AWSClient.GetEC2Instances(context.Background())
	if err != nil {
		panic(fmt.Errorf("failed to fetch ec2: %v", err))
	}

	fmt.Printf("%v\n", s3)
	fmt.Printf("%v\n", ec2)

	// Write generated Go code to file
	// result, err := container.
	// if err != nil {
	// 	panic(fmt.Errorf("failed to write to file: %v", err))
	// }
	// fmt.Println(result)

}
