// internal/aws/ec2/repository.go
package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// EC2RepositoryInterface はEC2リソースへのアクセスを抽象化します。
type EC2RepositoryInterface interface {
	DescribeVpcs(ctx context.Context, filters []types.Filter) ([]types.Vpc, error)
	DescribeSecurityGroups(ctx context.Context, groupIds []string) ([]types.SecurityGroup, error)
}

// EC2Repository はEC2RepositoryInterfaceを実装します。
type EC2Repository struct {
	client *ec2.Client
}

// NewEC2Repository は新しいEC2Repositoryを生成します。
func NewEC2Repository(client *ec2.Client) *EC2Repository {
	return &EC2Repository{
		client: client,
	}
}

// DescribeVpcs はAWSからVPCのリストを取得します。
func (r *EC2Repository) DescribeVpcs(ctx context.Context, filters []types.Filter) ([]types.Vpc, error) {
	input := &ec2.DescribeVpcsInput{
		Filters: filters,
	}
	result, err := r.client.DescribeVpcs(ctx, input)
	if err != nil {
		return nil, err
	}
	return result.Vpcs, nil
}

// DescribeSecurityGroups はAWSからSecurityGroupのリストを取得します。
func (r *EC2Repository) DescribeSecurityGroups(ctx context.Context, groupIds []string) ([]types.SecurityGroup, error) {
	input := &ec2.DescribeSecurityGroupsInput{
		GroupIds: groupIds,
	}
	result, err := r.client.DescribeSecurityGroups(ctx, input)
	if err != nil {
		return nil, err
	}
	return result.SecurityGroups, nil
} 
