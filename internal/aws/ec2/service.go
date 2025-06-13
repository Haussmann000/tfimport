// internal/aws/ec2/service.go
package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// Vpc はHCL生成に必要なVPCの情報を保持します。
type Vpc struct {
	ID        string
	CidrBlock string
	Tags      map[string]string
}

// EC2ServiceInterface はEC2関連のビジネスロジックを定義します。
type EC2ServiceInterface interface {
	ListVpcs(ctx context.Context, resourceName string) ([]Vpc, error)
}

// EC2Service はEC2ServiceInterfaceを実装します。
type EC2Service struct {
	repo EC2RepositoryInterface
}

// NewEC2Service は新しいEC2Serviceを生成します。
func NewEC2Service(repo EC2RepositoryInterface) *EC2Service {
	return &EC2Service{
		repo: repo,
	}
}

// ListVpcs はVPCのリストを取得し、ドメインオブジェクトに変換します。
func (s *EC2Service) ListVpcs(ctx context.Context, resourceName string) ([]Vpc, error) {
	var filters []types.Filter
	if resourceName != "" {
		filters = append(filters, types.Filter{
			Name:   aws.String("tag:Name"),
			Values: []string{resourceName + "*"},
		})
	}

	awsVpcs, err := s.repo.DescribeVpcs(ctx, filters)
	if err != nil {
		return nil, err
	}

	var vpcs []Vpc
	for _, v := range awsVpcs {
		vpcs = append(vpcs, Vpc{
			ID:        *v.VpcId,
			CidrBlock: *v.CidrBlock,
			Tags:      convertTags(v.Tags),
		})
	}
	return vpcs, nil
}

func convertTags(tags []types.Tag) map[string]string {
	m := make(map[string]string)
	for _, t := range tags {
		m[*t.Key] = *t.Value
	}
	return m
} 
