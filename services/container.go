// container.go
package service

import (
	"github.com/aws/aws-sdk-go-v2/aws"
)

// Container holds references to services
type Container struct {
	AWSClient ServiceClient
}

// NewContainer creates a new container with provided AWS configuration
func NewContainer(cfg aws.Config) *Container {
	return &Container{
		AWSClient: NewAWSClient(cfg),
	}
}
