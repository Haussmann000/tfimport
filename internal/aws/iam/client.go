package iam

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

func NewIAMClient(cfg aws.Config) *iam.Client {
	return iam.NewFromConfig(cfg)
} 
