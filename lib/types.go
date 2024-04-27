package lib

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type TfOutput interface {
	OutputTfFile()
}

type Output interface {
	NewOutput() any
	OutputFile() Result
	OutputTfFile() TfOutput
}

func NewOutput(output Output, resource string) Output {
	return output
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
