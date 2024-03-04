package lib

import "github.com/aws/aws-sdk-go-v2/service/ec2/types"

type TfOutput interface {
	OutputTfFile()
}

type VpcTfOutput struct {
	Index     int
	CidrBlock string
	Tags      []types.Tag
	TfOutput
}

type ImportBlock struct {
	Id *string `json:"id"`
	To string  `json:"to"`
}