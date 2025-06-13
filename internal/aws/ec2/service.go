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

// Service はEC2関連のビジネスロジックを定義します。
type Service interface {
	ListVpcs(ctx context.Context, resourceName string) ([]Vpc, error)
}

// VPCService はServiceを実装します。
type VPCService struct {
	repo EC2RepositoryInterface
}

// NewVPCService は新しいVPCServiceを生成します。
func NewVPCService(repo EC2RepositoryInterface) *VPCService {
	return &VPCService{
		repo: repo,
	}
}

// ListVpcs はVPCのリストを取得し、ドメインオブジェクトに変換します。
func (s *VPCService) ListVpcs(ctx context.Context, resourceName string) ([]Vpc, error) {
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
