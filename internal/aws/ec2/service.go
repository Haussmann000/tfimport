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

// SecurityGroup はHCL生成に必要なセキュリティグループの情報を保持します。
type SecurityGroup struct {
	ID          string
	Name        string
	Description string
	Tags        map[string]string
}

// Service はEC2関連のビジネスロジックを定義します。
type Service interface {
	ListVpcs(ctx context.Context, resourceName string) ([]Vpc, error)
	ListSecurityGroups(ctx context.Context, groupIDs []string) ([]SecurityGroup, error)
}

// EC2Service はServiceを実装します。
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

// ListSecurityGroups は指定されたIDのセキュリティグループを取得します。
func (s *EC2Service) ListSecurityGroups(ctx context.Context, groupIDs []string) ([]SecurityGroup, error) {
	if len(groupIDs) == 0 {
		return nil, nil
	}
	awsSgs, err := s.repo.DescribeSecurityGroups(ctx, groupIDs)
	if err != nil {
		return nil, err
	}

	var sgs []SecurityGroup
	for _, sg := range awsSgs {
		sgs = append(sgs, SecurityGroup{
			ID:          *sg.GroupId,
			Name:        *sg.GroupName,
			Description: *sg.Description,
			Tags:        convertTags(sg.Tags),
		})
	}

	return sgs, nil
}

func convertTags(tags []types.Tag) map[string]string {
	m := make(map[string]string)
	for _, t := range tags {
		m[*t.Key] = *t.Value
	}
	return m
} 
